package utils

import (
	"os"
    "net/http"
	"io"
	"fmt"
)

//DownloadFile download a file from url into filepath 
func DownloadFile(filepath string, url string) (error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil  {
		return fmt.Errorf("%s could not be created: %s", filepath, err)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("URL %s could not be retrieved: %s", url, err)
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return fmt.Errorf("Error while writing the content: %s", err)
	}

	return nil
}