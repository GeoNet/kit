package geonet_header

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed header.html
var headerHTML string
var headerTmpl = template.Must(template.New("header").Parse(headerHTML))

//go:embed icons/home.svg
var homeIcon template.HTML

//go:embed icons/news.svg
var newsIcon template.HTML

//go:embed images/geonet_logo_white.svg
var geonetLogo template.HTML

type HeaderConfig struct {
	// The origin to be used at the beginning of the links in the header.
	// If nil, relative links are used.
	Origin string
	// A struct that defines which dropdown menu is currently "active", ie: the
	// dropdown menu that the current page belongs to.
	Active
	// The HTML for the GeoNet logo and home/news icons. These should not be changed.
	HomeIcon   template.HTML
	NewsIcon   template.HTML
	GeoNetLogo template.HTML
}

type Active struct {
	Home, Earthquake, Landslide, Tsunami, Volcano, DataDiscovery, DataTypes, DataAccess bool
}

// ReturnGeoNetHeader returns HTML for the main GeoNet header that
// can be inserted into a webpage. The config is used to set certain properties.
func ReturnGeoNetHeader(config HeaderConfig) (template.HTML, error) {
	var b bytes.Buffer
	var contents template.HTML

	config.HomeIcon = homeIcon
	config.NewsIcon = newsIcon
	config.GeoNetLogo = geonetLogo

	if err := headerTmpl.ExecuteTemplate(&b, "header", config); err != nil {
		fmt.Println(err)
		return contents, err
	}
	return template.HTML(b.String()), nil // nolint: gosec // The source is our HTML file.
}
