# map180

map180 is a Go (golang) library for drawing SVG maps with markers on EPSG:3857.

* It handles maps that cross the 180 meridian
* Allows high zoom level data in a specified region.  
* Land and lake layer queries are cached for  speed.

Once you have a database with the map data loaded (see below) drawing a map is as simple as:

```
wm, err = map180.Init(db, `public.map180_layers`, map180.Region(`newzealand`), 256000000)
b, err := wm.SVG(...)
```

See the Go docs for further details.

## Database

Postgres 9.* with Postgis 2.*

##  Using the Assembled Data

* Create the tables (and associated indexed) `public.map180_layers` and `public.map180_labels`.  See `etc/nz_map180_layer.ddl`
* Load the data e.g., for the fits db: `psql -h 127.0.0.1 fits postgres -f data/new_zealand_map_layers.ddl`

If necessary change the schema, table, and user access as required.  They can be specificed in map180.Init()

## Assembing Data

The goal is to end up with land and lakes multi polygon on EPSG:3857 entered into `public.map180_layers` and labels in 
`public.map180_labels`.  The zoom region should include data for your region of interest at higher zoom levels.  

The assembled New Zealand data set (`data/new_zealand_map_layers.ddl`) 
was made from shape files that where loaded into the DB and then cut and transformed into `public.map180_layers` and `public.map180_layers`.  

The files `etc/load-nz-shp.sh` and `etc/nz_map180_layer.ddl` document the process of creating `public.map180_layers`.  This was then dumped using:

```
pg_dump -h 127.0.0.1 --table=public.map180_layers --table=public.map180_labels --data-only -U postgres  fits -f data/new_zealand_map_layers.ddl
```

The assembled New Zealand data set uses data sourced from:

* Natural Earth - http://www.naturalearthdata.com/
* the LINZ Data Service http://data.linz.govt.nz which is licensed by LINZ for re-use under the Creative Commons Attribution 3.0 New Zealand licence.
