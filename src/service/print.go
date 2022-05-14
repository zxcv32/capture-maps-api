package service

import (
	"encoding/json"
	gim "github.com/ozankasikci/go-image-merge"
	"io/ioutil"
	"net/http"
	"zxcv32/capture-maps-api/src/earth"
	"zxcv32/capture-maps-api/src/file"
)

type printRequest struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Zoom   int     `json:"zoom"`
	Radius int     `json:"radius"`
}

// HandleRequest POST request
func HandleRequest(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var task printRequest
	err := decoder.Decode(&task)
	if err != nil {
		panic(err)
	}

	lat, lng, zoom, radius := task.Lat, task.Lng, task.Zoom, task.Radius
	filename := captureTiles(lat, lng, zoom, radius)
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/octet-stream")
	_, error := writer.Write(fileBytes)
	if error != nil {
		http.Error(writer, error.Error(), 500)
		return
	}
	file.DeleteFile(filename)
}

func captureTiles(lat float64, lng float64, zoom int, radius int) string {
	centreTileX, centreTileY, scale := earth.CalcCentreTile(lat, lng, zoom)
	var grids []*gim.Grid
	latitudinalTiles := min(radius*2-1, scale)
	var xIndex = centreTileX - (radius - 1)
	lastXIndexTile := latitudinalTiles + xIndex
	var j = centreTileY
	var gimGridCentre gim.Grid
	gimGridCentre.ImageFilePath = westToEast(xIndex, j, lastXIndexTile, zoom)

	j = centreTileY - radius
	repeat := radius - 1
	for repeat > 0 {
		repeat--
		j++ // go from north to south towards centre
		var gimGridNorth gim.Grid
		gimGridNorth.ImageFilePath = westToEast(xIndex, j, lastXIndexTile, zoom)
		grids = append(grids, &gimGridNorth)
	}

	grids = append(grids, &gimGridCentre)

	j = centreTileY // reset j
	repeat = radius - 1
	for repeat > 0 {
		repeat--
		j++ // go south
		var gimGridSouth gim.Grid
		gimGridSouth.ImageFilePath = westToEast(xIndex, j, lastXIndexTile, zoom)
		grids = append(grids, &gimGridSouth)
	}

	flashbang := file.Merge(grids, 1, radius*2-1)
	file.DeleteFiles(grids)
	return flashbang
}

func westToEast(xIndex int, j int, lastXIndexTile int, zoom int) string {
	var grids []*gim.Grid
	// West to east loop
	for i := xIndex; i < lastXIndexTile; i++ {
		tileLat := earth.Tile2lat(float64(j), zoom)
		tileLng := earth.Tile2long(float64(i), zoom)
		filename := earth.ComputeImages(tileLat, tileLng, zoom)
		var gimGrid gim.Grid
		gimGrid.ImageFilePath = filename
		grids = append(grids, &gimGrid)
	}
	merged := file.Merge(grids, len(grids), 1)
	file.DeleteFiles(grids)
	return merged
}

func min(a int, b int) int {
	min := a
	if b < min {
		min = b
	}
	return min
}
