#!/bin/bash -e

# Copy required CSS/JS files to assets folder
mkdir -p assets/local
cp ../geonet_header/header.css assets/local/
cp ../geonet_header/header.js assets/local/
cp ../geonet_header_basic/header_basic.css assets/local/
cp ../geonet_footer/footer.js assets/local/
cp ../geonet_footer/footer.css assets/local/