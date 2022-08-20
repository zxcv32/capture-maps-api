package file

import (
	gim "github.com/ozankasikci/go-image-merge"
	log "github.com/sirupsen/logrus"
	"image/png"
	"os"
)

func Merge(grids []*gim.Grid, x int, y int) string {
	rgba, err := gim.New(grids, x, y).Merge()
	if err != nil {
		log.Errorln(err)
	}
	filename := generateFilename()
	file, err := os.Create(filename)
	err = png.Encode(file, rgba)
	return filename
}
