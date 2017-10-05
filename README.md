# Kit

Common packages for GeoNet applications.

## metrics

metrics is for gathering metrics.


## mmi

mmi is for Modificed Mercalli Intensity calculations in New Zealand.


## mseed

mseed is a Go wrapper for libmseed.  It is for working with miniSEED data.  If using mseed then libmseed will need
to be explicitly vendored and compiled  This will require a C compiler (eg., gcc) and make 
(possibly other packages depending on your system).  Alpine requires musl-dev

```
govendor add github.com/GeoNet/kit/cvendor/libmseed/^
cd vendor/github.com/GeoNet/kit/cvendor/libmseed
make clean 
make
```


## sc3ml

sc3ml is for working with SeisComPML files.


## shake

shake is for PGA, PGV, and MMI calculations


## slink

slink is a Go wrapper for libslink.  It is for working with SEEDlink servers.  If using slink then libslink will need
to be explicitly vendored and compiled  This will require a C compiler (eg., gcc) and make 
(possibly other packages depending on your system).  Alpine requires musl-dev

```
govendor add github.com/GeoNet/kit/cvendor/libslink/^
cd vendor/github.com/GeoNet/kit/cvendor/libslink
make clean 
make
```

## wgs84

wgs84 is for distance, bearing, and locality calculations on the WGS84 ellipsoid.