package mvt

import (
	"github.com/go-courier/geography/encoding/mvt/vector_tile"
)

type mvtReader struct {
	*vector_tile.Tile
}

type featurePayload struct {
	Layer   *vector_tile.Tile_Layer
	Feature *vector_tile.Tile_Feature
}

func (f *featurePayload) MVTGeometryType() vector_tile.Tile_GeomType {
	return *f.Feature.Type
}

func (f featurePayload) MVTGeometry() []uint32 {
	return f.Feature.Geometry
}

func (f *featurePayload) Properties() map[string]interface{} {
	properties := map[string]interface{}{}

	for i := 0; i < len(f.Feature.Tags)/2; i++ {
		keyIdx := f.Feature.Tags[i*2]
		valueIdx := f.Feature.Tags[i*2+1]
		key := f.Layer.Keys[keyIdx]
		tv := f.Layer.Values[valueIdx]

		if tv.BoolValue != nil {
			properties[key] = *tv.BoolValue
		}
		if tv.StringValue != nil {
			properties[key] = *tv.StringValue
		}
		if tv.IntValue != nil {
			properties[key] = *tv.IntValue
		}
		if tv.FloatValue != nil {
			properties[key] = *tv.FloatValue
		}
		if tv.DoubleValue != nil {
			properties[key] = *tv.DoubleValue
		}
		if tv.UintValue != nil {
			properties[key] = *tv.UintValue
		}
		if tv.SintValue != nil {
			properties[key] = *tv.SintValue
		}
	}

	return properties
}

func (w *mvtReader) RangeFeature(read func(name string, extent uint32, feature Feature) error) error {
	for i := range w.Tile.Layers {
		layer := w.Tile.Layers[i]
		name := *layer.Name
		extent := *layer.Extent

		for i := range layer.Features {
			if err := read(name, extent, &featurePayload{layer, layer.Features[i]}); err != nil {
				return err
			}
		}
	}
	return nil
}
