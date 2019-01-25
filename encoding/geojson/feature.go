package geojson

import (
	"encoding/json"
)

type Feature struct {
	ID         interface{}            `json:"id,omitempty"`
	Type       string                 `json:"type"`
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
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

	return json.Marshal(fea)
}
