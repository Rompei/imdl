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
func Download(url string, fnameCh chan string, errCh chan error, x, y uint, compress bool, m *sync.Mutex) {
	ext := filepath.Ext(url)
	if ext == "" {
		errCh <- fmt.Errorf("Extention was not detected.")
		return
	}
	img, data, err := getImage(url)
	if err != nil {
		errCh <- err
		return
	}
	if x != 0 && y != 0 {
		img = resize.Resize(x, y, img, resize.Lanczos3)
	}

	var dir string
	if dir = os.Getenv("IMAGE_DIR"); dir == "" {
		dir = "."
	}

	fname := fmt.Sprintf("%x", md5.Sum(data))
	path := fmt.Sprintf("%s/%s", dir, fname)
	if compress {
		fname += ".jpg"
		path += ".jpg"
	} else {
		fname += ".png"
		path += ".png"
	}
	if err = saveImage(path, img, compress, m); err != nil {
		errCh <- err
		return
	}
	fnameCh <- fname
}

// DownloadNorm is normal download.
func DownloadNorm(url string, x, y uint, compress bool, m *sync.Mutex) (string, error) {
	ext := filepath.Ext(url)
	if ext == "" {
		return "", fmt.Errorf("Extention was not detected.")
	}
	img, data, err := getImage(url)
	if err != nil {
		return "", err
	}
	if x != 0 && y != 0 {
		img = resize.Resize(x, y, img, resize.Lanczos3)
	}

	var dir string
	if dir = os.Getenv("IMAGE_DIR"); dir == "" {
		dir = "."
	}

	fname := fmt.Sprintf("%x", md5.Sum(data))
	path := fmt.Sprintf("%s/%s", dir, fname)
	if compress {
		fname += ".jpg"
		path += ".jpg"
	} else {
		fname += ".png"
		path += ".png"
	}
	if err = saveImage(path, img, compress, m); err != nil {
		return "", err
	}
	return fname, nil
}

// DownloadToPath downloads image to specified path.
func DownloadToPath(url, dir string, fnameCh chan string, errCh chan error, x, y uint, compress bool, m *sync.Mutex) {
	ext := filepath.Ext(url)
	if ext == "" {
		errCh <- fmt.Errorf("Extention was not detected.")
		return
	}
	img, data, err := getImage(url)
	if err != nil {
		errCh <- err
		return
	}
	if x != 0 && y != 0 {
		img = resize.Resize(x, y, img, resize.Lanczos3)
	}

	if dir == "" {
		dir = "."
	}

	if !fileIsExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			errCh <- err
			return
		}
	}

	fname := fmt.Sprintf("%x", md5.Sum(data))
	path := fmt.Sprintf("%s/%s", dir, fname)
	if compress {
		fname += ".jpg"
		path += ".jpg"
	} else {
		fname += ".png"
		path += ".png"
	}
	if err = saveImage(path, img, compress, m); err != nil {
		errCh <- err
		return
	}
	fnameCh <- fname
}

// DownloadToPathNoem download normally to directory.
func DownloadToPathNoem(url, dir string, x, y uint, compress bool, m *sync.Mutex) (string, error) {
	ext := filepath.Ext(url)
	if ext == "" {
		return "", fmt.Errorf("Extention was not detected.")
	}
	img, data, err := getImage(url)
	if err != nil {
		return "", err
	}
	if x != 0 && y != 0 {
		img = resize.Resize(x, y, img, resize.Lanczos3)
	}

	if dir == "" {
		dir = "."
	}

	if !fileIsExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}

	fname := fmt.Sprintf("%x", md5.Sum(data))
	path := fmt.Sprintf("%s/%s", dir, fname)
	if compress {
		fname += ".jpg"
		path += ".jpg"
	} else {
		fname += ".png"
		path += ".png"
	}
	if err = saveImage(path, img, compress, m); err != nil {
		return "", err
	}
	return fname, nil
}

func getImage(url string) (image.Image, []byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	buf := bytes.NewBuffer(data)

	img, _, err := image.Decode(buf)
	if err != nil {
		return nil, nil, err
	}
	return img, data, nil
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

func fileIsExist(fname string) bool {
	_, err := os.Stat(fname)
	return err == nil
}
