#!/bin/sh

# adjust the database being connected to as you require.  This uses 'fits'. 

shp2pgsql -d -s 4326 -I /work/shp/ne_50m_land.shp  public.ne50land | psql -h localhost fits postgres
shp2pgsql -d -W "LATIN1" -s 4326 -I /work/shp/ne_50m_lakes.shp  public.ne50lakes | psql -h localhost fits postgres
shp2pgsql -d -s 4326 -I /work/shp/ne_10m_minor_islands.shp  public.ne10minorislands | psql -h localhost fits postgres
shp2pgsql -d -s 4326 -I /work/shp/ne_10m_land.shp  public.ne10land | psql -h localhost fits postgres
shp2pgsql -d -W "LATIN1" -s 4326 -I /work/shp/ne_10m_lakes.shp  public.ne10lakes | psql -h localhost fits postgres

shp2pgsql -d -s 4326 -I /work/shp/nz-chatham-is-island-polygons-topo-1250k.shp public.nztopo_1250k_chathams_land | psql -h localhost fits postgres
shp2pgsql -d -s 4326 -I /work/shp/nz-chatham-is-lagoon-polygons-topo-1250k.shp public.nztopo_1250k_chathams_lagoon | psql -h localhost fits postgres
shp2pgsql -d -s 4326 -I /work/shp/nz-coastlines-and-islands-polygons-topo-1500k.shp  public.nztopo_1500k_land | psql -h localhost fits postgres
shp2pgsql -d -s 4326 -I /work/shp/nz-coastlines-and-islands-polygons-topo-150k.shp public.nztopo_150k_land | psql -h localhost fits postgres
shp2pgsql -d -s 4326 -I /work/shp/nz-lake-polygons-topo-1500k.shp public.nztopo_1500k_lakes | psql -h localhost fits postgres
shp2pgsql -d -s 4326 -I /work/shp/nz-mainland-lake-polygons-topo-150k.shp  public.nztopo_150k_lakes| psql -h localhost fits postgres

shp2pgsql -d -s 4326 -I /work/shp/nz-kermadec-is-lake-polygons-topo-125k/nz-kermadec-is-lake-polygons-topo-125k.shp public.nztopo_125k_kermadec_lakes | psql -h localhost fits postgres

shp2pgsql -d -s 4326 -I /work/shp/nz-geographic-names-topo-1500k.shp  public.nztopo_1500k_names | psql -h localhost fits postgres
shp2pgsql -d -s 4326 -I /work/shp/nz-chatham-is-geographic-names-topo-1250k.shp  public.nztopo_1250k_chatham_names | psql -h localhost fits postgres
shp2pgsql -d -s 4326 -I /work/shp/nz-kermadec-is-geographic-names-topo-125k/nz-kermadec-is-geographic-names-topo-125k.shp  public.nztopo_125k_kermadec_names | psql -h localhost fits postgres
