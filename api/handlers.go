package api

import (
	"encoding/json"
	"net/http"
)

func GetRawDataHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT osm_id, name, ST_AsGeoJSON(geom), tags FROM juho_test.parks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var osmID int
		var name string
		var geom string
		var tags string

		if err := rows.Scan(&osmID, &name, &geom, &tags); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var tagMap map[string]interface{}
		json.Unmarshal([]byte(tags), &tagMap)

		results = append(results, map[string]interface{} {
			"osm_id": osmID,
			"name":   name,
			"geom":   json.RawMessage(geom),
			"tags":   tagMap,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func SpatialQueryHandler(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	if region == "" {
		http.Error(w, "Mission 'region' parameter", http.StatusBadRequest)
		return
	}

	query := `
        SELECT osm_id, name, ST_AsGeoJSON(geom), tags
        FROM juho_test.parks
        WHERE ST_Within(geom, (SELECT geom FROM juho_test.boundary WHERE name = $1));
    `

	rows, err := db.Query(query, region)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var osmID int
		var name string
		var geom string
		var tags string

		if err := rows.Scan(&osmID, &name, &geom, &tags); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var tagMap map[string]interface{}
		json.Unmarshal([]byte(tags), &tagMap)

		results = append(results, map[string]interface{} {
			"osm_id": osmID,
			"name":   name,
			"geom":   json.RawMessage(geom),
			"tags":   tagMap,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
