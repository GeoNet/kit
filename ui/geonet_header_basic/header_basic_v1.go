package geonet_header_basic

import (
	"bytes"
	_ "embed"
	"html/template"
)

type HeaderBasicConfig struct {
	// The HTML for the logo to use.
	Logo template.HTML
	// Items to display for navigation. Can be either a link
	// or a dropdown. Note: If the first item is a link, it's
	// considered the 'home' page link.
	Items []HeaderBasicItem
	// The HTML for the home icon. This should not be changed.
	HomeIcon template.HTML
}

type HeaderBasicItem interface {
	GetTitle() string
	GetURL() string
	External() bool
	GetLinks() []HeaderBasicLink
}

// Defines a link that is displayed on the header for navigation.
type HeaderBasicLink struct {
	Title string
	URL   string
	// Whether or not the link is external (displays an external icon next to it).
	IsExternal bool
}

func (l HeaderBasicLink) GetTitle() string {
	return l.Title
}
func (l HeaderBasicLink) GetURL() string {
	return l.URL
}
func (l HeaderBasicLink) External() bool {
	return l.IsExternal
}
func (l HeaderBasicLink) GetLinks() []HeaderBasicLink {
	return []HeaderBasicLink{l}
}

// Defines a dropdown that is displayed on the header for navigation.
// Contains a number of links.
type HeaderBasicDropdown struct {
	Title string
	Links []HeaderBasicLink
}

func (d HeaderBasicDropdown) GetTitle() string {
	return d.Title
}
func (d HeaderBasicDropdown) GetURL() string {
	return ""
}
func (d HeaderBasicDropdown) External() bool {
	return false
}
func (d HeaderBasicDropdown) GetLinks() []HeaderBasicLink {
	return d.Links
}

//go:embed header-basic-v1.html
var headerBasicHtmlV1 string
var headerBasicTmplV1 = template.Must(template.New("headerbasic").Parse(headerBasicHtmlV1))

//go:embed icons/home.svg
var homeIcon template.HTML

// ReturnGeoNetHeaderBasic returns HTML for the basic GeoNet header that
// can be inserted into a webpage. The config is used to set certain properties.
func ReturnGeoNetHeaderBasic(config HeaderBasicConfig) (template.HTML, error) {
	var b bytes.Buffer
	var contents template.HTML

	config.HomeIcon = homeIcon

	if err := headerBasicTmplV1.ExecuteTemplate(&b, "headerbasic", config); err != nil {
		return contents, err
	}
	return template.HTML(b.String()), nil // nolint: gosec // The source is our HTML file.
}
