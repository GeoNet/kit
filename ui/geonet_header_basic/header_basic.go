package geonet_header_basic

import (
	"bytes"
	_ "embed"
	"html/template"
)

type HeaderBasicConfig struct {
	// The HTML for the logo to use.
	Logo template.HTML
	// Links to display for navigation. Note: The first link is
	// considered the 'home' page link.
	Links []HeaderBasicLink
	// The HTML for the home icon. This should not be changed.
	HomeIcon template.HTML
}

// Defines a link that is displayed on the header for navigation.
type HeaderBasicLink struct {
	Title string
	URL   string
	// Whether or not the link is external (displays an external icon next to it).
	IsExternal bool
}

//go:embed header_basic.html
var headerBasicHTML string
var headerBasicTmpl = template.Must(template.New("headerbasic").Parse(headerBasicHTML))

//go:embed icons/home.svg
var homeIcon template.HTML

// ReturnGeoNetHeaderBasic returns HTML for the basic GeoNet header that
// can be inserted into a webpage. The config is used to set certain properties.
func ReturnGeoNetHeaderBasic(config HeaderBasicConfig) (template.HTML, error) {
	var b bytes.Buffer
	var contents template.HTML

	config.HomeIcon = homeIcon

	if err := headerBasicTmpl.ExecuteTemplate(&b, "headerbasic", config); err != nil {
		return contents, err
	}
	return template.HTML(b.String()), nil // nolint: gosec // The source is our HTML file.
}
