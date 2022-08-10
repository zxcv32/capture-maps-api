package file

import (
	"errors"
	gim "github.com/ozankasikci/go-image-merge"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const EXT = ".png"

func GetFilename() string {
	var filename string
	for true {
		filename = generateFilename()
		if !exists(filename) {
			break
		}
	}
	return filename
}

func generateFilename() string {
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(rand.Int()) + EXT
	return fileName
}

func exists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		if err != nil {
			log.Errorln(err)
		}
		log.Errorf("Unknown error occurred while verifying if the file '%s' exists", filename)
	}
	return false
}

func DeleteFiles(grids []*gim.Grid) {
	for i := 0; i < len(grids); i++ {
		file := grids[i].ImageFilePath
		DeleteFile(file)
	}
}
func DeleteFile(file string) {
	var e = os.Remove(file)
	if e != nil {
		log.Errorln(e)
	}
}
