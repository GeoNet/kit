package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	footer "github.com/GeoNet/kit/ui/geonet_footer"
	header "github.com/GeoNet/kit/ui/geonet_header"
	header_basic "github.com/GeoNet/kit/ui/geonet_header_basic"
)

//go:embed assets
var assets embed.FS

//go:embed assets/example_header_logo.svg
var logo template.HTML

func main() {

	// Setup for serving embedded asset files.
	var assetFS = fs.FS(assets)
	htmlContent, err := fs.Sub(assetFS, "assets")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.FS(htmlContent))

	http.Handle("/", fs) // Serve static files
	http.Handle("/geonetheader", http.HandlerFunc(testUIhandler))
	http.Handle("/geonetheaderv2", http.HandlerFunc(testUIhandler))
	http.Handle("/geonetheaderv1", http.HandlerFunc(testUIhandler))
	http.Handle("/geonetfooter", http.HandlerFunc(testUIhandler))
	http.Handle("/geonetfooterv2", http.HandlerFunc(testUIhandler))
	http.Handle("/geonetfooterv1", http.HandlerFunc(testUIhandler))
	http.Handle("/geonetheaderbasic", http.HandlerFunc(testUIhandler))
	http.Handle("/geonetheaderbasicv2", http.HandlerFunc(testUIhandler))
	http.Handle("/geonetheaderbasicv1", http.HandlerFunc(testUIhandler))

	log.Println("starting server")
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func testUIhandler(w http.ResponseWriter, req *http.Request) {

	var html template.HTML
	var err error

	leadingHTML := `<!DOCTYPE html><html><head>`

	path := req.URL.Path
	switch path {
	case "/geonetfooter", "/geonetfooterv2":
		config := footer.FooterConfigV2{
			Origin: "https://www.geonet.org.nz",
			ExtraLogos: []footer.FooterLogoV2{
				{
					URL:     "https://www.rcet.science",
					LogoURL: "/example_extra_logo.png",
					Alt:     "RCET logo",
					Width:   260,
					Height:  45,
				},
			},
			IconPath: "/dependencies/geonet-design-system/icons",
		}
		leadingHTML += `<link rel="stylesheet" href="/dependencies/geonet-design-system/css/geonet-design-system.css">
		<link rel="stylesheet" href="/dependencies/geonet-fonts/css/Aspekta.css">
		<link rel="stylesheet" href="/dependencies/geonet-fonts/css/Soehne.css">
		<script src="/dependencies/geonet-design-system/js/geonet-design-system.js"></script>`

		html, err = footer.ReturnGeoNetFooterV2(config)
		if err != nil {
			log.Println(err)
		}
	case "/geonetfooterv1":
		config := footer.FooterConfig{
			ExtraLogos: []footer.FooterLogo{
				{
					URL:     "https://www.rcet.science",
					LogoURL: "/example_extra_logo.png",
					Alt:     "RCET logo",
				},
			},
		}
		leadingHTML += `<link rel="stylesheet" href="/dependencies/geonet-bootstrap/bootstrap.v5.min.css">
		<link rel="stylesheet" href="/dependencies/@fortawesome/fontawesome-free/css/all.min.css">
		<link rel="stylesheet" href="/local/footer-v1.css">
		<script src="/dependencies/geonet-bootstrap/bootstrap.bundle.v5.min.js"></script>
		<script src="/local/footer-v1.js"></script>`

		html, err = footer.ReturnGeoNetFooter(config)
		if err != nil {
			log.Println(err)
		}
	case "/geonetheader", "/geonetheaderv2":
		config := header.HeaderConfigV2{
			Origin:      "https://www.geonet.org.nz",
			CurrentItem: 1,
			IconPath:    "/dependencies/geonet-design-system/icons",
		}
		leadingHTML += `<link rel="stylesheet" href="/dependencies/geonet-design-system/css/geonet-design-system.css">
		<link rel="stylesheet" href="/dependencies/geonet-fonts/css/Aspekta.css">
		<link rel="stylesheet" href="/dependencies/geonet-fonts/css/Soehne.css">
		<script type="module" src="/dependencies/geonet-design-system/js/geonet-design-system.js"></script>
		<script type="module" src="/local/header-v2.js"></script>`

		html, err = header.ReturnGeoNetHeaderV2(config)
		if err != nil {
			log.Println(err)
		}
	case "/geonetheaderv1":
		config := header.HeaderConfig{
			Origin: "https://www.geonet.org.nz",
			Active: header.Active{
				Home: true,
			},
		}
		leadingHTML += `<link rel="stylesheet" href="/dependencies/geonet-bootstrap/bootstrap.v5.min.css">
		<link rel="stylesheet" href="/dependencies/@fortawesome/fontawesome-free/css/all.min.css">
		<link rel="stylesheet" href="/local/header-v1.css">
		<script src="/dependencies/geonet-bootstrap/bootstrap.bundle.v5.min.js"></script>
		<script src="/local/header-v1.js"></script>`

		html, err = header.ReturnGeoNetHeader(config)
		if err != nil {
			log.Println(err)
		}
	case "/geonetheaderbasic", "/geonetheaderbasicv2":
		items := []header_basic.HeaderBasicItemV2{
			header_basic.HeaderBasicLinkV2{
				Title:      "Test Home",
				URL:        "https://www.geonet.org.nz",
				IsExternal: false,
			},
			header_basic.HeaderBasicLinkV2{
				Title:      "Test External",
				URL:        "https://www.geonet.org.nz",
				IsExternal: true,
			},
			header_basic.HeaderBasicLinkV2{
				Title:      "Test Not External",
				URL:        "https://www.geonet.org.nz",
				IsExternal: false,
			},
		}
		config := header_basic.HeaderBasicConfigV2{
			Items:    items,
			IconPath: "/dependencies/geonet-design-system/icons",
		}
		leadingHTML += `<link rel="stylesheet" href="/dependencies/geonet-design-system/css/geonet-design-system.css">
		<link rel="stylesheet" href="/dependencies/geonet-fonts/css/Aspekta.css">
		<link rel="stylesheet" href="/dependencies/geonet-fonts/css/Soehne.css">
		<script type="module" src="/dependencies/geonet-design-system/js/geonet-design-system.js"></script>
		<script type="module" src="/local/header-basic-v2.js"></script>`

		html, err = header_basic.ReturnGeoNetHeaderBasicV2(config)
		if err != nil {
			log.Println(err)
		}
	case "/geonetheaderbasicv1":
		items := []header_basic.HeaderBasicItem{
			header_basic.HeaderBasicLink{
				Title:      "Test Home",
				URL:        "https://www.geonet.org.nz",
				IsExternal: false,
			},
			header_basic.HeaderBasicDropdown{
				Title: "Test Dropdown",
				Links: []header_basic.HeaderBasicLink{
					{
						Title:      "Test Dropdown External",
						URL:        "https://www.geonet.org.nz",
						IsExternal: true,
					},
					{
						Title:      "Test Dropdown Not External",
						URL:        "https://www.geonet.org.nz",
						IsExternal: false,
					},
				},
			},
			header_basic.HeaderBasicLink{
				Title:      "Test External",
				URL:        "https://www.geonet.org.nz",
				IsExternal: true,
			},
			header_basic.HeaderBasicLink{
				Title:      "Test Not External",
				URL:        "https://www.geonet.org.nz",
				IsExternal: false,
			},
		}
		config := header_basic.HeaderBasicConfig{
			Logo:  logo,
			Items: items,
		}
		leadingHTML += `<link rel="stylesheet" href="/dependencies/geonet-bootstrap/bootstrap.v5.min.css">
		<link rel="stylesheet" href="/dependencies/@fortawesome/fontawesome-free/css/all.min.css">
		<link rel="stylesheet" href="/local/header-basic-v1.css">
		<script src="/dependencies/geonet-bootstrap/bootstrap.bundle.v5.min.js"></script>`

		html, err = header_basic.ReturnGeoNetHeaderBasic(config)
		if err != nil {
			log.Println(err)
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}

	leadingHTML += `</head><body>`

	// Write leading HTML to writer
	_, err = w.Write([]byte(leadingHTML))
	if err != nil {
		log.Println(err)
	}

	// Write HTML UI fragment to writer.
	_, err = w.Write([]byte(html))
	if err != nil {
		log.Println(err)
	}

	// Add trailing HTML
	trailingHTML := `</body></html>`
	_, err = w.Write([]byte(trailingHTML))
	if err != nil {
		log.Println(err)
	}
}
