package map180

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type label struct {
	x, y, featureType int
	label             string
}

func (m *map3857) labels() (l []label, err error) {
	rows, err := db.Query(`with l as (
		select st_transScale(geom, $5, $6, $7, $8) as pt, type, name from public.map180_labels
		where
		ST_Within(geom, ST_MakeEnvelope($1,$2,$3,$4, 3857))
		AND zoom = $9
		)
 		select round(ST_X(pt)), round(ST_Y(pt)*-1), type,name from l`, m.llx, m.lly, m.urx, m.ury, m.xshift, m.yshift, m.dx, m.dx, m.zoom)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		lb := label{}
		err = rows.Scan(&lb.x, &lb.y, &lb.featureType, &lb.label)
		if err != nil {
			return
		}
		l = append(l, lb)
	}

	return
}

func labelsToSVG(labels []label) string {
	var b bytes.Buffer

	for _, l := range labels {
		l.label = strings.ReplaceAll(l.label, `Mount`, `Mt`)
		caser := cases.Title(language.BritishEnglish)
		l.label = caser.String(strings.ToLower(l.label))

		switch l.featureType {
		case 0:
			fmt.Fprintf(&b, "<circle cx=\"%d\" cy=\"%d\" r=\"1\" stroke=\"grey\" stroke-width=\"1\" fill=\"lightgrey\" />", l.x, l.y)
			fmt.Fprintf(&b, "<text fill=\"grey\" font-style=\"italic\" x=\"%d\" y=\"%d\" font-size=\"%d\" text-anchor=\"start\">%s</text>", l.x+3, l.y+5, 11, l.label)
		case 1:
			fmt.Fprintf(&b, "   <circle cx=\"%d\" cy=\"%d\" r=\"1\" stroke=\"deepskyblue\" stroke-width=\"1\" fill=\"deepskyblue\" />", l.x, l.y)
			fmt.Fprintf(&b, "<text fill=\"deepskyblue\" font-style=\"italic\" x=\"%d\" y=\"%d\" font-size=\"%d\" text-anchor=\"start\">%s</text>", l.x+3, l.y+5, 11, l.label)
		case 3:
			fmt.Fprintf(&b, "<circle cx=\"%d\" cy=\"%d\" r=\"1\" stroke=\"grey\" stroke-width=\"1\" fill=\"lightgrey\" />", l.x, l.y)
			fmt.Fprintf(&b, "<text fill=\"grey\" font-style=\"italic\" x=\"%d\" y=\"%d\" font-size=\"%d\" text-anchor=\"start\">%s</text>", l.x+3, l.y+5, 11, l.label)
		case 4:
			fmt.Fprintf(&b, "<circle cx=\"%d\" cy=\"%d\" r=\"1\" stroke=\"darkslategrey\" stroke-width=\"1\" fill=\"darkslategrey\" />", l.x, l.y)
			fmt.Fprintf(&b, "<text fill=\"darkslategrey\" font-style=\"italic\" x=\"%d\" y=\"%d\" font-size=\"%d\" text-anchor=\"start\">%s</text>", l.x+3, l.y+5, 11, l.label)
		}
	}

	return b.String()
}
