package geography

import (
	"github.com/go-courier/geography/encoding/wkb"
)

type Transform func(point Point) Point

type Geom interface {
	Type() string
	ToGeom() Geom
	Clip(b Bound) Geom
	Project(transform Transform) Geom
	Bound() Bound
	Equal(g Geom) bool
}

func scan(src interface{}, g Geom) error {
	if src == nil {
		return nil
	}
	if data, ok := src.([]byte); ok {
		if err := wkb.UnmarshalWKB(data, g); err != nil {
			return err
		}
	}
	return nil
}
