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
	http.Handle("/geonetfooter", http.HandlerFunc(testUIhandler))
	http.Handle(footer.ReturnFooterAssetPattern(), footer.ReturnFooterAssetServer())
	http.Handle("/geonetheaderbasic", http.HandlerFunc(testUIhandler))

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

	// Add leading HTML. Every page includes Bootstrap CSS/JS and fonts.
	// Each page can then also have its own custom CSS/JS files.
	leadingHTML := `<!DOCTYPE html><html><head>
	<link rel="stylesheet" href="/bootstrap.v5.min.css">
	<link rel="stylesheet" href="/font/css/font-awesome-6.1.1.min.css">`

	path := req.URL.Path
	switch path {
	case "/geonetfooter":
		config := footer.FooterConfig{
			ExtraLogos: []footer.FooterLogo{
				{
					URL:     "https://www.rcet.science",
					LogoURL: "/example_extra_logo.png",
				},
			},
		}
		leadingHTML += `<link rel="stylesheet" href="/local/footer.css">
		<script src="/local/footer.js"></script>`

		html, err = footer.ReturnGeoNetFooter(config)
		if err != nil {
			log.Println(err)
		}
	case "/geonetheader":
		config := header.HeaderConfig{
			Origin: "https://www.geonet.org.nz",
			Active: header.Active{
				Home: true,
			},
		}
		leadingHTML += `<link rel="stylesheet" href="/local/header.css">
		<script src="/local/header.js"></script>`

		html, err = header.ReturnGeoNetHeader(config)
		if err != nil {
			log.Println(err)
		}
	case "/geonetheaderbasic":
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
		leadingHTML += `<link rel="stylesheet" href="/local/header_basic.css">`

		html, err = header_basic.ReturnGeoNetHeaderBasic(config)
		if err != nil {
			log.Println(err)
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}

	leadingHTML += `<script src="/bootstrap.v5.bundle.min.js"></script></head><body>`

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
