package geojson

import (
	"encoding/json"
)

type Feature struct {
	ID         interface{}            `json:"id,omitempty"`
	Type       string                 `json:"type"`
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
	CRS        map[string]interface{} `json:"crs,omitempty"`
}

func (f *Feature) MarshalJSON() ([]byte, error) {
	type feature Feature
	fea := &feature{
		ID:       f.ID,
		Type:     "Feature",
		Geometry: f.Geometry,
	}

	if f.Properties != nil && len(f.Properties) != 0 {
		fea.Properties = f.Properties
	}

	if f.CRS != nil && len(f.CRS) != 0 {
		fea.CRS = f.CRS
	}

	return json.Marshal(fea)
}

func UnmarshalFeature(data []byte) (*Feature, error) {
	f := &Feature{}
	err := json.Unmarshal(data, f)
	if err != nil {
		return nil, err
	}

	return f, nil
}
