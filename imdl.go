package imdl

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// Download stores an image from url.
func Download(url string, fnameCh chan string, x, y uint, compress bool, m *sync.Mutex) {
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

	buf := bytes.NewBuffer(data)

	img, _, err := image.Decode(buf)
	if err != nil {
		fnameCh <- ""
		return
	}
	if x != 0 && y != 0 {
		img = resize.Resize(x, y, img, resize.Lanczos3)
	}

	var dir string
	if dir = os.Getenv("IMAGE_DIR"); dir == "" {
		dir = "."
	}

	path := fmt.Sprintf("%s/%x", dir, md5.Sum(data))
	if compress {
		path += ".jpg"
	} else {
		path += ".png"
	}
	if err = saveImage(path, img, compress, m); err != nil {
		fnameCh <- ""
		return
	}
	fnameCh <- path
}

func saveImage(path string, img image.Image, compress bool, m *sync.Mutex) error {
	m.Lock()
	defer m.Unlock()
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Compless image.
	if compress {
		err = jpeg.Encode(file, img, &jpeg.Options{jpeg.DefaultQuality})
	} else {
		err = png.Encode(file, img)
	}
	return err
}
