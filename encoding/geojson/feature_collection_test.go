package geojson

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
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

func TestUnmarshalFeatureCollection(t *testing.T) {
	rawJSON := `
{
    "type": "FeatureCollection",
    "features": [
        {
            "id": 1,
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    110.00424010553802,
                    19.997877614004096
                ]
            },
            "properties": {
                "age": 11,
                "name": "张三",
                "sex": "男"
            }
        },
        {
            "id": 1,
            "type": "Feature",
            "geometry": {
                "type": "LineString",
                "coordinates": [
                    [
                        110.00434539290096,
                        21.99724189683811
                    ],
                    [
                        111.0048398494024,
                        22.997247735807395
                    ]
                ]
            },
            "properties": {
                "age": 11,
                "name": "张三",
                "sex": "男"
            }
        },
        {
            "id": 1,
            "type": "Feature",
            "geometry": {
                "type": "Polygon",
                "coordinates": [
                    [
                        [
                            110.00445903309571,
                            23.997262049847986
                        ],
                        [
                            110.00445903309571,
                            23.997262049847986
                        ],
                        [
                            110.00445903309571,
                            23.997262049847986
                        ]
                    ]
                ]
            },
            "properties": {
                "age": 11,
                "name": "张三",
                "sex": "男"
            }
        }
    ]
}`

	fc, err := UnmarshalFeatureCollection([]byte(rawJSON))
	spew.Dump(fc)
	if err != nil {
		t.Fatalf("unmarshal feature collection without issue, err %v", err)
	}

	if fc.Type != "FeatureCollection" {
		t.Errorf("should have type of FeatureCollection, got %v", fc.Type)
	}

	if len(fc.Features) != 3 {
		t.Errorf("should have 3 features but got %d", len(fc.Features))
	}
}

func TestUnmarshalFeatureCollection2(t *testing.T) {
	var rawJSON = `{
    "type": "FeatureCollection",
    "features": [
        {
            "id": 4182641,
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    116.68759999999999,
                    39.891329999999996
                ]
            },
            "properties": {
                "PM10": "325.0",
                "PM2_5": "209.0",
                "VEHICLE_DIRECTION": "278.0",
                "VEHICLE_SPEED": "0.0",
                "collectTime": "2019-03-19T06:02:32+08:00",
                "deviceID": "B600-01E4",
                "name": "京B-7Y686",
                "plateNO": "京B-7Y686"
            }
        },
        {
            "id": 4182629,
            "type": "Feature",
            "geometry": {
                "type": "Point",
                "coordinates": [
                    116.67637,
                    39.89466
                ]
            },
            "properties": {
                "PM10": "285.0",
                "PM2_5": "185.0",
                "VEHICLE_DIRECTION": "356.0",
                "VEHICLE_SPEED": "1.0",
                "collectTime": "2019-03-19T06:05:18+08:00",
                "deviceID": "B600-01DA",
                "name": "京B-7Y477",
                "plateNO": "京B-7Y477"
            }
        }
    ]
}
`

	fc, err := UnmarshalFeatureCollection([]byte(rawJSON))
	require.NoError(t, err)
	spew.Dump(fc)
}

func TestUnmarshalFeatureCollection3(t *testing.T) {
	var rawJSON = `
{
    "type": "FeatureCollection",
    "name": "CountryPolygon",
    "crs": {
        "type": "name",
        "properties": {
            "name": "urn:ogc:def:crs:OGC:1.3:CRS84"
        }
    },
    "features": [{
            "type": "Feature",
            "properties": { "Id": 1234},
            "geometry": {
                "type": "Polygon","coordinates": [[
                        [113.9334716796875,34.87127685546875],
                        [113.9375,34.873901367187543],
                        [113.937683105468764,34.86669921875],
                        [113.939270019531236,34.865905761718722],
                        [113.938476562500043,34.863281250000014]
                    ]
                ]
            }
        }
    ]
}
`

	fc, err := UnmarshalFeatureCollection([]byte(rawJSON))
	require.NoError(t, err)
	spew.Dump(fc)
}
