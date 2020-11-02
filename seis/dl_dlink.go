package seis

import (
	"time"
)

// DLink is a wrapper around a DLConn connection.
type DLink struct {
	server  string
	timeout time.Duration

	program  string
	username string
}

// DLinkOpt is a function for setting DLink internal parameters.
type DLinkOpt func(*DLink)

// SetDLTimeout sets the timeout for seedlink server commands and packet requests.
func SetDLTimeout(t time.Duration) DLinkOpt {
	return func(d *DLink) {
		d.timeout = t
	}
}

// SetDLProgram sets the program name for datalink connections.
func SetDLProgram(s string) DLinkOpt {
	return func(d *DLink) {
		d.program = s
	}
}

// SetDLUsername sets the username for datalink connections.
func SetDLUsername(s string) DLinkOpt {
	return func(d *DLink) {
		d.username = s
	}
}

// NewDLink returns a DLink pointer for the given server, optional settings can be passed
// as DLinkOpt functions.
func NewDLink(server string, opts ...DLinkOpt) *DLink {
	dl := DLink{
		server:   server,
		timeout:  5 * time.Second,
		program:  "seis",
		username: "seis",
	}
	for _, opt := range opts {
		opt(&dl)
	}
	return &dl
}

// SetTimeout sets the timeout value used for connection requests.
func (d *DLink) SetTimeout(t time.Duration) {
	d.timeout = t
}

// SetProgram sets the program name used for connection requests.
func (d *DLink) SetProgram(s string) {
	d.program = s
}

// SetUsername sets the username used for connection requests.
func (d *DLink) SetUsername(s string) {
	d.username = s
}

// Connect returns a DLConn pointer on a successful connection to a datalink server.
func (d *DLink) Connect() (*DLConn, error) {
	conn, err := NewDLConn(d.server, d.timeout)
	if err != nil {
		return nil, err
	}

	dl := DLConn{
		Conn: conn,
	}

	if err := dl.SetId(d.program, d.username); err != nil {
		_ = dl.Close()
		return nil, err
	}

	return &dl, nil
}
