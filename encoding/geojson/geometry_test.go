package geojson_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/go-courier/geography/encoding/geojson"
	"testing"
)

func TestUnmarshalGeometryCollection(t *testing.T) {
	rawJSON := `{"type": "GeometryCollection", "geometries": [
		{"type": "Point", "coordinates": [102.0, 0.5]},
		{"type": "MultiPoint", "coordinates": [[102.0, 0.5],[111,222]]},
		{"type": "LineString", "coordinates": [[102.0, 0.5],[111,222]]},
		{"type": "MultiLineString", "coordinates": [[[11,22],[33,44]],[[55,66],[77,888]]]},
		{"type": "Polygon", "coordinates": [[[11,22],[33,44]],[[55,66],[77,888]]]},
		{"type": "MultiPolygon", "coordinates": [[[[11,22],[33,44]],[[55,66],[77,888]]]]}

	]}`

	g, err := geojson.UnmarshalGeometry([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal geometry without issue, err %v", err)
	}

	if g.Type != "GeometryCollection" {
		t.Errorf("incorrect type, got %v", g.Type)
	}

	if len(g.Geometries) != 6 {
		t.Errorf("should have 6 geometries but got %d", len(g.Geometries))
	}
	spew.Dump(g)
}
