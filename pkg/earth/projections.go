package earth

import "math"

// CalcCentreTile Convert latitude and longitude to tile coordinates at given zoom
func CalcCentreTile(lat float64, lng float64, zoom int) (int, int, int) {
	scale := int(1) << uint(zoom)
	worldCoordinateX, worldCoordinateY := project(lat, lng)
	tileX := math.Floor((worldCoordinateX * float64(scale)) / TileSize)
	tileY := math.Floor((worldCoordinateY * float64(scale)) / TileSize)
	return int(tileX), int(tileY), scale
}

// Tile2long convert tile coordinates to longitude
func Tile2long(x float64, zoom int) float64 {
	return x/math.Pow(2, float64(zoom))*360 - 180
}

// Tile2lat convert tile coordinates to longitude
func Tile2lat(y float64, zoom int) float64 {
	var n = math.Pi - 2*math.Pi*y/math.Pow(2, float64(zoom))
	return 180 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
}

// The mapping between latitude, longitude and pixels is defined by the web
// mercator projection.
func project(lat float64, lng float64) (float64, float64) {
	sinY := math.Sin((lat * math.Pi) / 180)

	// Truncating to 0.9999 effectively limits latitude to 89.189. This is
	// about a third of a tile past the edge of the world tile.
	sinY = math.Min(math.Max(sinY, -0.9999), 0.9999)
	return TileSize * (0.5 + lng/360), TileSize * (0.5 - math.Log((1+sinY)/(1-sinY))/(4*math.Pi))
}
