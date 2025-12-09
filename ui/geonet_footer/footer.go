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

//go:embed images/toka_tu_ake_nhc_logo.svg
var nhcLogo template.HTML

//go:embed images/toka_tu_ake_nhc_logo_stacked.svg
var nhcLogoStacked template.HTML

//go:embed images/esnz_logo.svg
var esnzLogo template.HTML

//go:embed images/esnz_logo_stacked.svg
var esnzLogoStacked template.HTML

type FooterConfig struct {
	// Whether to use relative links in footer. If false, uses www.geonet.org.nz.
	UseRelativeLinks bool
	// The origin to be used at the beginning of GeoNet links in the footer.
	// Cannot be changed.
	Origin string
	// The GeoNet, ESNZ, and NHC logos are fixed and cannot be changed.
	GeoNetLogo      template.HTML
	EsnzLogo        template.HTML
	EsnzLogoStacked template.HTML
	NhcLogo         template.HTML
	NhcLogoStacked  template.HTML
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
	config.EsnzLogo = esnzLogo
	config.EsnzLogoStacked = esnzLogoStacked
	config.NhcLogo = nhcLogo
	config.NhcLogoStacked = nhcLogoStacked

	config.Origin = "https://www.geonet.org.nz"
	if config.UseRelativeLinks {
		config.Origin = ""
	}

	if err := footerTmpl.ExecuteTemplate(&b, "footer", config); err != nil {
		return contents, err
	}
	return template.HTML(b.String()), nil // nolint: gosec // The source is our HTML file.
}
