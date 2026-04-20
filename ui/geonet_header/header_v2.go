package geonet_header

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed header-v2.html
var headerHtmlV2 string
var headerTmplV2 = template.Must(template.New("header").Parse(headerHtmlV2))

type HeaderConfigV2 struct {
	// The origin to be used at the beginning of the links in the header.
	// If nil, relative links are used.
	Origin string
	// The index of the header item that should be highlighted (based on current page).
	// From 0 to n, counting left to right through the primary nav, and then the secondary nav.
	// For example:
	//   0: home page (nothing highlighted)
	//   1: earthquake page (any page accessed from the earthquake dropdown)
	//   8: news page
	//   11: search page.
	CurrentItem int
	// The location of the geonet-design-system icons/logos.
	IconPath string
}

// ReturnGeoNetHeaderV2 returns HTML for version 2 of the main GeoNet header that
// can be inserted into a webpage. The config is used to set certain properties.
func ReturnGeoNetHeaderV2(config HeaderConfigV2) (template.HTML, error) {
	var b bytes.Buffer
	var contents template.HTML

	if err := headerTmplV2.ExecuteTemplate(&b, "header", config); err != nil {
		fmt.Println(err)
		return contents, err
	}
	return template.HTML(b.String()), nil // nolint: gosec // The source is our HTML file.
}
