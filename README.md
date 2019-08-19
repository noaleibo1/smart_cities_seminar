# Smart Cities Seminar

This code was written as a research assignment for Smart Cities Seminar by Dr. Efrat Blumenfeld Lieberthal.

## Steps to create database table:
1. Download OSM dataset of Illinois from https://download.geofabrik.de/north-america/us/illinois.html
2. Use osm2psql tool to translate to postgres:

    `osm2pgsql --slim --username noa --database chicago illinois-latest.osm.pbf`
    
    Database password: postgres.
3. Use only data in wanted relevant polygon - use `Select by Location` feature. Reference: https://gis.stackexchange.com/questions/61753/selecting-features-within-polygon-from-another-layer-using-qgis
4. Create a new layer out of selected features: https://gis.stackexchange.com/questions/26198/creating-new-layer-from-selection-in-qgis
5. Split with lines (Toolbox) - break lines at intersections. Reference: https://gis.stackexchange.com/questions/247013/splitting-a-polyline-at-intersections
6. Explode lines (Toolbox) - exploding lines into segments. Reference: https://gis.stackexchange.com/questions/271806/exploding-line-into-segments-using-qgis
7. Calculate length of each polyline into a new field.
8. Remove lines with length < 1.
9. Locate points along lines (external plugin) - create a layer of points from polylines. Tick "add endpoint" and "keep attributes".
10. Add x and y coordinates - in `Toolbox` search for `Add geometry attributes`.
