package sc3ml

import (
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestFromSC3ML(t *testing.T) {
	for _, input := range []string{"2015p768477_0.7.xml", "2015p768477_0.8.xml", "2015p768477_0.9.xml", "2015p768477_0.10.xml"} {
		r, err := os.Open("testdata/" + input)
		if err != nil {
			t.Fatal(err)
		}

		e, err := FromSC3ML(r)
		if err != nil {
			t.Errorf("%s: %s", input, err.Error())
		}
		r.Close()

		if e.PublicID != "2015p768477" {
			t.Errorf("%s: expected publicID 2015p768477 got %s", input, e.PublicID)
		}

		if e.Type != "earthquake" {
			t.Errorf("%s: expected type earthquake got %s", input, e.Type)
		}

		if e.Time.Format(time.RFC3339Nano) != "2015-10-12T08:05:01.717692Z" {
			t.Errorf("%s: expected 2015-10-12T08:05:01.717692Z, got %s", input, e.Time.Format(time.RFC3339Nano))
		}

		if e.Latitude != -40.57806609 {
			t.Errorf("%s: Latitude expected -40.57806609 got %f", input, e.Latitude)
		}

		if e.Longitude != 176.3257242 {
			t.Errorf("%s: Longitude expected 176.3257242 got %f", input, e.Longitude)
		}

		if e.Depth != 23.28125 {
			t.Errorf("%s: Depth expected 23.28125 got %f", input, e.Depth)
		}

		if e.MethodID != "NonLinLoc" {
			t.Errorf("%s: MethodID expected NonLinLoc got %s", input, e.MethodID)
		}

		if e.EarthModelID != "nz3drx" {
			t.Errorf("%s: EarthModelID expected NonLinLoc got %s", input, e.EarthModelID)
		}

		if e.StandardError != 0.5592857863 {
			t.Errorf("%s: StandardError expected 0.5592857863 got %f", input, e.StandardError)
		}

		if e.AzimuthalGap != 166.4674465 {
			t.Errorf("%s: AzimuthalGap expected 166.4674465 got %f", input, e.AzimuthalGap)
		}

		if e.MinimumDistance != 0.1217162272 {
			t.Errorf("%s: MinimumDistance expected 0.1217162272 got %f", input, e.MinimumDistance)
		}

		if e.UsedPhaseCount != 44 {
			t.Errorf("%s: UsedPhaseCount expected 44 got %d", input, e.UsedPhaseCount)
		}

		if e.UsedStationCount != 32 {
			t.Errorf("%s: UsedStationCount expected 32 got %d", input, e.UsedStationCount)
		}

		if e.Magnitude != 5.691131913 {
			t.Errorf("%s: Magnitude expected M got %f", input, e.Magnitude)
		}

		if e.MagnitudeUncertainty != 0 {
			t.Errorf("%s: uncertainty expected 0 got %f", input, e.MagnitudeUncertainty)
		}

		if e.MagnitudeStationCount != 171 {
			t.Errorf("%s: e.MagnitudeStationCount expected 171 got %d", input, e.MagnitudeStationCount)
		}

		if e.ModificationTime.Format(time.RFC3339Nano) != "2015-10-12T22:46:41.228824Z" {
			t.Errorf("%s: Modification time expected 2015-10-12T22:46:41.228824Z got %s", input, e.ModificationTime.Format(time.RFC3339Nano))
		}
	}
}

func TestManual(t *testing.T) {
	in := []struct {
		id     string
		q      Quake
		manual bool
	}{
		{id: loc(), q: Quake{}, manual: false},
		{id: loc(), q: Quake{Type: "not existing"}, manual: true},
		{id: loc(), q: Quake{Type: "duplicate"}, manual: true},
		{id: loc(), q: Quake{EvaluationMode: "manual"}, manual: true},
		{id: loc(), q: Quake{EvaluationStatus: "confirmed"}, manual: true},
	}

	for _, v := range in {
		if v.q.Manual() != v.manual {
			t.Errorf("%s expected manual %t got %t", v.id, v.manual, v.q.Manual())
		}
	}
}

func TestPublish(t *testing.T) {
	in := []struct {
		id      string
		q       Quake
		publish bool
	}{
		{id: loc(), q: Quake{Site: "backup"}, publish: false},
		{id: loc(), q: Quake{Type: "not existing", Site: "backup"}, publish: true},
		{id: loc(), q: Quake{Type: "duplicate", Site: "backup"}, publish: true},
		{id: loc(), q: Quake{EvaluationMode: "manual", Site: "backup"}, publish: true},
		{id: loc(), q: Quake{EvaluationStatus: "confirmed", Site: "backup"}, publish: true},
		{id: loc(), q: Quake{Site: "primary"}, publish: true},
		{id: loc(), q: Quake{Type: "not existing", Site: "primary"}, publish: true},
		{id: loc(), q: Quake{Type: "duplicate", Site: "primary"}, publish: true},
		{id: loc(), q: Quake{EvaluationMode: "manual", Site: "primary"}, publish: true},
		{id: loc(), q: Quake{EvaluationStatus: "confirmed", Site: "primary"}, publish: true},
	}

	for _, v := range in {
		if v.q.Publish() != v.publish {
			t.Errorf("%s expected publish %t got %t", v.id, v.publish, v.q.Publish())
		}
	}
}

func TestStatus(t *testing.T) {
	in := []struct {
		id     string
		q      Quake
		status string
	}{
		{id: loc(), q: Quake{}, status: "automatic"},
		{id: loc(), q: Quake{Type: "not existing"}, status: "deleted"},
		{id: loc(), q: Quake{Type: "not existing", EvaluationMode: "manual"}, status: "deleted"},
		{id: loc(), q: Quake{Type: "not existing", EvaluationStatus: "confirmed"}, status: "deleted"},
		{id: loc(), q: Quake{Type: "not existing", EvaluationMode: "manual", EvaluationStatus: "confirmed"}, status: "deleted"},
		{id: loc(), q: Quake{Type: "duplicate"}, status: "duplicate"},
		{id: loc(), q: Quake{Type: "duplicate", EvaluationMode: "manual"}, status: "duplicate"},
		{id: loc(), q: Quake{Type: "duplicate", EvaluationStatus: "confirmed"}, status: "duplicate"},
		{id: loc(), q: Quake{Type: "duplicate", EvaluationMode: "manual", EvaluationStatus: "confirmed"}, status: "duplicate"},
		{id: loc(), q: Quake{EvaluationMode: "manual"}, status: "reviewed"},
		{id: loc(), q: Quake{EvaluationStatus: "confirmed"}, status: "reviewed"},
	}

	for _, v := range in {
		if v.q.Status() != v.status {
			t.Errorf("%s expected status %s got %s", v.id, v.status, v.q.Status())
		}
	}
}

func TestQuality(t *testing.T) {
	in := []struct {
		id      string
		q       Quake
		quality string
	}{
		{id: loc(), q: Quake{}, quality: "caution"},
		{id: loc(), q: Quake{EvaluationMode: "manual"}, quality: "best"},
		{id: loc(), q: Quake{Type: "not existing", EvaluationMode: "manual"}, quality: "deleted"},
		{id: loc(), q: Quake{Type: "not existing", EvaluationStatus: "confirmed"}, quality: "deleted"},
		{id: loc(), q: Quake{Type: "not existing", EvaluationMode: "manual", EvaluationStatus: "confirmed"}, quality: "deleted"},
		{id: loc(), q: Quake{UsedPhaseCount: 10, MagnitudeStationCount: 4}, quality: "caution"},
		{id: loc(), q: Quake{UsedPhaseCount: 19, MagnitudeStationCount: 10}, quality: "caution"},
		{id: loc(), q: Quake{UsedPhaseCount: 20, MagnitudeStationCount: 9}, quality: "caution"},
		{id: loc(), q: Quake{UsedPhaseCount: 20, MagnitudeStationCount: 10}, quality: "good"},
	}

	for _, v := range in {
		if v.q.Quality() != v.quality {
			t.Errorf("%s expected quality %s got %s", v.id, v.quality, v.q.Quality())
		}
	}
}

func TestAlert(t *testing.T) {
	in := []struct {
		id    string
		q     Quake
		alert bool
	}{
		{id: loc(), alert: false,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "automatic", UsedPhaseCount: 8, MagnitudeStationCount: 8, Depth: 10.0, AzimuthalGap: 200, MinimumDistance: 1.0}},
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "manual", UsedPhaseCount: 8, MagnitudeStationCount: 8, Depth: 10.0, AzimuthalGap: 200, MinimumDistance: 1.0}},
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "", EvaluationStatus: "confirmed", UsedPhaseCount: 28, MagnitudeStationCount: 28, Depth: 10.0, AzimuthalGap: 200, MinimumDistance: 1.0}},
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "", EvaluationStatus: "confirmed", UsedPhaseCount: 8, MagnitudeStationCount: 8, Depth: 10.0, AzimuthalGap: 200, MinimumDistance: 1.0}},
		{id: loc(), alert: false,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "automatic", UsedPhaseCount: 22, MagnitudeStationCount: 8, Depth: 10.0, AzimuthalGap: 200, MinimumDistance: 1.0}},
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "automatic", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 10.0, AzimuthalGap: 200, MinimumDistance: 1.0}},
		{id: loc(), alert: false,
			q: Quake{Time: time.Now().UTC(), Type: "not existing", EvaluationMode: "automatic", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 10.0, AzimuthalGap: 200, MinimumDistance: 1.0}},
		{id: loc(), alert: false,
			q: Quake{Time: time.Now().UTC(), Type: "duplicate", EvaluationMode: "automatic", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 10.0, AzimuthalGap: 200, MinimumDistance: 1.0}},
		{id: loc(), alert: false, // shallow automatic shouldn't alert
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "automatic", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 0.01, AzimuthalGap: 200, MinimumDistance: 1.0}},
		{id: loc(), alert: false, // high azimuthal gap automatic shouldn't alert
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "automatic", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 10.0, AzimuthalGap: 330, MinimumDistance: 1.0}},
		{id: loc(), alert: false, // far outside network automatic shouldn't alert
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "automatic", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 10.0, AzimuthalGap: 200, MinimumDistance: 3.0}},
		//	 combinations of bad quality parameters shouldn't alert
		{id: loc(), alert: false,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "automatic", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 10.0, AzimuthalGap: 320, MinimumDistance: 3.0}},
		{id: loc(), alert: false,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "automatic", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 0.01, AzimuthalGap: 320, MinimumDistance: 3.0}},
		//	bad quality parameters but confirmed.  Should all alert.
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "", EvaluationStatus: "confirmed", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 0.01, AzimuthalGap: 200, MinimumDistance: 1.0}},
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "", EvaluationStatus: "confirmed", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 10.0, AzimuthalGap: 330, MinimumDistance: 1.0}},
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "", EvaluationStatus: "confirmed", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 10.0, AzimuthalGap: 200, MinimumDistance: 3.0}},
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "", EvaluationStatus: "confirmed", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 10.0, AzimuthalGap: 320, MinimumDistance: 3.0}},
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "", EvaluationStatus: "confirmed", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 0.01, AzimuthalGap: 320, MinimumDistance: 3.0}},
		//	bad quality parameters but manual.  Should all alert.
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "manual", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 0.01, AzimuthalGap: 200, MinimumDistance: 1.0}},
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "manual", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 10.0, AzimuthalGap: 330, MinimumDistance: 1.0}},
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "manual", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 10.0, AzimuthalGap: 200, MinimumDistance: 3.0}},
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "manual", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 10.0, AzimuthalGap: 320, MinimumDistance: 3.0}},
		{id: loc(), alert: true,
			q: Quake{Time: time.Now().UTC(), Type: "earthquake", EvaluationMode: "manual", UsedPhaseCount: 22, MagnitudeStationCount: 12, Depth: 0.01, AzimuthalGap: 320, MinimumDistance: 3.0}},
	}

	for _, v := range in {
		alert, _ := v.q.Alert()

		if alert != v.alert {
			t.Errorf("%s incorrect alert quality got %t expected %t", v.id, alert, v.alert)
		}
	}

	for _, v := range in {
		v.q.Time = v.q.Time.Add(time.Minute * -61)

		alert, _ := v.q.Alert()

		if alert != false {
			t.Errorf("%s expected false alert quality any old event", v.id)
		}
	}
}

func TestCertainty(t *testing.T) {
	in := []struct {
		id        string
		q         Quake
		certainty string
	}{
		{id: loc(), q: Quake{}, certainty: "Possible"},
		{id: loc(), q: Quake{EvaluationMode: "manual"}, certainty: "Observed"},
		{id: loc(), q: Quake{Type: "not existing", EvaluationMode: "manual"}, certainty: "Unlikely"},
		{id: loc(), q: Quake{Type: "not existing", EvaluationStatus: "confirmed"}, certainty: "Unlikely"},
		{id: loc(), q: Quake{Type: "not existing", EvaluationMode: "manual", EvaluationStatus: "confirmed"}, certainty: "Unlikely"},
		{id: loc(), q: Quake{UsedPhaseCount: 10, MagnitudeStationCount: 4}, certainty: "Possible"},
		{id: loc(), q: Quake{UsedPhaseCount: 19, MagnitudeStationCount: 10}, certainty: "Possible"},
		{id: loc(), q: Quake{UsedPhaseCount: 20, MagnitudeStationCount: 9}, certainty: "Possible"},
		{id: loc(), q: Quake{UsedPhaseCount: 20, MagnitudeStationCount: 10}, certainty: "Likely"},
	}

	for _, v := range in {
		if v.q.Certainty() != v.certainty {
			t.Errorf("%s expected certainty %s got %s", v.id, v.certainty, v.q.Certainty())
		}
	}
}

func loc() string {
	_, _, l, _ := runtime.Caller(1)
	return "L" + strconv.Itoa(l)
}
