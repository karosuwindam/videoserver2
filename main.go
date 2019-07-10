package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"./dirread"
	"./jsonread"
)

//LISTPATH リストのHTML
const LISTPATH = "list.html" //list.html
//PLAYPATH 再生用ののHTMLテンプレート
const PLAYPATH = "play.html" //play.html
//JSONPATH 設定ファイルのパス
const JSONPATH = "config.json"

//動画パスフォルダ指定
const DATAPATH = "video"

//Datamap <%%>で囲まれたデータを変換させるデータ
var Datamap map[string]string

//PlaylistData Dirファイル結果からmp4だけ検出結果入力
type PlayListA struct {
	PlayData []string
}

var PlaylistData []PlayListA

//Cdata 設定ファイルをJSONで読み取る結果
var Cdata jsonread.Setupdata

//deflistcreat PlaylistDataを作るプログラム
func deflistcreat() {
	dirClass := new(dirread.Dirtype)
	countmaxb := 0
	countmax := 0
	//Todo:
	for j := 0; j < len(Cdata.Videopath); j++ {
		//for j := 0; j < 1; j++ {
		dirClass.Setup(Cdata.Videopath[j])
		dirClass.Read("")
		countmax = len(dirClass.Data)
		for i := countmaxb; i < countmax; i++ {
			if dirClass.Data[i].Folder {
				tmp := dirClass.Data[i].Name
				dirClass.Read(tmp + "/")
			}
		}
		for {
			countmaxb = countmax
			if countmaxb < len(dirClass.Data) {
				countmax = len(dirClass.Data)
				for i := countmaxb; i < countmax; i++ {
					if dirClass.Data[i].Folder {
						tmp := dirClass.Data[i].Name
						dirClass.Read(tmp + "/")
					}
				}
			} else {
				break
			}
		}
		//countTmp = append(bb, countmaxb)
		dirClass.Count = append(dirClass.Count, countmaxb)
	}
	//	dirClass.Count = countTmp
	PlaylistData = listup(*dirClass, ".mp4")
	//fmt.Println(PlaylistData)
}

//main実行
func main() {
	tempMap := map[string]string{}
	Datamap = tempMap

	jsonClassR := new(jsonread.Jsondata)
	//dirClass := new(dirread.Dirtype)

	jsonClassR.Setup("./" + JSONPATH)
	jsonClassR.Read(&Cdata)
	//dirClass.Setup(Cdata.Videopath)
	//dirClass.Read()

	//mp4ファイルを取得
	//PlaylistData = listup(*dirClass, ".mp4")
	deflistcreat()
	for i := 0; i < len(PlaylistData); i++ {
		if i == 0 {
			Datamap["namedata"] = ConvertCsv(PlaylistData[i].PlayData, DATAPATH+"/")
		} else {
			num := strconv.Itoa(i)
			tmp := ConvertCsv(PlaylistData[i].PlayData, DATAPATH+num+"/")
			if tmp != "" {
				Datamap["namedata"] += "," + tmp
			}
		}
	}

	//http setup
	http.HandleFunc("/play", playnow)
	http.HandleFunc("/list", playlist)
	//ToDo:
	for i := 0; i < len(PlaylistData); i++ {
		if i == 0 { //フォルダ登録
			vfp := http.StripPrefix("/"+DATAPATH+"/", http.FileServer(http.Dir(Cdata.Videopath[i])))
			http.Handle("/"+DATAPATH+"/", vfp)
		} else {
			num := strconv.Itoa(i)
			http.Handle("/"+DATAPATH+num+"/", http.StripPrefix("/"+DATAPATH+num+"/", http.FileServer(http.Dir(Cdata.Videopath[i]))))
		}
	}
	fp2 := http.StripPrefix("/", http.FileServer(http.Dir(Cdata.Path)))
	http.Handle("/", fp2)
	fmt.Printf("Server Start\nServer IP:%s\nPort:%s\n", Cdata.Server, Cdata.Port)
	http.ListenAndServe(Cdata.Server+":"+Cdata.Port, nil)
}

//MapCGid ?以降のデータを解析

func MapCGIparser(s string) map[string]string {
	output := map[string]string{}
	if strings.Index(s, "?") >= 0 {
		tmp := strings.Split(s, "?")[1]
		if strings.Index(tmp, "=") >= 0 {
			if strings.Index(tmp, "&") >= 0 {
				tmpd := strings.Split(tmp, "&")
				for i := 0; i < len(tmpd); i++ {
					stmp := strings.Split(tmpd[i], "=")
					output[stmp[0]] = stmp[1]
					Datamap[stmp[0]] = stmp[1]
				}
			} else {
				//fmt.Println(tmp)
				stmp := strings.Split(tmp, "=")
				fmt.Println(stmp)
				output[stmp[0]] = stmp[1]
				Datamap[stmp[0]] = stmp[1]

			}
		}
	} else {

	}
	return output

}

//CGI playlistを作成する
func playnow(w http.ResponseWriter, r *http.Request) {
	//str := strings.Split(r.RequestURI, "?")[1]
	if r.Method == "POST" {
		//Datamap["namedata"] = r.FormValue("list")
		if r.FormValue("list") != "" {
			Datamap["namedata"] = ListPoerser(PlaylistData, r.FormValue("list"), DATAPATH)
		} else {
			for i := 0; i < len(PlaylistData); i++ {
				if i == 0 {
					Datamap["namedata"] = ConvertCsv(PlaylistData[i].PlayData, DATAPATH+"/")
				} else {
					num := strconv.Itoa(i)
					tmp := ConvertCsv(PlaylistData[i].PlayData, DATAPATH+num+"/")
					if tmp != "" {
						Datamap["namedata"] += "," + tmp
					}
				}
			}
		}

		Datamap["id"] = "0"
		fmt.Println(Datamap["namedata"])
		//fmt.Println(r.FormValue("list") == "")
	} else {
		for i := 0; i < len(PlaylistData); i++ {
			if i == 0 {
				Datamap["namedata"] = ConvertCsv(PlaylistData[i].PlayData, DATAPATH+"/")
			} else {
				num := strconv.Itoa(i)
				tmp := ConvertCsv(PlaylistData[i].PlayData, DATAPATH+num+"/")
				if tmp != "" {
					Datamap["namedata"] += "," + tmp
				}
			}
		}
	}
	fmt.Println(MapCGIparser(r.RequestURI))
	//fmt.Println(MapCGIparser("gggb?id=ggg"))

	output := ConvertData(ReadHtml(Cdata.Path+PLAYPATH), Datamap)
	if Datamap["namedata"] != "" {
		fmt.Fprintf(w, output)
	} else {
		fmt.Fprintf(w, "No mp4 data")
	}
}

//CGI Playlistを表示する
func playlist(w http.ResponseWriter, r *http.Request) {
	var fp *os.File
	var err error
	var tmp string

	deflistcreat()
	fmt.Println(r)
	fp, err = os.Open(Cdata.Path + LISTPATH)
	if err != nil {

	}
	reader := bufio.NewReaderSize(fp, 4096)
	for line := ""; err == nil; line, err = reader.ReadString('\n') {
		if strings.Index(line, "<%output%>") >= 0 {
			//line = cgiEditOutput()
			line = strings.Replace(line, "<%output%>", cgiEditOutput(), 1)

		}
		tmp += line
		//fmt.Print(line)
	}
	if err != io.EOF {
		panic(err)
	}
	//tmp := r.URL.RawQuery
	fmt.Fprintf(w, "%s", tmp)

}

//リストをHTML出力
func cgiEditOutput() string {
	var output string
	count := 0
	for j := 0; j < len(PlaylistData); j++ {
		for i := 0; i < len(PlaylistData[j].PlayData); i++ {
			tmpdata := PlaylistData[j].PlayData[i]
			if strings.Index(tmpdata, ".mp4") >= 0 {
				ai := strconv.Itoa(count)
				output += "<div>"
				output += "<input type=\"checkbox\" name=\"list\" id=\"ck" + ai + "\">"
				output += "<a href=\"play?id=" + ai + "\">" + tmpdata + "</a></div>\n"
			} else {
				//output += "<div>" + tmpdata + "<br><li>" + "</li></div>\n"
			}
			count++

		}
	}
	return output
}
