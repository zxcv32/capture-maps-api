package earth

import "math"

// The mapping between latitude, longitude and pixels is defined by the web
// mercator projection.
func project(lat float64, lng float64) (float64, float64) {
	sinY := math.Sin((lat * math.Pi) / 180)

	// Truncating to 0.9999 effectively limits latitude to 89.189. This is
	// about a third of a tile past the edge of the world tile.
	sinY = math.Min(math.Max(sinY, -0.9999), 0.9999)
	return TileSize * (0.5 + lng/360), TileSize * (0.5 - math.Log((1+sinY)/(1-sinY))/(4*math.Pi))
}
