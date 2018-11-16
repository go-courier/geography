package maptile

import (
	"fmt"
	"github.com/go-courier/geography"
	"github.com/go-courier/geography/encoding/mvt"
	"github.com/go-courier/geography/encoding/mvt/vector_tile"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTile(t *testing.T) {
	// https://overpass-api.de/api/map?bbox=101.250,21.943,112.500,31.952
	mt := NewMapTile(5, 25, 13)

	min := geography.Point{101.250, 21.94304553343818}
	max := geography.Point{112.5, 31.952162238024968}

	require.Equal(t, geography.Point{0, 4096}, mt.NewTransform(4096)(min))
	require.Equal(t, geography.Point{4096, 0}, mt.NewTransform(4096)(max))

	mt.AddTileLayers(LayerPoi{})

	v, _ := mvt.ToMVT(mt)
	tile := vector_tile.Tile{}
	if err := proto.Unmarshal(v.Bytes(), &tile); err != nil {
		panic(err)
	}

	for i := range tile.Layers {
		layer := tile.Layers[i]

		fmt.Printf("%s\n", *layer.Name)
		fmt.Printf("\t%d\n", *layer.Version)

		fmt.Printf("\t%v\n", layer.Keys)
		fmt.Printf("\t%v\n", layer.Values)

		for j := range layer.Features {
			f := layer.Features[j]
			fmt.Printf("\t\t%v\n", f.Type)
			fmt.Printf("\t\t%v\n", f.Tags)
		}
	}
}

type LayerPoi struct {
}

func (LayerPoi) Name() string {
	return "poi"
}

func (LayerPoi) Features(tile *MapTile) ([]Feature, error) {
	return []Feature{
		FeaturePoi{Point: geography.Point{1, 1}},
		FeaturePoi{Point: geography.Point{1, 2}},
		FeaturePoi{Point: geography.Point{1, 3}},
	}, nil
}

type FeaturePoi struct {
	geography.Point
}

func (w FeaturePoi) ToGeom() geography.Geom {
	return w.Point
}

func (FeaturePoi) Properties() map[string]interface{} {
	return map[string]interface{}{
		"name": "string",
	}
}
