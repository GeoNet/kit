package seis

import (
	"strings"
)

type slStream struct {
	network   string
	station   string
	selection string
}

/*
func (s slStream) stream() string {
	switch {
	case s.network != "" && s.station != "":
		return s.network + "_" + s.station
	case s.network != "":
		return s.network + "_*"
	case s.station != "":
		return "*_" + s.station
	default:
		return "*_*"
	}
}
*/

/*
func (s slStream) encodeStream() (string, string) {

}
*/

func decodeStreams(streams, selectors string) ([]slStream, error) {

	var list []slStream
	for _, sl := range strings.Split(streams, ",") {
		stnSplit := strings.Split(sl, ":")
		var selectCmd []string
		switch {
		case len(stnSplit) > 1:
			selectCmd = strings.Fields(stnSplit[1])
		case selectors != "":
			selectCmd = strings.Split(selectors, " ")
		default:
			selectCmd = []string{"?????"}
		}

		var network, station string
		switch netSplit := strings.Split(stnSplit[0], "_"); {
		case len(netSplit) == 1:
			station, network = netSplit[0], "*"
		default:
			station, network = netSplit[1], netSplit[0]
		}

		for _, sel := range selectCmd {
			list = append(list, slStream{
				station:   station,
				network:   network,
				selection: sel,
			})
		}
	}

	return list, nil
}
