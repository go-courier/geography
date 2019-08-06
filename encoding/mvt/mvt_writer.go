package mvt

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"

	"github.com/go-courier/geography/encoding/mvt/vector_tile"
	"github.com/go-courier/ptr"
)

func NewMVTWriter() MVTWriter {
	return &mvtWriter{
		layers: map[string]*vector_tile.Tile_Layer{},
	}
}

type mvtWriter struct {
	layerNames []string
	layers     map[string]*vector_tile.Tile_Layer
}

func (w *mvtWriter) Layers() []*vector_tile.Tile_Layer {
	layers := make([]*vector_tile.Tile_Layer, len(w.layerNames))
	for i, name := range w.layerNames {
		layers[i] = w.layers[name]
	}
	return layers
}

func (w *mvtWriter) WriteFeature(name string, extent uint32, feature Feature) {
	if feature == nil {
		return
	}

	layer, ok := w.layers[name]
	if !ok {
		layer = &vector_tile.Tile_Layer{
			Version: ptr.Uint32(2),
			Name:    &name,
			Extent:  &extent,
		}
		w.layers[name] = layer
		w.layerNames = append(w.layerNames, name)
	}

	keyIdxSet := make(map[string]uint32, 0)
	valueIdxSet := make(map[string]uint32, 0)

	addKeyValue := func(f *vector_tile.Tile_Feature, key string, value interface{}) {
		tv := vectorTileValue(value)

		keyIdx, ok := keyIdxSet[key]
		if !ok {
			layer.Keys = append(layer.Keys, key)
			keyIdxSet[key] = uint32(len(layer.Keys) - 1)
			keyIdx = keyIdxSet[key]
		}

		valueKey := tv.String()
		valueIdx, ok := valueIdxSet[valueKey]
		if !ok {
			layer.Values = append(layer.Values, tv)
			valueIdxSet[valueKey] = uint32(len(layer.Values) - 1)
			valueIdx = valueIdxSet[valueKey]
		}

		f.Tags = append(f.Tags, keyIdx, valueIdx)
	}

	geomType := feature.MVTGeometryType()

	feat := &vector_tile.Tile_Feature{
		Type:     &geomType,
		Geometry: feature.MVTGeometry(),
	}

	if featureID, ok := feature.(FeatureID); ok {
		id := featureID.ID()
		feat.Id = &id
	}

	for k, v := range feature.Properties() {
		addKeyValue(feat, k, v)
	}

	layer.Features = append(layer.Features, feat)
}

func vectorTileValue(i interface{}) *vector_tile.Tile_Value {
	tv := new(vector_tile.Tile_Value)
	switch t := i.(type) {
	default:
		buff := new(bytes.Buffer)
		err := binary.Write(buff, binary.BigEndian, t)
		if err == nil {
			tv.XXX_unrecognized = buff.Bytes()
		}
	case encoding.TextMarshaler:
		data, err := t.MarshalText()
		if err == nil {
			str := string(data)
			tv.StringValue = &str
		}
	case string:
		tv.StringValue = &t
	case fmt.Stringer:
		str := t.String()
		tv.StringValue = &str
	case bool:
		tv.BoolValue = &t
	case int8:
		intv := int64(t)
		tv.SintValue = &intv
	case int16:
		intv := int64(t)
		tv.SintValue = &intv
	case int32:
		intv := int64(t)
		tv.SintValue = &intv
	case int64:
		tv.IntValue = &t
	case uint8:
		intv := int64(t)
		tv.SintValue = &intv
	case uint16:
		intv := int64(t)
		tv.SintValue = &intv
	case uint32:
		intv := int64(t)
		tv.SintValue = &intv
	case uint64:
		tv.UintValue = &t
	case float32:
		tv.FloatValue = &t
	case float64:
		tv.DoubleValue = &t
	}
	return tv
}
