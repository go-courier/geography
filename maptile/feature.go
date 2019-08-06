package maptile

import (
	"github.com/go-courier/geography"
	"github.com/go-courier/geography/encoding/mvt"
	"github.com/go-courier/geography/encoding/mvt/vector_tile"
)

func featureFromMVTFeature(f mvt.Feature, pixelXYToLonLat func(p geography.Point) geography.Point) (Feature, error) {
	feat := &feature{}
	feat.props = f.Properties()

	reader := mvt.NewMVTGeometryReader(f.MVTGeometry())

	switch f.MVTGeometryType() {
	case vector_tile.Tile_POINT:
		points, err := reader.ReadPoints()
		if err != nil {
			return nil, err
		}

		if len(points) == 1 {
			feat.Geom = geography.Point{points[0].X(), points[0].Y()}
		} else {
			mp := geography.MultiPoint{}
			for i := range points {
				mp = append(mp, geography.Point{points[i].X(), points[i].Y()})
			}
			feat.Geom = mp
		}
	case vector_tile.Tile_LINESTRING:
		lines, err := reader.ReadLineStrings()
		if err != nil {
			return nil, err
		}

		if len(lines) == 1 {
			feat.Geom = readLine(lines[0])
		} else {
			mls := geography.MultiLineString{}
			for i := range lines {
				mls = append(mls, readLine(lines[i]))
			}
			feat.Geom = mls
		}
	case vector_tile.Tile_POLYGON:
		p, err := reader.ReadPolygon()
		if err != nil {
			return nil, err
		}
		feat.Geom = readPolygon(p)
	}

	if feat.Geom != nil && pixelXYToLonLat != nil {
		feat.Geom = feat.Geom.Project(pixelXYToLonLat)
	}

	return feat, nil
}

func readPoint(point mvt.Coord) geography.Point {
	return geography.Point{point.X(), point.Y()}
}

func readLine(points []mvt.Coord) geography.LineString {
	line := geography.LineString{}
	for i := range points {
		line = append(line, readPoint(points[i]))
	}
	return line
}

func readPolygon(points [][]mvt.Coord) geography.Polygon {
	p := geography.Polygon{}
	for i := range points {
		p = append(p, readLine(points[i]))
	}
	return p
}

type feature struct {
	geography.Geom
	props map[string]interface{}
}

func (f *feature) ToGeom() geography.Geom {
	return f.Geom
}

func (f *feature) Properties() map[string]interface{} {
	return f.props
}

type featurePayload struct {
	Feature
	lonLatToPixelXY func(p geography.Point) geography.Point
}

func (f *featurePayload) IsValid() bool {
	return f.Feature != nil && f.ToGeom() != nil
}

func (f *featurePayload) MVTGeometryType() vector_tile.Tile_GeomType {
	return mvt.TileGeomTypeFromGeoType(f.ToGeom().Type())
}

func (f *featurePayload) MVTGeometry() []uint32 {
	g := f.ToGeom().Project(f.lonLatToPixelXY)

	mvtGeometryMarshaller, ok := g.(mvt.MVTGeometryMarshaller)
	if !ok {
		return nil
	}

	w := mvt.NewMVTGeometryWriter(mvtGeometryMarshaller.Cap())
	mvtGeometryMarshaller.MarshalMVTGeometry(w)

	return w.Data()
}

func (f *featurePayload) Properties() map[string]interface{} {
	return f.Feature.Properties()
}
