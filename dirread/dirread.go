package dirread

import (
	"io/ioutil"
	"log"
	"time"
)

type filedata struct {
	Name     string
	Folder   bool
	Size     int64
	Time     time.Time
	RootPath string
}

// Dirtype implements buffering for an []filedata object.
//Dirtypeは[]filedata objectをバッファする必要あり
type Dirtype struct {
	path  string
	Data  []filedata
	Count []int
	Renew bool
}

func (t *Dirtype) Setup(s string) {
	t.path = s
	var tmp []filedata
	var tmp2 []int
	if (len(t.Data) == 0) || (t.Renew) {
		t.Data = tmp
		t.Count = tmp2
		t.Renew = false
	}
}

func (t *Dirtype) Read() int {
	var tmp []filedata
	if t.path == "" {
		return -1
	}
	files, err := ioutil.ReadDir(t.path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
<<<<<<< HEAD
		tmp2 := filedata{s + f.Name(), f.IsDir(), f.Size(), f.ModTime(), t.path + s}
=======
		tmp2 := filedata{f.Name(), f.IsDir(), f.Size()}
>>>>>>> parent of c9276db... ok data
		tmp = append(tmp, tmp2)
	}
	t.Data = tmp
	return 0

}
