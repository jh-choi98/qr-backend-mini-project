package data

import (
	"encoding/json"

	geojson "github.com/paulmach/go.geojson"
)

type OverpassResponse struct {
	Elements []Element `json:"elements"`
}
// `json:"elements"`: a struct tage in Go. It provides metadata about how the struct fileds should be processed.
// This tag is used by Go's encoding/json package to map JSON keys to struct fileds when parsing JSON into a struct or converting a struct back to JSON
// When reading JSON, map the 3JSON field named 'elements' to the Elements field
// When writing JSON, use 'elements' as the key in th JSON output

type Element struct {
	Type string `json:"type"`
	ID int `json:"id"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Tags map[string]string `json:"tags"`
	Nodes[]int `json:"nodes"`
}

func ConvertToGeoJSON(overpassJSON []byte) (*geojson.FeatureCollection, error) {
	var response OverpassResponse

	err := json.Unmarshal(overpassJSON, &response)
	if err != nil {
		return nil, err
	}

	featureCollection := geojson.NewFeatureCollection()
	for _, element := range response.Elements {
		var feature *geojson.Feature

		if element.Type == "node" {
			geometry := geojson.NewPointGeometry([]float64{element.Lon, element.Lat})
			feature = geojson.NewFeature(geometry)
		}

		if feature != nil {
			properties := make(map[string]interface{})
			for key, value := range element.Tags {
				properties[key] = value
			}
			feature.Properties = properties
			featureCollection.AddFeature(feature)
		}
	}
	return featureCollection, nil
}
