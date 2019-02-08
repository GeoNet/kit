drop table public.map180_layers;
drop table public.map180_labels;

CREATE TABLE public.map180_layers (
	mapPK SERIAL PRIMARY KEY,
	region INT NOT NULL,
	zoom INT NOT NULL,
	type INT NOT NULL
);

SELECT addgeometrycolumn('public', 'map180_layers', 'geom', 3857, 'MULTIPOLYGON', 2);

CREATE INDEX ON public.map180_layers (zoom);
CREATE INDEX ON public.map180_layers (region);
CREATE INDEX ON public.map180_layers (type);
CREATE INDEX ON public.map180_layers USING gist (geom);

GRANT SELECT ON public.map180_layers TO PUBLIC;

CREATE TABLE public.map180_labels (
	labelPK SERIAL PRIMARY KEY,
	zoom INT NOT NULL,
	type INT NOT NULL,
	name text
);

SELECT addgeometrycolumn('public', 'map180_labels', 'geom', 3857, 'POINT', 2);

CREATE INDEX ON public.map180_labels (zoom);
CREATE INDEX ON public.map180_labels USING gist (geom);

GRANT SELECT ON public.map180_labels TO PUBLIC;

-- land = type 0
-- lakes = type 1

-- World. Region 0
insert into public.map180_layers (region,zoom,type,geom) select 0,0,0,  ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(-180,-85,180,85, 4326), geom),3857))
 from public.ne50land;
 -- lakes
 insert into public.map180_layers (region,zoom,type,geom) select 0,0,0,  ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(-180,-85,180,85, 4326), geom),3857))
 from public.ne50lakes;


-- New Zealand.  Region 1 
-- NE50  Left and right of 180.  
insert into public.map180_layers (region,zoom,type,geom) select 1,0,0,  ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(160,-48,180,-27, 4326), geom),3857))
 from public.ne50land where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(160,-48,180,-27, 4326), geom));
 insert into public.map180_layers (region,zoom,type,geom) select 1,0,0,  ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(-180,-48,-175,-27, 4326), geom),3857)) 
 	from public.ne50land where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(-180,-48,-175,-27, 4326), geom));
 insert into public.map180_layers (region,zoom,type,geom) select 1,0,1,  ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(160,-48,180,-27, 4326), geom),3857))
 from public.ne50lakes where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(160,-48,180,-27, 4326), geom));
 insert into public.map180_layers (region,zoom,type,geom) select 1,0,1,  ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(-180,-48,-160,-27, 4326), geom),3857)) 
 	from public.ne50lakes where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(-180,-48,-160,-27, 4326), geom));	

-- NE10  Left and right of 180.  
insert into public.map180_layers (region,zoom,type,geom) select 1,1,0,  ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(160,-48,180,-27, 4326), geom),3857))
 from public.ne10land where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(160,-48,180,-27, 4326), geom));
 insert into public.map180_layers (region,zoom,type,geom) select 1,1,0,  ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(-180,-48,-175,-27, 4326), geom),3857)) 
 	from public.ne10land where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(-180,-48,-175,-27, 4326), geom)); 	
 insert into public.map180_layers (region,zoom,type,geom) select 1,1,1,  ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(160,-48,180,-27, 4326), geom),3857))
 from public.ne10lakes where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(160,-48,180,-27, 4326), geom));
 insert into public.map180_layers (region,zoom,type,geom) select 1,1,1,  ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(-180,-48,-175,-27, 4326), geom),3857)) 
 	from public.ne10lakes where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(-180,-48,-175,-27, 4326), geom)); 	

-- NZTOPO
--
-- zoom 2: 1500k small feaures removed for performance 
--
insert into public.map180_layers (region,zoom,type,geom) select 1,2,0,  
	ST_Multi(ST_Transform(geom,3857)) from public.nztopo_1500k_land where st_area(geom) *111*111 > 0.5 ;
insert into public.map180_layers (region,zoom,type,geom) select 1,2,1,  
	ST_Multi(ST_Transform(geom,3857)) from public.nztopo_1500k_lakes  where st_area(geom) *111*111 > 0.5;

-- CVZ small water features
insert into public.map180_layers (region,zoom,type,geom) select 1,2,1,  
	ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(175.45,-39.4,175.8,-39.0, 4326), geom),3857)) from public.nztopo_1500k_lakes  
 where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(175.45,-39.4,175.8,-39.0, 4326), geom)) and st_area(geom) *111*111 <= 0.5;


--  Raoul is missing from 1500k.  Add it using filtered 50k
insert into public.map180_layers (region,zoom,type,geom) select 1,2,0,  
	ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(182.0,-29.30,182.14,-29.22, 4326), geom),3857)) 
 from public.nztopo_150k_land where 
  not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(182.0,-29.30,182.14,-29.22, 4326), geom)) and st_area(geom) *111 *111 > 0.5;


-- Chathams missing from 1500k, use 1250k files.
insert into public.map180_layers (region,zoom,type,geom) select 1,2,0,  ST_Multi(ST_Transform(geom,3857))
 from public.nztopo_1250k_chathams_land where st_area(geom) *111*111 > 0.5 ;
 insert into public.map180_layers (region,zoom,type,geom) select 1,2,1,  ST_Multi(ST_Transform(geom,3857))
 from public.nztopo_1250k_chathams_lagoon ;

--
-- zoom 3: 150k small feaures removed for performance 
--
insert into public.map180_layers (region,zoom,type,geom) select 1,3,0,  
	ST_Multi(ST_Transform(geom,3857)) from public.nztopo_150k_land where st_area(geom) *111*111 > 0.5;
insert into public.map180_layers (region,zoom,type,geom) select 1,3,1,  
	ST_Multi(ST_Transform(geom,3857)) from public.nztopo_150k_lakes where st_area(geom) *111*111 > 0.5 ;

-- CVZ small water features.
insert into public.map180_layers (region,zoom,type,geom) select 1,3,1,  
	ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(175.45,-39.4,175.8,-39.0, 4326), geom),3857)) from public.nztopo_150k_land
where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(175.45,-39.4,175.8,-39.0, 4326), geom)) and st_area(geom) *111*111 <= 0.5;


-- Raoul small features.
insert into public.map180_layers (region,zoom,type,geom) select 1,3,0,  
	ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(182.0,-29.30,182.14,-29.22, 4326), geom),3857)) 
 from public.nztopo_150k_land where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(182.0,-29.30,182.14,-29.22, 4326), geom)) 
 and st_area(geom) *111*111 <= 0.5;	
insert into public.map180_layers (region,zoom,type,geom) select 1,3,1,  
	ST_Multi(ST_Transform(geom,3857)) 
from public.nztopo_125k_kermadec_lakes;

-- White Island small features.
insert into public.map180_layers (region,zoom,type,geom) select 1,3,0,  
	ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(177.164,-37.54,177.20,-37.505, 4326), geom),3857)) 
 from public.nztopo_150k_land where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(177.164,-37.54,177.20,-37.505, 4326), geom)) 
  and st_area(geom) *111*111 <= 0.5;	
insert into public.map180_layers (region,zoom,type,geom) select 1,3,1,  
	ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(177.164,-37.54,177.20,-37.505, 4326), geom),3857)) 
 from public.nztopo_150k_lakes where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(177.164,-37.54,177.20,-37.505, 4326), geom)) 
 and st_area(geom) *111*111 <= 0.5;

-- Chathams small features.
insert into public.map180_layers (region,zoom,type,geom) select 1,3,0,  
	ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(183,-44.5,184,-43.5, 4326), geom),3857)) 
 from public.nztopo_150k_land where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(183,-44.5,184,-43.5, 4326), geom)) 
 and st_area(geom) *111*111 <= 0.5;	
insert into public.map180_layers (region,zoom,type,geom) select 1,3,1,  
	ST_Multi(ST_Transform(ST_Intersection(ST_MakeEnvelope(183,-44.5,184,-43.5, 4326), geom),3857)) 
 from public.nztopo_150k_lakes where Not ST_IsEmpty(ST_Intersection(ST_MakeEnvelope(183,-44.5,184,-43.5, 4326), geom)) 
 and st_area(geom) *111*111 <= 0.5;

-- map labels
insert into public.map180_labels (type,zoom,name,geom) select 0, 2, name, ST_Transform(geom, 3857) 
	from public.nztopo_1500k_names where desc_code = 'HILL' AND size >5;
insert into public.map180_labels (type,zoom,name,geom) select 1, 2, name, ST_Transform(geom, 3857) 
	from public.nztopo_1500k_names where desc_code = 'LAKE' AND size >5;
insert into public.map180_labels (type,zoom,name,geom) select 3, 2, name, ST_Transform(geom, 3857) 
	from public.nztopo_1500k_names where desc_code = 'ISLD' AND size >5;
insert into public.map180_labels (type,zoom,name,geom) select 4, 2, name, ST_Transform(geom, 3857) 
	from public.nztopo_1500k_names where desc_code = 'TOWN' AND size >5;

insert into public.map180_labels (type,zoom,name,geom) select 0, 3, name, ST_Transform(geom, 3857) 
	from public.nztopo_1500k_names where desc_code = 'HILL' AND size >=4;
insert into public.map180_labels (type,zoom,name,geom) select 1, 3, name, ST_Transform(geom, 3857) 
	from public.nztopo_1500k_names where desc_code = 'LAKE' AND size >=4;
insert into public.map180_labels (type,zoom,name,geom) select 3, 3, name, ST_Transform(geom, 3857) 
	from public.nztopo_1500k_names where desc_code = 'ISLD' AND size >=4;
insert into public.map180_labels (type,zoom,name,geom) select 4, 3, name, ST_Transform(geom, 3857) 
	from public.nztopo_1500k_names where desc_code = 'TOWN' AND size >5;

-- Add a few things at lower zoom level.
insert into public.map180_labels (type,zoom,name,geom) select 0, 2, name, ST_Transform(geom, 3857) 
	from public.nztopo_1500k_names where desc_code = 'HILL' AND name = 'Mount Tongariro';
insert into public.map180_labels (type,zoom,name,geom) select 0, 2, name, ST_Transform(geom, 3857) 
	from public.nztopo_1500k_names where desc_code = 'HILL' AND name = 'Mount Ngauruhoe';

update public.map180_labels set name = 'Mount Taranaki/Egmont' where name = 'Mount Taranaki or Mount Egmont';

-- Raoul
insert into public.map180_labels (type,zoom,name,geom) select 3, 2, name, ST_Transform(geom, 3857) 
	from public.nztopo_125k_kermadec_names where desc_code = 'ISLD' AND size >7;

insert into public.map180_labels (type,zoom,name,geom) select 0, 3, name, ST_Transform(geom, 3857) 
	from public.nztopo_125k_kermadec_names where desc_code = 'HILL' AND size >5;
insert into public.map180_labels (type,zoom,name,geom) select 1, 3, name, ST_Transform(geom, 3857) 
	from public.nztopo_125k_kermadec_names where desc_code = 'LAKE' AND size >4.5;
insert into public.map180_labels (type,zoom,name,geom) select 3, 3, name, ST_Transform(geom, 3857) 
	from public.nztopo_125k_kermadec_names where desc_code = 'ISLD' AND size >7;

delete from map180_labels where name = 'Kermadec Islands';
delete from public.map180_labels where name like '%Raoul%' and ST_X(geom)::int = -19802996;

-- Chatham
insert into public.map180_labels (type,zoom,name,geom) select 3, 2, name, ST_Transform(geom, 3857) 
	from public.nztopo_1250k_chatham_names where desc_code = 'ISLD' AND size >7;
insert into public.map180_labels (type,zoom,name,geom) select 4, 2, name, ST_Transform(geom, 3857) 
	from public.nztopo_1250k_chatham_names where desc_code = 'POPL' AND size >4.5;

insert into public.map180_labels (type,zoom,name,geom) select 3, 3, name, ST_Transform(geom, 3857) 
	from public.nztopo_1250k_chatham_names where desc_code = 'ISLD' AND size >=7;
insert into public.map180_labels (type,zoom,name,geom) select 4, 3, name, ST_Transform(geom, 3857) 
	from public.nztopo_1250k_chatham_names where desc_code = 'POPL' AND size >4.5;

DROP TABLE public.ne50land;
DROP TABLE public.ne50lakes;
DROP TABLE public.ne10minorislands;
DROP TABLE public.ne10land;
DROP TABLE public.ne10lakes;
DROP TABLE public.nztopo_1250k_chathams_land;
DROP TABLE public.nztopo_1250k_chathams_lagoon;
DROP TABLE public.nztopo_1500k_land;
DROP TABLE public.nztopo_150k_land;
DROP TABLE public.nztopo_1500k_lakes;
DROP TABLE public.nztopo_150k_lakes;
DROP TABLE public.nztopo_1500k_names;
DROP TABLE public.nztopo_125k_kermadec_names;
DROP TABLE public.nztopo_125k_kermadec_lakes;
DROP TABLE public.nztopo_1250k_chatham_names;