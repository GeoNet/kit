package seis

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"sort"
	"sync"
	"time"
)

// SLStation stores the latest state information for the given network and station combination.
type SLStation struct {
	Network   string    `json:"network"`
	Station   string    `json:"station"`
	Sequence  int       `json:"sequence"`
	Timestamp time.Time `json:"timestamp"`
}

// Key returns a blank SLStation except for the Network and Station entries, this useful as a map key.
func (s SLStation) Key() SLStation {
	return SLStation{
		Network: s.Network,
		Station: s.Station,
	}
}

// SLState maintains the current state information for a seedlink connection.
type SLState struct {
	mu   sync.Mutex
	once sync.Once

	state map[SLStation]SLStation
}

// Stations returns a sorted slice of current station state information.
func (s *SLState) Stations() []SLStation {
	s.mu.Lock()
	defer s.mu.Unlock()

	var stations []SLStation
	for _, v := range s.state {
		stations = append(stations, v)
	}
	sort.Slice(stations, func(i, j int) bool {
		switch {
		case stations[i].Network < stations[j].Network:
			return true
		case stations[i].Network > stations[j].Network:
			return false
		case stations[i].Station < stations[j].Station:
			return true
		default:
			return false
		}
	})

	return stations
}

// Add inserts or updates the station collection details into the connection state.
func (s *SLState) Add(station SLStation) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.once.Do(func() {
		s.state = make(map[SLStation]SLStation)
	})

	// there is an edge case when using wildcard options are in use and
	// different sampling rates may generate timestamp mismatches.
	s.state[station.Key()] = station
}

func (s *SLState) Find(stn SLStation) *SLStation {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range s.state {
		if ok, err := path.Match(stn.Network, k.Network); err != nil || !ok {
			continue
		}
		if ok, err := path.Match(stn.Station, k.Station); err != nil || !ok {
			continue
		}
		return &v
	}

	return nil
}

func (s *SLState) Unmarshal(data []byte) error {

	var stations []SLStation
	if err := json.Unmarshal(data, &stations); err != nil {
		return err
	}

	for _, v := range stations {
		s.Add(v)
	}

	return nil
}

func (s *SLState) Marshal() ([]byte, error) {

	data, err := json.MarshalIndent(s.Stations(), "", "  ")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *SLState) ReadFile(path string) error {

	if path == "" {
		return nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := s.Unmarshal(data); err != nil {
		return err
	}

	return nil
}

func (s *SLState) WriteFile(path string) error {

	if path == "" {
		return nil
	}

	data, err := s.Marshal()
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return err
	}

	return nil
}
