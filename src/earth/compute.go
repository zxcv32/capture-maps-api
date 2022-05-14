package earth

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"zxcv32/capture-maps-api/src/file"
)

const TileSize = 256
const TileScale = 2

func ComputeImages(lat float64, lng float64, zoom int) string {
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

// CalcCentreTile Convert latitude and longitude to tile coordinates at given zoom
func CalcCentreTile(lat float64, lng float64, zoom int) (int, int, int) {
	scale := int(1) << uint(zoom)
	worldCoordinateX, worldCoordinateY := project(lat, lng)
	tileX := math.Floor((worldCoordinateX * float64(scale)) / TileSize)
	tileY := math.Floor((worldCoordinateY * float64(scale)) / TileSize)
	return int(tileX), int(tileY), scale
}

func Tile2long(x float64, zoom int) float64 {
	return x/math.Pow(2, float64(zoom))*360 - 180
}

func Tile2lat(y float64, zoom int) float64 {
	var n = math.Pi - 2*math.Pi*y/math.Pow(2, float64(zoom))
	return 180 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
}
