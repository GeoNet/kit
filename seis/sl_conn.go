package seis

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const slVersionFinderString = `^SeedLink v(\d)\.(\d)`
const slTimeFormat = "2006,01,02,15,04,05"

var slVersionFinder = regexp.MustCompile(slVersionFinderString)

const (
	slCmdHello = "HELLO"
	slCmdCat   = "CAT" //Not implemented by Ringserver
	slCmdClose = "BYE"

	slCmdStation = "STATION" //Enables multi-station mode: STATION station code [network code]
	slCmdEnd     = "END"     //End of handshaking for multi-station mode

	slCmdSelect = "SELECT" //   SELECT [pattern]
	slCmdData   = "DATA"   // DATA [n [begin time]]
	//slCmdFetch  = "FETCH"  // FETCH [n [begin time]]
	slCmdTime = "TIME" // TIME [begin time [end time]]

	slCmdInfoId           = "INFO ID"
	slCmdInfoCapabilities = "INFO CAPABILITIES"
	slCmdInfoStations     = "INFO STATIONS"
	slCmdInfoStreams      = "INFO STREAMS"
	slCmdInfoGaps         = "INFO GAPS"
	slCmdInfoConnections  = "INFO CONNECTIONS"
	slCmdInfoAll          = "INFO ALL"

	slCmdCrLf = "\r\n"
)

var slInfoLevel = map[string]struct {
	capability string
	command    string
}{
	"ID":           {"info:id", slCmdInfoId},
	"CAPABILITIES": {"info:capabilities", slCmdInfoCapabilities},
	"STATIONS":     {"info:stations", slCmdInfoStations},
	"STREAMS":      {"info:streams", slCmdInfoStreams},
	"GAPS":         {"info:gaps", slCmdInfoGaps},
	"CONNECTIONS":  {"info:connections", slCmdInfoConnections},
	"ALL":          {"info:all", slCmdInfoAll},
}

type SLConn struct {
	net.Conn
	timeout time.Duration

	rawVersion string
	version    struct {
		major, minor int
	}

	capabilities map[string]bool
}

// NewSLConn returns a new connection to the named seedlink server with a given command timeout. It is expected that the
// Close function be called when the connection is no longer required.
func NewSLConn(service string, timeout time.Duration) (*SLConn, error) {
	if !strings.Contains(service, ":") {
		service = net.JoinHostPort(service, "18000")
	}

	client, err := net.Dial("tcp", service)
	if err != nil {
		return nil, err
	}

	conn := SLConn{
		Conn:    client,
		timeout: timeout,
	}

	if err := conn.getCapabilities(); err != nil {
		_ = conn.Close()

		return nil, err
	}

	return &conn, nil
}

func (c *SLConn) setDeadline() error {
	if !(c.timeout > 0) {
		return nil
	}
	return c.SetDeadline(time.Now().Add(c.timeout))
}

func (c *SLConn) readPacket() (*SLPacket, error) {

	var buf bytes.Buffer
	if _, err := io.CopyN(&buf, c, SLPacketSize); err != nil {
		return nil, err
	}

	pkt, err := NewSLPacket(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return pkt, nil
}

func (c *SLConn) writeString(str string) (int, error) {
	if err := c.setDeadline(); err != nil {
		return 0, err
	}
	return c.Write([]byte(str + slCmdCrLf))
}

func (c *SLConn) infoCommand(cmd string) ([]byte, error) {

	if _, err := c.writeString(cmd); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	for {
		pkt, err := c.readPacket()
		if err != nil {
			return nil, err
		}
		offset := binary.BigEndian.Uint16(pkt.Data[44:46])
		buf.WriteString(string(pkt.Data[offset:]))

		if pkt.Seq[5] != '*' {
			break
		}

	}

	return buf.Bytes(), nil
}

func (c *SLConn) issueCommand(cmd string) ([]byte, error) {

	if _, err := c.writeString(cmd); err != nil {
		return nil, err
	}

	b := make([]byte, 512)

	i, err := c.Read(b)
	if err != nil {
		return nil, err
	}

	if s := string(b[:i]); strings.HasPrefix(s, "ERROR") {
		return nil, fmt.Errorf("got ERROR response: %v", s)
	}

	return b[:i], nil
}

func (c *SLConn) modifierCommand(cmd string) error {

	if _, err := c.writeString(cmd); err != nil {
		return err
	}

	b := make([]byte, 10)

	i, err := c.Read(b)
	if err != nil {
		return err
	}

	if s := string(b[:i]); !strings.HasPrefix(s, "OK") {
		return fmt.Errorf("non-OK response from server: %v", strings.TrimSpace(s))
	}

	return nil
}

func (c *SLConn) actionCommand(cmd string) error {

	if _, err := c.writeString(cmd); err != nil {
		return err
	}

	return nil
}

func parseSeedlinkVersion(hello string) (int, int) {
	match := slVersionFinder.FindStringSubmatch(hello)

	if len(match) == 0 {
		return 0, 0
	}

	major, _ := strconv.ParseInt(match[1], 10, 32)
	minor, _ := strconv.ParseInt(match[2], 10, 32)

	return int(major), int(minor)
}

func (c *SLConn) getCapabilities() error {
	hello, err := c.issueCommand(slCmdHello) // Use this to get some initial version/capability information.
	if err != nil {
		return fmt.Errorf("failed to issue a 'hello' command: %v", err)
	}

	c.rawVersion = string(hello)
	c.capabilities = make(map[string]bool)

	// h is like:
	// SeedLink v3.1 (2017.052 RingServer) :: SLPROTO:3.1 CAP EXTREPLY NSWILDCARD BATCH WS:13
	// GeoNet SeedLink Server
	// TODO: Can we implement EXTREPLY CAP reporting?
	// TODO: Investigate BATCH

	c.version.major, c.version.minor = parseSeedlinkVersion(string(hello))

	if caps := strings.Split(strings.Split(string(hello), slCmdCrLf)[0], "::"); len(caps) == 2 {
		for _, hc := range strings.Split(caps[1], " ") {
			c.capabilities[hc] = true
		}
	}

	capinfo, err := c.infoCommand(slCmdInfoCapabilities)
	if err != nil {
		return fmt.Errorf("unable to list capabilities: %v", err)
	}

	var info SLInfo
	if err := info.Unmarshal(capinfo); err != nil {
		return fmt.Errorf("could not parse capabilities XML: %v", err)
	}

	for _, i := range info.Capability {
		c.capabilities[i.Name] = true
	}

	return nil
}

// GetInfoLevel requests the seedlink server return an INFO request for the given level.
func (c *SLConn) GetInfo(level string) ([]byte, error) {
	info, ok := slInfoLevel[strings.ToUpper(level)]
	if !ok {
		return nil, fmt.Errorf("unknown info level: %v", level)
	}
	if !c.capabilities[info.capability] {
		return nil, fmt.Errorf("capability %s not present", info.capability)
	}

	return c.infoCommand(info.command)
}

// GetInfo requests the seedlink server return an INFO request for the given level. The results
// are returned as a decoded SLInfo pointer, or an error otherwise.
func (c *SLConn) GetSLInfo(level string) (*SLInfo, error) {
	data, err := c.GetInfo(level)
	if err != nil {
		return nil, err
	}

	var info SLInfo
	if err := info.Unmarshal(data); err != nil {
		return nil, err
	}

	return &info, nil
}

// CommandId sends an INFO ID command to the seedlink server.
func (c *SLConn) CommandId() ([]byte, error) {
	return c.infoCommand(slCmdInfoId)
}

// CommandHello sends a HELLO command to the seedlink server.
func (c *SLConn) CommandHello() ([]byte, error) {
	return c.infoCommand(slCmdHello)
}

// CommandClose sends a BYE command to the seedlink server.
func (c *SLConn) CommandClose() ([]byte, error) {
	return c.infoCommand(slCmdClose)
}

// CommandStationList sends a CAT command to the seedlink server.
func (c *SLConn) CommandCat() ([]byte, error) {
	return c.infoCommand(slCmdCat)
}

// CommandStation sends a STATION command to the seedlink server.
func (c *SLConn) CommandStation(station, network string) error {
	if strings.ContainsAny(station, "*?") && !c.capabilities["NSWILDCARD"] {
		return fmt.Errorf("station selector '%s' contains wildcards but the server does not report capability NSWILDCARD", station)
	}
	if strings.ContainsAny(network, "*?") && !c.capabilities["NSWILDCARD"] {
		return fmt.Errorf("network selector '%s' contains wildcards but the server does not report capability NSWILDCARD", network)
	}
	switch {
	case network != "":
		if err := c.modifierCommand(fmt.Sprintf("%s %s %s", slCmdStation, station, network)); err != nil {
			return fmt.Errorf("error sending STATION %s %s: %v", station, network, err)
		}
	default:
		if err := c.modifierCommand(fmt.Sprintf("%s %s", slCmdStation, station)); err != nil {
			return fmt.Errorf("error sending STATION %s: %v", station, err)
		}
	}

	return nil
}

// CommandSelect sends a SELECT command to the seedlink server.
func (c *SLConn) CommandSelect(selection string) error {

	if err := c.modifierCommand(fmt.Sprintf("%s %s", slCmdSelect, selection)); err != nil {
		return fmt.Errorf("error sending SELECT %s: %v", selection, err)
	}

	return nil
}

// CommandData sends a DATA command to the seedlink server.
func (c *SLConn) CommandData(sequence string, starttime time.Time) error {

	var dc string
	switch {
	case sequence == "":
		dc = slCmdData
	case starttime.IsZero():
		dc = fmt.Sprintf("%s %s\n", slCmdData, sequence)
	default:
		dc = fmt.Sprintf("%s %s %s\n", slCmdData, sequence, starttime.Format(slTimeFormat))
	}

	if err := c.modifierCommand(dc); err != nil {
		return fmt.Errorf("error sending DATA: %v", err)
	}

	return nil
}

// CommandTime sends a TIME command to the seedlink server.
func (c *SLConn) CommandTime(starttime, endtime time.Time) error {

	if starttime.IsZero() {
		return nil
	}

	var tc string
	switch {
	case endtime.IsZero():
		tc = fmt.Sprintf("%s %s\n", slCmdTime, starttime.Format(slTimeFormat))
	default:
		tc = fmt.Sprintf("%s %s %s\n", slCmdTime, starttime.Format(slTimeFormat), endtime.Format(slTimeFormat))
	}

	if err := c.modifierCommand(tc); err != nil {
		return fmt.Errorf("error sending TIME: %v", err)
	}

	return nil
}

// CommandEnd sends an END command to the seedlink server.
func (c *SLConn) CommandEnd() error {
	if err := c.actionCommand(slCmdEnd); err != nil {
		return fmt.Errorf("error sending END: %v", err)
	}
	return nil
}

// Collect returns a seedlink packet if available within the optional timout. Any error returned should be
// checked that it isn't a timeout, this should be handled as appropriate for the request.
func (c *SLConn) Collect() (*SLPacket, error) {
	if err := c.setDeadline(); err != nil {
		return nil, err
	}
	pkt, err := c.readPacket()
	if err != nil {
		return nil, err
	}

	return pkt, nil
}
