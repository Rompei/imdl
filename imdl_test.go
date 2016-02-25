package imdl

import (
	"os"
	"runtime"
	"sync"
	"testing"
)

var URLs = []string{
	"http://img.cupo.cc/wp-content/uploads/comment/2015/09/159755f270a3492d7.jpg",
	"https://i.ytimg.com/vi/e12aBCICrl4/maxresdefault.jpg",
	"http://blog-imgs-43.fc2.com/t/n/7/tn777dome/20141004232349bf4.jpg",
	"http://blog-imgs-65.fc2.com/h/a/t/hatsunemiku1006/20140624193649b83.jpg",
	"http://blog-imgs-65.fc2.com/a/r/a/arakawapoke/20140705154411872.jpg",
	"http://livedoor.4.blogimg.jp/anico_bin/imgs/1/7/1700d7e7.jpg",
	"http://articleimage.nicoblomaga.jp/image/56/2015/5/b/5b43c22afefec2d2c74dc4f32a690f127de043211426906261.jpg",
	"http://bpmaker-resource.giffy.me/userdata/user/38/38342/5/b1-1404387709.jpg",
	"http://livedoor.blogimg.jp/rytescarlet/imgs/4/1/411f8cb2.jpg",
}

func TestStoreImage(t *testing.T) {
	os.Setenv("IMAGE_DIR", "images")
	c := make(chan string, len(URLs))
	errCh := make(chan error, runtime.NumCPU())
	var m sync.Mutex
	for _, u := range URLs {
		go Download(u, c, errCh, 680, 480, false, &m)
	}

	for i := 0; i < len(URLs); i++ {
		fname := <-c
		t.Log(fname)
	}
}

func TestStoreImageCompress(t *testing.T) {
	os.Setenv("IMAGE_DIR", "images")
	c := make(chan string, len(URLs))
	errCh := make(chan error, runtime.NumCPU())
	var m sync.Mutex
	for _, u := range URLs {
		go Download(u, c, errCh, 680, 480, true, &m)
	}

	for i := 0; i < len(URLs); i++ {
		fname := <-c
		t.Log(fname)
	}
}
