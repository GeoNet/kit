package geonet_footer

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed footer.html
var footerHTML string
var footerTmpl = template.Must(template.New("footer").Parse(footerHTML))

//go:embed images/geonet_logo.svg
var geonetLogo template.HTML

//go:embed images/gns_logo.svg
var gnsLogo template.HTML

//go:embed images/toka_tu_ake_eqc_logo.svg
var eqcLogo template.HTML

//go:embed images/toka_tu_ake_eqc_logo_stacked.svg
var eqcLogoStacked template.HTML

type FooterConfig struct {
	// The GeoNet, GNS, and EQC logos are fixed and cannot be changed.
	GeoNetLogo     template.HTML
	GnsLogo        template.HTML
	EqcLogo        template.HTML
	EqcLogoStacked template.HTML
	// URLs for extra logos to be added to the footer can be passed in.
	ExtraLogos []FooterLogo
	// Set whether footer is a stripped down, basic version.
	Basic bool
}

// Defines a logo to be displayed in the footer along with the default logos.
type FooterLogo struct {
	// The URL to link to when the image is clicked.
	URL string
	// The URL to the logo image.
	LogoURL string
}

// ReturnGeoNetFooter returns HTML for the main GeoNet footer that
// can be inserted into a webpage.
func ReturnGeoNetFooter(config FooterConfig) (template.HTML, error) {
	var b bytes.Buffer
	var contents template.HTML

	config.GeoNetLogo = geonetLogo
	config.GnsLogo = gnsLogo
	config.EqcLogo = eqcLogo
	config.EqcLogoStacked = eqcLogoStacked

	if err := footerTmpl.ExecuteTemplate(&b, "footer", config); err != nil {
		return contents, err
	}
	return template.HTML(b.String()), nil // nolint: gosec // The source is our HTML file.
}
