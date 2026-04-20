#!/bin/bash -e

npm install-clean

OUT_DIR=assets/dependencies
mkdir -p $OUT_DIR

MODULE=geonet-fonts
mkdir -p $OUT_DIR/$MODULE/css
mkdir -p $OUT_DIR/$MODULE/webfonts/Aspekta
mkdir -p $OUT_DIR/$MODULE/webfonts/Soehne
cp -r node_modules/$MODULE/dist/css/ $OUT_DIR/$MODULE/css
cp node_modules/$MODULE/dist/webfonts/Aspekta/Aspekta-300.woff2 $OUT_DIR/$MODULE/webfonts/Aspekta
cp -r node_modules/$MODULE/dist/webfonts/Soehne/ $OUT_DIR/$MODULE/webfonts/Soehne

MODULE=geonet-design-system
mkdir -p $OUT_DIR/$MODULE
mkdir -p $OUT_DIR/$MODULE/css
mkdir -p $OUT_DIR/$MODULE/js
mkdir -p $OUT_DIR/$MODULE/icons

cp -r node_modules/$MODULE/dist/css/ $OUT_DIR/$MODULE/css
cp -r node_modules/$MODULE/dist/js/ $OUT_DIR/$MODULE/js
icons=(
  footer-graphic.svg
  external.svg
  android-app.svg
  ios-app.svg
  chevron-down.svg
  esnz-logo.svg
  facebook.svg
  x.svg
  instagram.svg
  youtube.svg
  github.svg
  geonet.svg
  geonet-on-dark.svg
  government.svg
  close.svg
  menu.svg
  right-arrow.svg
  search.svg
  shielded.svg
  toka-tu-ake-nhc-logo.svg
)
for icon in "${icons[@]}"; do
  cp "node_modules/$MODULE/dist/icons/$icon" "$OUT_DIR/$MODULE/icons/"
done

# V1 header/footer stuff
MODULE=@fortawesome/fontawesome-free
mkdir -p $OUT_DIR/$MODULE/css
mkdir -p $OUT_DIR/$MODULE/webfonts
cp node_modules/$MODULE/css/all.min.css $OUT_DIR/$MODULE/css/all.min.css
cp node_modules/$MODULE/webfonts/fa-brands-400.woff2 $OUT_DIR/$MODULE/webfonts/fa-brands-400.woff2
cp node_modules/$MODULE/webfonts/fa-brands-400.ttf $OUT_DIR/$MODULE/webfonts/fa-brands-400.ttf
cp node_modules/$MODULE/webfonts/fa-regular-400.woff2 $OUT_DIR/$MODULE/webfonts/fa-regular-400.woff2
cp node_modules/$MODULE/webfonts/fa-regular-400.ttf $OUT_DIR/$MODULE/webfonts/fa-regular-400.ttf
cp node_modules/$MODULE/webfonts/fa-solid-900.woff2 $OUT_DIR/$MODULE/webfonts/fa-solid-900.woff2
cp node_modules/$MODULE/webfonts/fa-solid-900.ttf $OUT_DIR/$MODULE/webfonts/fa-solid-900.ttf
cp node_modules/$MODULE/webfonts/fa-v4compatibility.woff2 $OUT_DIR/$MODULE/webfonts/fa-v4compatibility.woff2
cp node_modules/$MODULE/webfonts/fa-v4compatibility.ttf $OUT_DIR/$MODULE/webfonts/fa-v4compatibility.ttf

MODULE=geonet-bootstrap
mkdir -p $OUT_DIR/$MODULE
cp node_modules/$MODULE/dist/js/bootstrap.bundle.v5.min.js $OUT_DIR/$MODULE/bootstrap.bundle.v5.min.js
cp node_modules/$MODULE/dist/js/bootstrap.bundle.v5.min.js.map $OUT_DIR/$MODULE/bootstrap.bundle.v5.min.js.map
cp node_modules/$MODULE/dist/css/bootstrap.v5.min.css $OUT_DIR/$MODULE/bootstrap.v5.min.css
cp node_modules/$MODULE/dist/css/bootstrap.v5.min.css.map $OUT_DIR/$MODULE/bootstrap.v5.min.css.map

# Copy required CSS/JS files to assets folder
mkdir -p assets/local
cp ../geonet_header_basic/header_basic.css assets/local/
cp ../geonet_footer/footer.js assets/local/
cp ../geonet_footer/footer.css assets/local/
cp ../geonet_header/header-v1.css assets/local/
cp ../geonet_header/header-v1.js assets/local/
cp ../geonet_header/header-v2.js assets/local/
