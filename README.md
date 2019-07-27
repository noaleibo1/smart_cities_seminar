# Smart Cities Seminar

This code was written as a research assignment for Smart Cities Seminar by Dr. Efrat Blumenfeld Lieberthal.

## Steps to create database table:
1. ...
2. Split with lines (Toolbox) - break lines at intersections. Reference: https://gis.stackexchange.com/questions/247013/splitting-a-polyline-at-intersections
3. Explode lines (Toolbox) - exploding lines into segments. Reference: https://gis.stackexchange.com/questions/271806/exploding-line-into-segments-using-qgis
4. Calculate length of each polyline into a new field.
5. Remove lines with length < 1.
6. Locate points along lines (external plugin) - create a layer of points from polylines. Tick "add endpoint" and "keep attributes".
7. Add x and y coordinates - in `Toolbox` search for `Add geometry attributes`.
