package earth

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"zxcv32/capture-maps-api/src/file"
)

const TileSize = 256 // Map image pixels
const TileScale = 2  // Tile pixel multiplier

// ComputeImage Retrieve the tiled image from the input latitude and longitude
func ComputeImage(lat float64, lng float64, zoom int) string {
	if err := godotenv.Load(); err != nil {
		log.Errorln("No .env file found")
	}
	key := os.Getenv("API_KEY")
	URL := fmt.Sprintf("https://maps.googleapis.com/maps/api/staticmap?center=%f,%f&key=%s&zoom=%d&size=%dx%d&scale=%d&maptype=%s&region=%s",
		lat, lng, key, zoom, TileSize, TileSize, TileScale, "hybrid", "IN")
	filename, err := file.DownloadImage(URL)
	if err != nil {
		log.Fatal(err)
	}
	return filename
}
