package mvt

import (
	"bytes"

	"github.com/go-courier/geography/encoding/mvt/vector_tile"
	"github.com/golang/protobuf/proto"
)

func UnmarshalMVT(data []byte, v MVTUnmarshaller) error {
	t := &vector_tile.Tile{}
	if err := proto.Unmarshal(data, t); err != nil {
		return err
	}
	w := &mvtReader{Tile: t}
	if err := v.UnmarshalMVT(w); err != nil {
		return err
	}
	return nil
}

func MarshalMVT(v MVTMarshaller) ([]byte, error) {
	w := NewMVTWriter()
	if err := v.MarshalMVT(w); err != nil {
		return nil, err
	}
	return proto.Marshal(&vector_tile.Tile{Layers: w.Layers()})
}

func ToMVT(v MVTMarshaller) (*MVT, error) {
	data, err := MarshalMVT(v)
	if err != nil {
		return nil, err
	}
	return &MVT{
		Buffer: bytes.NewBuffer(data),
	}, nil
}

type MVT struct {
	*bytes.Buffer
}

func (MVT) ContextType() string {
	return "application/vnd.mapbox-vector-tile"
}