package earth

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zxcv32/capture-maps-api/pkg/file"
)

const TileSize = 256 // Map image pixels
const TileScale = 2  // Tile pixel multiplier

// ComputeImage Retrieve the tiled image from the input latitude and longitude
func ComputeImage(lat float64, lng float64, zoom int, mapTypeId string, apiKey string) (string, error) {
	URL := fmt.Sprintf("https://maps.googleapis.com/maps/api/staticmap?center=%f,%f&key=%s&zoom=%d&size=%dx%d&scale=%d&maptype=%s&region=%s",
		lat, lng, apiKey, zoom, TileSize, TileSize, TileScale, mapTypeId, "IN")
	filename, err := file.DownloadImage(URL)
	if err != nil {
		log.Errorln(err)
		return "", err
	}
	return filename, nil
}
