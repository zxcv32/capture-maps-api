package file

import (
	"errors"
	"io"
	"net/http"
	"os"
)

//DownloadImage The function downloads the image from the given URL
func DownloadImage(URL string) (string, error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", errors.New("received non 200 response code from Google maps static API")
	}
	// Create an empty file
	file, err := os.Create(GetFilename())
	if err != nil {
		return "", err
	}
	defer file.Close()
	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}
