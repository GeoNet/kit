package dl

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	cmdId    = "ID"
	cmdWrite = "WRITE"
)

// DLConn provides connection information to a datalink service.
type DLConn struct {
	net.Conn
	timeout time.Duration

	id       string
	writable bool
	size     int
}

// NewDLConn makes a connection to a datalink service, the function
// SetId should be run after the connection is made and Close
// should be called when the link is no longer required.
func NewDLConn(service string, timeout time.Duration) (*DLConn, error) {
	if !strings.Contains(service, ":") {
		service = net.JoinHostPort(service, "18000")
	}
	client, err := net.Dial("tcp", service)
	if err != nil {
		return nil, err
	}

	conn := DLConn{
		Conn:    client,
		timeout: timeout,
	}

	return &conn, nil
}

// Id returns the client identification sent to the server.
func (d *DLConn) Id() string {
	return d.id
}

// Writable returns whether packets can be written over the connection.
func (d *DLConn) Writable() bool {
	return d.writable
}

// Size returns the current expected packet size.
func (d *DLConn) Size() int {
	return d.size
}

// SetId sens an ID message to the remote connection and decodes the connection capabilities.
func (d *DLConn) SetId(program, username string) error {

	id := fmt.Sprintf("%v:%v:%v:%v-%v", program, username, os.Getpid(), runtime.GOOS, runtime.GOARCH)

	dlp := Packet{
		Header: []byte(fmt.Sprintf("%s %s", cmdId, id)),
	}

	resp, err := d.sendPacket(dlp)
	if err != nil {
		return err
	}

	d.id = id

	// Like: ID DataLink 2014.269 :: DLPROTO:1.0 PACKETSIZE:512 WRITE
	if split := strings.Split(resp.header(), "::"); len(split) > 1 {
		for _, f := range strings.Split(split[1], " ") {
			switch {
			case f == "WRITE":
				d.writable = true
			case strings.HasPrefix(f, "PACKETSIZE:"):
				size, err := strconv.ParseInt(strings.Split(f, ":")[1], 10, 32)
				if err != nil {
					return fmt.Errorf("failed to parse packetsize: %v", err)
				}
				d.size = int(size)
			}
		}
	}

	return nil
}

// WriteMS sends a raw miniseed packet to a datalink service, it will need to be decoded
// prior to sending to allow for the required message identification values to added.
func (d *DLConn) WriteMS(srcname string, start, end time.Time, data []byte) error {

	// sanity checks
	if !d.writable {
		return fmt.Errorf("connection is not writable")
	}
	if l := len(data); l != d.size {
		return fmt.Errorf("data has incorrect length, expected %d got %d", d.size, l)
	}

	/*
		// decode miniseed data
		var msr MSRecord
		if err := msr.Unpack(data, false); err != nil {
			return err
		}
	*/

	//TODO: The 'A' request an acknowledgement, currently because we always READ the connection we fail if this isn't in place
	//TODO: Do we need a way to write without an ack (performance or something?)
	dlp := Packet{
		Header: []byte(fmt.Sprintf("%s %s/MSEED %v %v A %v",
			cmdWrite, srcname,
			hpTime(start),
			hpTime(end),
			len(data)),
		),
		Body: data,
	}

	// send packet and wait for acknowledgement
	resp, err := d.sendPacket(dlp)
	if err != nil {
		return err
	}

	if s := resp.header(); !strings.HasPrefix(s, "OK") {
		return fmt.Errorf("non-OK response message: %v", s)
	}

	return nil
}

func (d *DLConn) setDeadline() error {
	if !(d.timeout > 0) {
		return nil
	}
	return d.SetDeadline(time.Now().Add(d.timeout))
}

func (d *DLConn) sendPacket(dlp Packet) (*Packet, error) {

	if err := d.setDeadline(); err != nil {
		return nil, err
	}

	out, err := packetToBytes(dlp)
	if err != nil {
		return nil, err
	}

	if _, err = d.Write(out); err != nil {
		return nil, err
	}

	b := make([]byte, 512)

	n, err := d.Read(b)
	if err != nil {
		return nil, err
	}

	dlp, err = packetFromBytes(b[:n])
	if err != nil {
		return nil, err
	}

	if dlp.header() == "" {
		return nil, fmt.Errorf("no response in header from server")
	}

	if hSplit := strings.Split(dlp.header(), " "); len(hSplit) > 0 && hSplit[0] == "ERROR" {
		switch {
		case len(hSplit) > 1:
			return nil, fmt.Errorf("error response (%v) from server: %v", hSplit[1], dlp.body())
		default:
			return nil, fmt.Errorf("error response (unknown) from server")
		}
	}

	return &dlp, nil
}
