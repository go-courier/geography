package geojson

import (
	"encoding/json"
	"github.com/go-courier/geography"
)

// 几何体关联对象
type Geometry struct {
	Type            string `json:"type"`
	Point           *geography.Point
	MultiPoint      *geography.MultiPoint
	LineString      *geography.LineString
	MultiLineString *geography.MultiLineString
	Polygon         *geography.Polygon
	MultiPolygon    *geography.MultiPolygon
	Geometries      []*Geometry
}

func (g *Geometry) MarshalJSON() ([]byte, error) {
	type geometry struct {
		Type        string      `json:"type"`
		Geometries  interface{} `json:"geometries,omitempty"`
		Coordinates interface{} `json:"coordinates,omitempty"`
	}

	geo := &geometry{
		Type: g.Type,
	}

	switch g.Type {
	case "Point":
		geo.Coordinates = g.Point
	case "MultiPoint":
		geo.Coordinates = g.MultiPoint
	case "LineString":
		geo.Coordinates = g.LineString
	case "MultiLineString":
		geo.Coordinates = g.MultiLineString
	case "Polygon":
		geo.Coordinates = g.Polygon
	case "MultiPolygon":
		geo.Coordinates = g.MultiPolygon
	case "GeometryCollection":
		geo.Geometries = g.Geometries
	}

	return json.Marshal(geo)
}
