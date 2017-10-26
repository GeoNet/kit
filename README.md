# Kit

Common packages for GeoNet applications.

## Protobuf Definitions

Protobuf definitions are stored in the `protobuf` dir.  Go pkgs are generated from these files
with a `_pb` suffix e.g., `haz_pb` is generated from `protobuf/haz`.  Go protobuf pkgs import
`github.com/golang/protobuf/proto` this is not vendored here.  Pkgs are tested against the latest
version of `github.com/golang/protobuf/proto`.

## Go Packages

### gloria_pb

a generated pkg for gloria protobuf messages.


### haz_pb

a generated pkg for haz protobuf messages.


### metrics

metrics is for gathering metrics.


### mmi

mmi is for Modificed Mercalli Intensity calculations in New Zealand.


### mseed

mseed is a Go wrapper for libmseed.  It is for working with miniSEED data.  If using mseed then libmseed will need
to be explicitly vendored and compiled  This will require a C compiler (eg., gcc) and make 
(possibly other packages depending on your system).  Alpine requires musl-dev

```
govendor add github.com/GeoNet/kit/cvendor/libmseed/^
cd vendor/github.com/GeoNet/kit/cvendor/libmseed
make clean 
make
```


### sc3ml

sc3ml is for working with SeisComPML files.


### shake

shake is for PGA, PGV, and MMI calculations


### slink

slink is a Go wrapper for libslink.  It is for working with SEEDlink servers.  If using slink then libslink will need
to be explicitly vendored and compiled  This will require a C compiler (eg., gcc) and make 
(possibly other packages depending on your system).  Alpine requires musl-dev

```
govendor add github.com/GeoNet/kit/cvendor/libslink/^
cd vendor/github.com/GeoNet/kit/cvendor/libslink
make clean 
make
```

### wgs84

wgs84 is for distance, bearing, and locality calculations on the WGS84 ellipsoid.