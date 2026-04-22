package geonet_header_basic

import (
	"bytes"
	_ "embed"
	"html/template"
)

type HeaderBasicConfigV2 struct {
	// Location of the geonet-design-system icons/logos.
	IconPath string
	// Items to display for navigation.
	Items []HeaderBasicItemV2
}

type HeaderBasicItemV2 interface {
	GetTitle() string
	GetURL() string
	External() bool
}

// Defines a link that is displayed on the header for navigation.
type HeaderBasicLinkV2 struct {
	Title string
	URL   string
	// Whether or not the link is external (displays an external icon next to it).
	IsExternal bool
}

func (l HeaderBasicLinkV2) GetTitle() string {
	return l.Title
}
func (l HeaderBasicLinkV2) GetURL() string {
	return l.URL
}
func (l HeaderBasicLinkV2) External() bool {
	return l.IsExternal
}

//go:embed header-basic-v2.html
var headerBasicHtmlV2 string
var headerBasicTmplV2 = template.Must(template.New("headerbasic").Parse(headerBasicHtmlV2))

// ReturnGeoNetHeaderBasicV2 returns HTML for version 2 of the basic GeoNet header that
// can be inserted into a webpage. The config is used to set certain properties.
func ReturnGeoNetHeaderBasicV2(config HeaderBasicConfigV2) (template.HTML, error) {
	var b bytes.Buffer
	var contents template.HTML

	if err := headerBasicTmplV2.ExecuteTemplate(&b, "headerbasic", config); err != nil {
		return contents, err
	}
	return template.HTML(b.String()), nil // nolint: gosec // The source is our HTML file.
}
