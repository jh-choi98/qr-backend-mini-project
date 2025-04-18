# OSM Geospatial API Backend

![Go CI](https://github.com/jh-choi98/qr-backend-mini-project/actions/workflows/go.yml/badge.svg)

This project provides a **backend API** for working with geospatial data from **OpenStreetMap (OSM)**.  
It fetches OSM data, converts it to **GeoJSON**, stores it in a **PostGIS** database, and serves spatial queries via REST API.

## Features

✅ Fetches geospatial data (e.g., parks) from **OSM Overpass API**  
✅ Converts OSM JSON response to **GeoJSON**  
✅ Stores geospatial data in a **PostGIS database** (Supabase)  
✅ Provides RESTful APIs to:

- Retrieve all stored geospatial data
- Perform spatial queries (e.g., find parks within a region)

---

## 🗂 Project Structure

```
/qr-backend-mini-project
│── /data                 # OSM data processing
│   ├── fetch_osm.go      # Fetches data from Overpass API
│   ├── convert_geojson.go # Converts OSM JSON → GeoJSON
│── /db                   # Database logic
│   ├── connect.go        # Database connection (PostGIS)
│   ├── load_data.go      # Loads GeoJSON into PostGIS
│── /api                  # API handlers
│   ├── handler.go        # API endpoints
│   ├── server.go         # Starts the API server
│── main.go               # Main execution file
```

---

## Setup & Installation

### Clone the Repository

````sh
git clone https://github.com/jh-choi98/qr-backend-mini-project.git
cd qr-backend-mini-project


### 2️⃣ Set Up Environment Variables

```ini
CONNECT_STRING="user=supabase_admin password=<your_password> host=<your_host> port=5432 dbname=<your_db> sslmode=disable"
OSM_API_URL="http://overpass-api.de/api/interpreter?data=[out:json];area[name=%22Toronto%22]-%3E.searchArea;(node[leisure=park](area.searchArea);way[leisure=park](area.searchArea);relation[leisure=park](area.searchArea););out%20body;%3E;out%20skel%20qt;"
````

---

### 3️⃣ Install Dependencies

```sh
go mod tidy
```

---

### 4️⃣ Run the Backend

```sh
go run main.go
```

---

## 📡 API Endpoints

### 1️⃣ Get All Geospatial Data

**Example Request:**

```sh
curl http://localhost:8080/get-raw-data
```

**Example Response:**

```json
[
  {
    "osm_id": 123456,
    "name": "High Park",
    "geom": { "type": "Point", "coordinates": [-79.466, 43.646] },
    "tags": { "leisure": "park", "name": "High Park" }
  }
]
```

---

### 2️⃣ Perform a Spatial Query

**Example Request:**

```sh
curl "http://localhost:8080/spatial-query?region=Toronto"
```

**Example Response:**

```json
[
  {
    "osm_id": 654321,
    "name": "Trinity Bellwoods Park",
    "geom": { "type": "Point", "coordinates": [-79.416, 43.645] },
    "tags": { "leisure": "park", "name": "Trinity Bellwoods Park" }
  }
]
```

---

## 🔧 Troubleshooting

### Database connection error (`pq: relation "juho_test.parks" does not exist`)

```sql
CREATE SCHEMA IF NOT EXISTS juho_test;

CREATE TABLE IF NOT EXISTS juho_test.parks (
    id SERIAL PRIMARY KEY,
    osm_id BIGINT UNIQUE,
    name TEXT,
    geom GEOMETRY(Point, 4326),
    tags JSONB
);
```

### Region Not Found in `/spatial-query`

```sql
CREATE TABLE IF NOT EXISTS juho_test.boundary (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE,
    geom GEOMETRY(Polygon, 4326)
);

INSERT INTO juho_test.boundary (name, geom)
VALUES (
    'Toronto',
    ST_GeomFromText(
        'POLYGON((-79.6393 43.5810, -79.1153 43.5810, -79.1153 43.8555, -79.6393 43.8555, -79.6393 43.5810))',
        4326
    )
);
```

### Too Many Features in `osm_data.geojson`

```go
maxInsert := 1000
if len(featureCollection.Features) < maxInsert {
    maxInsert = len(featureCollection.Features)
}
```

---

## 🛠 Tech Stack

- Backend: Go (`net/http`)
- Database: PostgreSQL + PostGIS (Supabase)
- Geospatial Data: OpenStreetMap (OSM) + Overpass API

---

## 📄 License

MIT License

---

## ✨ Contributors

👤 [jh-choi98](https://github.com/jh-choi98)

---

## 🌍 References

- OpenStreetMap Overpass API: [https://overpass-api.de](https://overpass-api.de)
- PostGIS Documentation: [https://postgis.net/documentation/](https://postgis.net/documentation/)
- Supabase PostgreSQL: [https://supabase.com/docs/guides/database](https://supabase.com/docs/guides/database)
