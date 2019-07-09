package jsonread

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"log"
)

// SetupdataはJosnファイルの読み込み設定
type Setupdata struct {
	Server    string   `json:"server"`    //server設定
	Port      string   `json:"port"`      //Port設定
	Path      string   `json:path`        //path設定
	Videopath []string `json:"videopath"` //videoパスの設定
}

//Jsondata はクラス
type Jsondata struct {
	name string //読み込みファイルネーム
}

//(t *Jsondata) Setup (s string)は読み込みファイルネーム設定 sに入力したファイルをJsondata.nameに設定
func (t *Jsondata) Setup(s string) {
	t.name = s
}

func (t Jsondata) Read(data *Setupdata) int {
	if t.name == "" {
		return -1
	}
	byte, err := ioutil.ReadFile(t.name)
	if err != nil {
		log.Fatal(err)
	}
	var persons Setupdata
	if err := json.Unmarshal(byte, &persons); err != nil {
		log.Fatal(err)
	}

	persons.Path = t.ckpath(persons.Path)
	for i := 0; i < len(persons.Videopath); i++ {
		persons.Videopath[i] = t.ckpath(persons.Videopath[i])
	}

	*data = persons
	//fmt.Println(persons)
	return 0
}

//(t Jsondata) ckpath(s string)　パス確認をして最後に/が抜けている場合は追加する
func (t Jsondata) ckpath(s string) string {
	tmp := s
	if tmp != "" {
		if tmp[len(tmp)-1:] == "/" {
		} else {
			tmp += "/"
		}
	}
	return tmp
}
