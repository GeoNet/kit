package seis

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"
)

// SLink is a wrapper around an SLConn to provide
// handling of timeouts and keep alive messages.
type SLink struct {
	server  string
	timeout time.Duration

	netto     time.Duration
	keepalive time.Duration

	state     []SLStation
	starttime time.Time
	endtime   time.Time
	sequence  int

	streams   string
	selectors string

	refresh   time.Duration
	statefile string
}

// SLink is a function for setting SLink internal parameters.
type SLinkOpt func(*SLink)

// SetSLTimeout sets the timeout for seedlink server commands and packet requests.
func SetSLTimeout(d time.Duration) SLinkOpt {
	return func(s *SLink) {
		s.timeout = d
	}
}

// SetSLRefresh sets the interval for state refreshing if enabled.
func SetSLRefresh(d time.Duration) SLinkOpt {
	return func(s *SLink) {
		s.refresh = d
	}
}

// SetSLStateFile sets the file for storing state information.
func SetSLStateFile(f string) SLinkOpt {
	return func(s *SLink) {
		s.statefile = f
	}
}

// SetSLNetTo sets the time to after which the connection is closed after no packets have been received.
func SetSLNetTo(d time.Duration) SLinkOpt {
	return func(s *SLink) {
		s.netto = d
	}
}

// SetSLKeepAlive sets the time to send an ID message to server if no packets have been received.
func SetSLKeepAlive(d time.Duration) SLinkOpt {
	return func(s *SLink) {
		s.keepalive = d
	}
}

// SetSLSequence sets the start sequence for the initial request.
func SetSLSequence(sequence int) SLinkOpt {
	return func(s *SLink) {
		s.sequence = sequence
	}
}

// SetSLStartTime sets the start of the initial request from the seedlink server.
func SetSLStartTime(t time.Time) SLinkOpt {
	return func(s *SLink) {
		s.starttime = t.UTC()
	}
}

// SetSLEndTime sets the end of the initial request from the seedlink server.
func SetSLEndTime(t time.Time) SLinkOpt {
	return func(s *SLink) {
		s.endtime = t.UTC()
	}
}

// SetSLStreams sets the list of stations and streams to from the seedlink server.
func SetSLStreams(streams string) SLinkOpt {
	return func(s *SLink) {
		s.streams = streams
	}
}

// SetSLSelectors sets the default list of selectors to use for seedlink stream requests.
func SetSLSelectors(selectors string) SLinkOpt {
	return func(s *SLink) {
		s.selectors = selectors
	}
}

// SetSLState sets the default list of station state information.
func SetSLState(stations ...SLStation) SLinkOpt {
	return func(s *SLink) {
		s.state = append(s.state, stations...)
	}
}

// NewSlink returns a SLink pointer for the given server, optional settings can be passed
// as SLinkOpt functions.
func NewSLink(server string, opts ...SLinkOpt) *SLink {
	sl := SLink{
		server:    server,
		streams:   "*_*",
		selectors: "???",
		timeout:   5 * time.Second,
		netto:     300 * time.Second,
		keepalive: 30 * time.Second,
		refresh:   300 * time.Second,
		sequence:  -1,
	}
	for _, opt := range opts {
		opt(&sl)
	}
	return &sl
}

// SetTimeout sets the timeout value used for connection requests.
func (s *SLink) SetTimeout(d time.Duration) {
	s.timeout = d
}

// SetRefresh sets the interval used for state refreshes if enabled.
func (s *SLink) SetRefresh(d time.Duration) {
	s.refresh = d
}

// SetStateFile sets the file for storing state information.
func (s *SLink) SetStateFile(f string) {
	s.statefile = f
}

// SetNetTo sets the overall timeout after which a reconnection is tried.
func (s *SLink) SetNetTo(d time.Duration) {
	s.netto = d
}

// SetKeepAlive sets the time interval needed without any packets for
// a check message is sent.
func (s *SLink) SetKeepAlive(d time.Duration) {
	s.keepalive = d
}

// SetSequence sets the start sequence for the initial request.
func (s *SLink) SetSequence(sequence int) {
	s.sequence = sequence
}

// SetStartTime sets the initial starting time of the request.
func (s *SLink) SetStartTime(t time.Time) {
	s.starttime = t.UTC()
}

// SetEndTime sets the initial end time of the request.
func (s *SLink) SetEndTime(t time.Time) {
	s.endtime = t.UTC()
}

// SetStreams sets the channel streams used for seedlink connections.
func (s *SLink) SetStreams(streams string) {
	s.streams = streams
}

// SetSelectors sets the channel selectors used for seedlink connections.
func (s *SLink) SetSelectors(selectors string) {
	s.selectors = selectors
}

// SetState sets the initial list of station state information.
func (s *SLink) SetState(stations ...SLStation) {
	s.state = append(s.state, stations...)
}

// SLCollectFunc is a function run on each returned seedlink packet. It should return a true value
// to stop collecting data without an error message. A non-nil returned error will also stop
// collection but with an assumed errored state.
type SLCollectFunc func(string, []byte) (bool, error)

// CollectWithContext makes a connection to the seedlink server, recovers initial client information and
// the sets the connection into streaming mode. Recovered packets are passed to a given function
// to process, if this function returns a true value or a non-nil error value the collection will
// stop and the function will return.
// If a call returns with a timeout error a check is made whether a keepalive is needed or whether
// the function should return as no data has been received for an extended period of time. It is
// assumed the calling function will attempt a reconnection with an updated set of options, specifically
// any start or end time parameters. The Context parameter can be used to to cancel the data collection
// independent of the function as this may never be called if no appropriate has been received.
func (s *SLink) CollectWithContext(ctx context.Context, fn SLCollectFunc) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var state SLState

	// possibly load a statefile, skip any errors to avoid looping endlessly due to a corrupt file.
	_ = state.ReadFile(s.statefile)

	for _, v := range s.state {
		state.Add(v)
	}

	list, err := decodeStreams(s.streams, s.selectors)
	if err != nil {
		return err
	}

	conn, err := NewSLConn(s.server, s.timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	refresh := s.refresh
	if !(refresh > 0) {
		refresh = time.Hour
	}

	ticker := time.NewTicker(refresh)

	for _, l := range list {
		if err := conn.CommandStation(l.station, l.network); err != nil {
			return err
		}

		if err := conn.CommandSelect(l.selection); err != nil {
			return err
		}

		sequence, starttime := s.sequence, s.starttime
		if v := state.Find(SLStation{Network: l.network, Station: l.station}); v != nil {
			sequence, starttime = v.Sequence, v.Timestamp
		}

		switch {
		// if an endtime is given then ignore statefile info.
		case !s.endtime.IsZero():
			if err := conn.CommandTime(s.starttime, s.endtime); err != nil {
				return err
			}
			// there may be a sequence number
		case !(sequence < 0):
			//convert the next sequence number into uppercase hex
			seq := fmt.Sprintf("%06X", (sequence+1)&0xffffff)
			if err := conn.CommandData(seq, starttime); err != nil {
				return err
			}
		default:
			// or check a possible start time
			if err := conn.CommandTime(starttime, time.Time{}); err != nil {
				return err
			}
		}
	}
	if err := conn.CommandEnd(); err != nil {
		return err
	}

	last := time.Now()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			if err := state.WriteFile(s.statefile); err != nil {
				return err
			}
		default:
			pkt, err := conn.Collect()
			switch {
			case err != nil:
				// could be a timeout
				switch e, ok := err.(net.Error); {
				case ok && e.Timeout():
					// hit the limit so close the connection
					if s.netto > 0 && s.netto < time.Since(last) {
						return err
					}
					// may be time for a keep alive
					if s.keepalive > 0 && s.keepalive < time.Since(last) {
						// send an ID request, ignore any results other than an error
						if _, err := conn.CommandId(); err != nil {
							return err
						}
						last = time.Now()
					}
				default:
					return err
				}
			case pkt != nil:
				// pass over to the handler function, exit if stop or an error
				if stop, err := fn(string(pkt.Seq[:]), pkt.Data[:]); err != nil || stop {
					return err
				}

				seq, err := strconv.ParseInt(string(pkt.Seq[:]), 16, 32)
				if err != nil {
					return err
				}

				var msr MSRecord
				if err := msr.Unpack(pkt.Data[:], false); err == nil {
					state.Add(SLStation{
						Network:   msr.Network(),
						Station:   msr.Station(),
						Sequence:  int(seq),
						Timestamp: msr.StartTime(),
					})
				}

				last = time.Now()
			}
		}
	}

	if err := state.WriteFile(s.statefile); err != nil {
		return err
	}

	return nil
}

// Collect calls CollectWithContext with a background Context and a handler function.
func (s *SLink) Collect(fn SLCollectFunc) error {
	return s.CollectWithContext(context.Background(), fn)
}
