package imdl

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// Download stores an image from url.
func Download(url string, fnameCh chan string) {
	ext := filepath.Ext(url)
	if ext == "" {
		fnameCh <- ""
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		fnameCh <- ""
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fnameCh <- ""
		return
	}

	var dir string
	if dir = os.Getenv("IMAGE_DIR"); dir == "" {
		dir = "."
	}
	fname := fmt.Sprintf("%x%s", md5.Sum(data), ext)

	file, err := os.Create(dir + "/" + fname)
	if err != nil {
		fnameCh <- ""
		return
	}
	defer file.Close()

	file.Write(data)
	fnameCh <- fname
}
