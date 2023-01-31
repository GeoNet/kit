# Kit

Common packages for GeoNet applications.

## Protobuf Definitions

Protobuf definitions are stored in the `protobuf` dir.  Go pkgs are generated from these files
with a `_pb` suffix e.g., `haz_pb` is generated from `protobuf/haz`.  Go protobuf pkgs import
`google.golang.org/protobuf/proto` this is not vendored here.  Pkgs are tested against the latest
version of `google.golang.org/protobuf/proto`.

(06/Aug/2021) Currently `github.com/golang/groupcache` is still using `github.com/golang/protobuf` thus an indirect still in `go.mod`. Need to keep watching if groupcache has migrated so we can get rid of this.

## Go Packages

### gloria_pb

a generated pkg for gloria protobuf messages.

### sit_pb

a generated pkg for protobufs created from the SIT repo.

### sit_delta_pb

a generated pkg for protobufs created from the Delta repo for SIT.

### haz_pb

a generated pkg for haz protobuf messages.


### metrics

metrics is for gathering metrics.


### mmi

mmi is for Modificed Mercalli Intensity calculations in New Zealand.


### sc3ml

sc3ml is for working with SeisComPML files.

### seis

- `seis/ms` is for working with miniSEED data.
- `seis/sl` is for working with SEEDlink servers.
- `seis/dl` is for working with datalink servers.


### shake

shake is for PGA, PGV, and MMI calculations


### weft

weft helps with web applications.


### wgs84

wgs84 is for distance, bearing, and locality calculations on the WGS84 ellipsoid.


### cfg
cfg helps read Postgres configuration from env var.
