package mvt

import "github.com/go-courier/geography/encoding/mvt/vector_tile"

type MVTWriter interface {
	WriteFeature(name string, extent uint32, feature Feature)
	Layers() []*vector_tile.Tile_Layer
}

type MVTReader interface {
	RangeFeature(read func(name string, extent uint32, feature Feature) error) error
}

type MVTMarshaller interface {
	MarshalMVT(w MVTWriter) error
}

type MVTUnmarshaller interface {
	UnmarshalMVT(r MVTReader) error
}

type FeatureID interface {
	ID() uint64
}

type Feature interface {
	MVTGeometryType() vector_tile.Tile_GeomType
	MVTGeometry() []uint32
	Properties() map[string]interface{}
}

type Coord interface {
	X() float64
	Y() float64
}

type MVTGeometryWriter interface {
	MoveTo(l int, getCoord func(i int) Coord)
	LineTo(l int, getCoord func(i int) Coord)
	ClosePath()
	Data() []uint32
}

type MVTGeometryReader interface {
	ReadPoints() ([]Coord, error)
	ReadLineStrings() ([][]Coord, error)
	ReadPolygon() ([][]Coord, error)
}

type MVTGeometryMarshaller interface {
	Cap() int
	MarshalMVTGeometry(w MVTGeometryWriter)
}
