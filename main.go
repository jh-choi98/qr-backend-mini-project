package main

import (
	"log"
	"qr-backend-mini-project/api"
	"qr-backend-mini-project/data"
	"qr-backend-mini-project/db"
)

const DATA_SOURCE = "http://overpass-api.de/api/interpreter?data=[out:json];area[name=%22Toronto%22]-%3E.searchArea;(node[leisure=park](area.searchArea);way[leisure=park](area.searchArea);relation[leisure=park](area.searchArea););out%20body;%3E;out%20skel%20qt;"

func main() {
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	api.InitDB(dbConn)

	err = data.FetchOSMData(DATA_SOURCE)
	if err != nil {
		log.Fatal(err)
	}

	err = db.LoadGeoJSONToPostGIS(dbConn, "osm_data.geojson")
	if err != nil {
		log.Fatal(err)
	}

	api.StartServer()
}
