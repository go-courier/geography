package geojson_test

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-courier/geography/encoding/geojson"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnmarshalFeature(t *testing.T) {
	rawJSON := `
	  { "type": "Feature",
	    "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
	    "properties": {"prop0": "value0"}
	  }`

	f, err := geojson.UnmarshalFeature([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal feature without issue, err %v", err)
	}

	if f.Type != "Feature" {
		t.Errorf("should have type of Feature, got %v", f.Type)
	}

	if len(f.Properties) != 1 {
		t.Errorf("should have 1 property but got %d", len(f.Properties))
	}
}

func TestMarshalFeatureID(t *testing.T) {
	f := &geojson.Feature{
		ID: "snail",
	}

	data, err := f.MarshalJSON()
	if err != nil {
		t.Fatalf("should marshal, %v", err)
	}

	if !bytes.Equal(data, []byte(`{"id":"snail","type":"Feature","geometry":{"type":""},"properties":null}`)) {
		t.Errorf("data not correct")
		t.Logf("%v", string(data))
	}
	f.ID = 123
	data, err = f.MarshalJSON()
	if err != nil {
		t.Fatalf("should marshal, %v", err)

	}

	if !bytes.Equal(data, []byte(`{"id":123,"type":"Feature","geometry":{"type":""},"properties":null}`)) {
		t.Errorf("data not correct")
		t.Logf("%v", string(data))
	}
}

func TestUnmarshalGeometry(t *testing.T) {
	var rawJSON = `{
    "type": "Feature",
    "properties": {
        "Id": 111
    },
    "geometry": {
        "type": "Polygon",
        "coordinates": [
            [
                [
                    11.11,
                    11.12
                ],
                [
                    12.11,
                    12.12
                ],
                [
                    13.11,
                    13.12
                ],
                [
                    14.11,
                    14.12
                ]
            ]
        ]
    }
}`
	g, err := geojson.UnmarshalFeature([]byte(rawJSON))
	require.NoError(t, err)
	spew.Dump(g)
}
