package geonet_footer

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed footer-v2.html
var footerHtmlV2 string
var footerTmplV2 = template.Must(template.New("footer").Parse(footerHtmlV2))

type FooterConfigV2 struct {
	// The origin to be used at the beginning of GeoNet links in the footer.
	Origin string
	// URLs for extra logos to be added to the footer can be passed in.
	ExtraLogos []FooterLogoV2
	// Set whether footer is a stripped down, basic version.
	Basic bool
	// The location of the geonet-design-system icons/logos.
	IconPath string
}

// Defines a logo to be displayed in the footer along with the default logos.
type FooterLogoV2 struct {
	// The URL to link to when the image is clicked.
	URL string
	// The URL to the logo image.
	LogoURL string
	// The alt text to add to the logo image.
	Alt string
	// Width in pixels
	Width int
	// Height in pixels
	Height int
}

// ReturnGeoNetFooter returns HTML for the main GeoNet footer that
// can be inserted into a webpage.
func ReturnGeoNetFooterV2(config FooterConfigV2) (template.HTML, error) {
	var b bytes.Buffer
	var contents template.HTML

	if err := footerTmplV2.ExecuteTemplate(&b, "footer", config); err != nil {
		return contents, err
	}

	return template.HTML(b.String()), nil // nolint: gosec // The source is our HTML file.
}
