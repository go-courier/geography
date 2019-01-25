package geojson

import (
	"fmt"
	"github.com/go-courier/geography"
	"github.com/go-courier/geography/coordstransform"
	"github.com/go-courier/geography/maptile"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewFeatureCollection(t *testing.T) {
	fc := NewFeatureCollection()

	if fc.Type != "FeatureCollection" {
		t.Errorf("should have type of FeatureCollection, got %v", fc.Type)
	}
}

func TestFeatureCollectionToJSON(t *testing.T) {
	fc := NewFeatureCollection()

	fc.SetCoordsTransform(&coordstransform.CoordsTransform{})

	data, err := fc.AddMapTileFeature([]maptile.Feature{
		&FeaturePoi{Geom: geography.Point{110, 20}},
		&FeaturePoi{Geom: geography.LineString{{110, 22}, {111, 23}}},
		&FeaturePoi{Geom: geography.Polygon{{{110, 24}, {110, 24}, {110, 24}}}},
	}...).ToJSON()

	require.NoError(t, err)

	fmt.Printf("%s\n", string(data))

}

type FeaturePoi struct {
	geography.Geom
}

func (*FeaturePoi) ID() uint64 {
	return 1
}

func (w *FeaturePoi) ToGeom() geography.Geom {
	return w.Geom
}

func (*FeaturePoi) Properties() map[string]interface{} {
	return map[string]interface{}{
		"name": "张三",
		"sex":  "男",
		"age":  11,
	}
}
