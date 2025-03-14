package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"

	geojson "github.com/paulmach/go.geojson"
)

func LoadGeoJSONToPostGIS(db *sql.DB, geojsonFile string) error {
	file, err := os.Open(geojsonFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	featureCollection :=&geojson.FeatureCollection{}
	if err := json.Unmarshal(data, featureCollection); err != nil {
		return err
	}

	maxInsert := 1000
	if len(featureCollection.Features) < maxInsert {
		maxInsert = len(featureCollection.Features)
	}

	stmt, err :=db.Prepare(`
	INSERT INTO juho_test.parks (osm_id, name, geom, tags)
        VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326), $5)
        ON CONFLICT (osm_id) DO NOTHING
	`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, feature := range featureCollection.Features {
		if i >= maxInsert {
			break
		}
		
		if feature.Geometry.IsPoint() {
			lon, lat := feature.Geometry.Point[0], feature.Geometry.Point[1]
			osmID := feature.Properties["id"]
			name := feature.Properties["name"]
			tags, _ := json.Marshal(feature.Properties)

			_, err := stmt.Exec(osmID, name, lon, lat, tags)
			if err != nil {
				fmt.Printf("Failed to insert feature: %v\n", err)
			}
		}
	}

	fmt.Println("GeoJSON data loaded successfully")
	return nil
}
