package data

import (
	"io"
	"log"
	"net/http"
	"os"
)

// "http://overpass-api.de/api/interpreter?data=[out:json];area[name=%22Toronto%22]-%3E.searchArea;(node[leisure=park](area.searchArea);way[leisure=park](area.searchArea);relation[leisure=park](area.searchArea););out%20body;%3E;out%20skel%20qt;"

func FetchOSMData(url string) error {
	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	jsonBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	geojsonData, err := ConvertToGeoJSON(jsonBytes)
	if err != nil {
		log.Fatal(err)
	}

	outputFile, err := os.Create("osm_data.geojson")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	geojsonBytes, err := geojsonData.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	_, err = outputFile.Write(geojsonBytes)
	return err
}
