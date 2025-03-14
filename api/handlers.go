package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"fmt"
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
		var osmID sql.NullInt64
		var name sql.NullString
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
		http.Error(w, "Missing 'region' parameter", http.StatusBadRequest)
		return
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM juho_test.boundary WHERE name = $1)", region).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, fmt.Sprintf("Region '%s' not found in boundary table", region), http.StatusNotFound)
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
		var osmID sql.NullInt64
		var name sql.NullString
		var geom string
		var tags string

		if err := rows.Scan(&osmID, &name, &geom, &tags); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var tagMap map[string]interface{}
		json.Unmarshal([]byte(tags), &tagMap)

		results = append(results, map[string]interface{}{
			"osm_id": osmID,
			"name":   name,
			"geom":   json.RawMessage(geom),
			"tags":   tagMap,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

