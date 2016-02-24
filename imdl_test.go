package imdl

import (
	"os"
	"sync"
	"testing"
)

var URLs = []string{
	"http://www.libsdl.org/projects/SDL_image/docs/demos/lena.jpg",
	"http://optipng.sourceforge.net/pngtech/img/lena.png",
}

func TestStoreImage(t *testing.T) {
	os.Setenv("IMAGE_DIR", "images")
	c := make(chan string, len(URLs))
	var m sync.Mutex
	go Download(URLs[0], c, 480, 360, true, &m)
	go Download(URLs[0], c, 480, 360, false, &m)

	for i := 0; i < len(URLs); i++ {
		fname := <-c
		t.Log(fname)
	}
}
