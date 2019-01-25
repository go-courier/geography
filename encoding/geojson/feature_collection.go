package geojson

import (
	"encoding/json"
	"github.com/go-courier/geography"
	"github.com/go-courier/geography/coordstransform"
	"github.com/go-courier/geography/maptile"
)

type FeatureCollection struct {
	coordsTransform *coordstransform.CoordsTransform
	Type            string     `json:"type"`
	Features        []*Feature `json:"features"`
}

// New FeatureCollection
func NewFeatureCollection() *FeatureCollection {
	return &FeatureCollection{
		Type:     "FeatureCollection",
		Features: make([]*Feature, 0),
	}
}

func (fc *FeatureCollection) SetCoordsTransform(coordsTransform *coordstransform.CoordsTransform) {
	fc.coordsTransform = coordsTransform
}

func (fc *FeatureCollection) AddMapTileFeature(features ...maptile.Feature) *FeatureCollection {
	for _, v := range features {
		fc.addMapTileFeature(v)
	}
	return fc
}

func (fc *FeatureCollection) addMapTileFeature(feature maptile.Feature) *FeatureCollection {
	feat := feature.ToGeom()
	geo := Geometry{
		Type: feat.Type(),
	}

	if fc.coordsTransform != nil {
		feat = feat.Project(fc.coordsTransform.ToMars)
	}

	switch feat.Type() {
	case "Point":
		point, _ := feat.(geography.Point)
		geo.Point = &point
		break
	case "MultiPoint":
		point, _ := feat.(geography.MultiPoint)
		geo.MultiPoint = &point
		break
	case "LineString":
		line, _ := feat.(geography.LineString)
		geo.LineString = &line
		break
	case "MultiLineString":
		line, _ := feat.(geography.MultiLineString)
		geo.MultiLineString = &line
		break
	case "Polygon":
		polygon, _ := feat.(geography.Polygon)
		geo.Polygon = &polygon
		break
	case "MultiPolygon":
		polygon, _ := feat.(geography.MultiPolygon)
		geo.MultiPolygon = &polygon
		break

	}

	fe := &Feature{
		ID:         feature.ID(),
		Type:       "Feature",
		Geometry:   geo,
		Properties: feature.Properties(),
	}

	fc.Features = append(fc.Features, fe)
	return fc
}

// MarshalJSON
func (fc *FeatureCollection) MarshalJSON() ([]byte, error) {
	type featureCollection FeatureCollection

	fcol := &featureCollection{
		Type: "FeatureCollection",
	}

	fcol.Features = fc.Features
	if fcol.Features == nil {
		fcol.Features = make([]*Feature, 0)
	}

	return json.Marshal(fcol)
}

func (fc *FeatureCollection) ToJSON() ([]byte, error) {
	return fc.MarshalJSON()
}
