package maptile

import (
	"sync"

	"github.com/go-courier/geography"
	"github.com/go-courier/geography/encoding/mvt"
)

func NewMapTile(z, x, y uint32) *MapTile {
	return &MapTile{
		Z:      z,
		X:      x,
		Y:      y,
		Layers: map[string]*Layer{},
	}
}

type MapTile struct {
	coordsTransform CoordsTransform
	Z               uint32
	X               uint32
	Y               uint32
	Layers          map[string]*Layer
	LayerNames      []string
}

type CoordsTransform interface {
	ToEarth(point geography.Point) geography.Point
	ToMars(point geography.Point) geography.Point
}

func (t *MapTile) SetCoordsTransform(coordsTransform CoordsTransform) {
	t.coordsTransform = coordsTransform
}

func (t *MapTile) MarshalMVT(mvtWriter mvt.MVTWriter) error {
	for i := range t.Layers {
		layer := t.Layers[i]
		if layer == nil || len(layer.Features) == 0 {
			continue
		}

		transform := t.NewLonLatToPixelXYTransform(layer.Extent)

		for i := range layer.Features {
			feat := layer.Features[i]
			if feat == nil {
				continue
			}

			f := &featurePayload{Feature: feat, lonLatToPixelXY: transform}

			if f.IsValid() {
				mvtWriter.WriteFeature(layer.Name, layer.Extent, f)
			}
		}

	}
	return nil
}

func (t *MapTile) UnmarshalMVT(r mvt.MVTReader) error {
	return r.RangeFeature(func(name string, extent uint32, feature mvt.Feature) error {

		layer, ok := t.Layers[name]
		if !ok {
			layer = &Layer{
				Name:   name,
				Extent: extent,
			}
			t.AddLayers(layer)
		}

		transform := t.NewPixelXYToLonLatTransform(layer.Extent)

		feat, err := featureFromMVTFeature(feature, transform)
		if err != nil {
			return err
		}
		layer.Features = append(layer.Features, feat)
		return nil
	})
}

func (t *MapTile) NewLonLatToPixelXYTransform(extent uint32) geography.Transform {
	projection := mvt.NewProjection(t.Z, t.X, t.Y, extent)

	return func(p geography.Point) geography.Point {
		if t.coordsTransform != nil {
			p = t.coordsTransform.ToMars(p)
		}
		x, y := projection.LonLatToTilePixelXY(p[0], p[1])
		return geography.Point{x, y}
	}
}

func (t *MapTile) NewPixelXYToLonLatTransform(extent uint32) geography.Transform {
	projection := mvt.NewProjection(t.Z, t.X, t.Y, extent)

	return func(p geography.Point) geography.Point {
		lon, lat := projection.TilePixelXYToLonLat(p[0], p[1])
		point := geography.Point{lon, lat}

		if t.coordsTransform != nil {
			return t.coordsTransform.ToEarth(p)
		}
		return point
	}
}

func (t *MapTile) BBox() geography.Bound {
	buffer := 0.0
	x := float64(t.X)
	y := float64(t.Y)

	minx := x - buffer

	miny := y - buffer
	if miny < 0 {
		miny = 0
	}

	lon1, lat1 := geography.TileXYToLonLat(minx, miny, uint32(t.Z))

	maxX := x + 1 + buffer

	maxTiles := float64(uint32(1 << t.Z))
	maxY := y + 1 + buffer
	if maxY > maxTiles {
		maxY = maxTiles
	}

	lon2, lat2 := geography.TileXYToLonLat(maxX, maxY, uint32(t.Z))

	if t.coordsTransform != nil {
		return geography.Bound{
			Min: t.coordsTransform.ToEarth(geography.Point{lon1, lat2}),
			Max: t.coordsTransform.ToEarth(geography.Point{lon2, lat1}),
		}
	}

	return geography.Bound{
		Min: geography.Point{lon1, lat2},
		Max: geography.Point{lon2, lat1},
	}
}

func (t *MapTile) AddLayers(layers ...*Layer) {
	for i := range layers {
		layer := layers[i]
		if _, ok := t.Layers[layer.Name]; !ok {
			t.Layers[layer.Name] = layer
			t.LayerNames = append(t.LayerNames, layer.Name)
		}
	}
}

func (t *MapTile) AddTileLayers(tileLayers ...TileLayer) (e error) {
	wg := sync.WaitGroup{}

	result := make(chan interface{})

	for i := range tileLayers {
		wg.Add(1)
		go func(tileLayer TileLayer) {
			defer wg.Done()
			features, err := tileLayer.Features(t)
			if err != nil {
				result <- err
				return
			}
			extend := uint32(0)

			if tileLayerExtentConf, ok := tileLayer.(TileLayerExtentConf); ok {
				extend = tileLayerExtentConf.Extent()
			}

			result <- NewLayer(tileLayer.Name(), extend, features...)

		}(tileLayers[i])
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	for r := range result {
		switch v := r.(type) {
		case error:
			e = v
		case *Layer:
			t.AddLayers(v)
		}
	}
	return
}
