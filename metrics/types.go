package metrics

// Message defines an interface for processing bytes.
type Processor interface {
	Process([]byte) error
}

// DoProcess executes m.Process with messaging metrics.
func DoProcess(m Processor, b []byte) error {
	MsgRx()
	t := Start()
	defer t.Track("process")

	err := m.Process(b)

	switch err {
	case nil:
		MsgProc()
	default:
		MsgErr()
	}

	return err
}
