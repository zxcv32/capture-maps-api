package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"

	gim "github.com/ozankasikci/go-image-merge"
	log "github.com/sirupsen/logrus"
	"github.com/zxcv32/capture-maps-api/pkg/earth"
	"github.com/zxcv32/capture-maps-api/pkg/file"
)

// printRequest struct to accept print request
type printRequest struct {
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Zoom      int     `json:"zoom"`
	Radius    int     `json:"radius"`
	MapTypeId string  `json:"mapTypeId"`
}

// setupResponse setups all the common HTTP response headers
func setupResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Content-Type", "image/png")
}

// HandleRequest print request
func PrintHandler(apiKey string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		setupResponse(&writer)
		if request.Method == "OPTIONS" {
			return
		}
		var requestId = rand.Int()
		log.Printf("Print request received: %d", requestId)
		body := request.Body
		decoder := json.NewDecoder(body)
		var task printRequest
		err := decoder.Decode(&task)
		if err != nil {
			log.Errorln(err)
			writer.WriteHeader(500)
			log.Printf("Print request not complete: %d", requestId)
			return
		}

		lat, lng, zoom, radius, mapTypeId := task.Lat, task.Lng, task.Zoom, task.Radius, task.MapTypeId
		validationError := validateRequest(lat, lng, zoom, radius, mapTypeId)
		if validationError != nil {
			http.Error(writer, validationError.Error(), 400)
			return
		}
		filename, err := captureTiles(lat, lng, zoom, radius, mapTypeId, apiKey)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}
		fileBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		writer.WriteHeader(http.StatusOK)
		_, error := writer.Write(fileBytes)
		if error != nil {
			http.Error(writer, error.Error(), 500)
			return
		}
		file.DeleteFile(filename)
		log.Printf("Print request complete: %d", requestId)
	}
}

func validateRequest(lat float64, lng float64, zoom int, radius int, mapTypeId string) error {
	if lat < -90 || lat > 90 {
		return errors.New("invalid latitude. Possible set: [-90,90]")
	}
	if lng < -180 || lng > 180 {
		return errors.New("invalid longitude. Possible set: [-180,80]")
	}
	if zoom < 0 || zoom > 21 {
		return errors.New("invalid zoom. Possible set: [0,21]")
	}
	if radius < 1 || radius > 15 {
		return errors.New("invalid radius. Possible set: [1,15]")
	}
	switch mapTypeId {
	case
		"hybrid",
		"satellite",
		"roadmap",
		"terrain":
	default:
		return errors.New("invalid mapTypeId. Possible values: hybrid,satellite,roadmap,terrain")
	}
	return nil
}

func captureTiles(lat float64, lng float64, zoom int, radius int, mapTypeId string, apiKey string) (string, error) {
	centreTileX, centreTileY, scale := earth.CalcCentreTile(lat, lng, zoom)
	var grids []*gim.Grid
	latitudinalTiles := min(radius*2-1, scale)
	var xIndex = centreTileX - (radius - 1)
	lastXIndexTile := latitudinalTiles + xIndex
	var j = centreTileY
	var gimGridCentre gim.Grid
	path, err := westToEast(xIndex, j, lastXIndexTile, zoom, mapTypeId, apiKey)
	if err != nil {
		return "", err
	}
	gimGridCentre.ImageFilePath = path
	j = centreTileY - radius
	repeat := radius - 1
	for repeat > 0 {
		repeat--
		j++ // go from north to south towards centre
		var gimGridNorth gim.Grid
		path, err := westToEast(xIndex, j, lastXIndexTile, zoom, mapTypeId, apiKey)
		if err != nil {
			return "", err
		}
		gimGridNorth.ImageFilePath = path
		grids = append(grids, &gimGridNorth)
	}

	grids = append(grids, &gimGridCentre)

	j = centreTileY // reset j
	repeat = radius - 1
	for repeat > 0 {
		repeat--
		j++ // go south
		var gimGridSouth gim.Grid
		path, err := westToEast(xIndex, j, lastXIndexTile, zoom, mapTypeId, apiKey)
		if err != nil {
			return "", err
		}
		gimGridSouth.ImageFilePath = path
		grids = append(grids, &gimGridSouth)
	}

	flashbang := file.Merge(grids, 1, radius*2-1)
	file.DeleteFiles(grids)
	return flashbang, nil
}

func westToEast(xIndex int, j int, lastXIndexTile int, zoom int, mapTypeId string, apiKey string) (string, error) {
	var grids []*gim.Grid
	// West to east loop
	for i := xIndex; i < lastXIndexTile; i++ {
		tileLat := earth.Tile2lat(float64(j), zoom)
		tileLng := earth.Tile2long(float64(i), zoom)
		filename, err := earth.ComputeImage(tileLat, tileLng, zoom, mapTypeId, apiKey)
		if err != nil {
			return "", err
		}
		var gimGrid gim.Grid
		gimGrid.ImageFilePath = filename
		grids = append(grids, &gimGrid)
	}
	merged := file.Merge(grids, len(grids), 1)
	file.DeleteFiles(grids)
	return merged, nil
}

func min(a int, b int) int {
	min := a
	if b < min {
		min = b
	}
	return min
}
