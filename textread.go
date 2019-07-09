package main

import "./dirread"
import (
	"log"
	"os"
	"strconv"
	"strings"
)

// listupはkeyに入力したデータがあるファイルを出力します。
func listup(str dirread.Dirtype, key string) []PlayListA {
	var name []string
	var list []PlayListA
	var tmp PlayListA
	num := 0

	for i := 0; i < len(str.Data); i++ {
		if !str.Data[i].Folder { //フォルダ以外
			if strings.Index(str.Data[i].Name, key) >= 0 { //keyワード検索
				if str.Data[i].Size > 0 { //ファイルサイズで0は無視
					name = append(name, str.Data[i].Name)
				}
			}
		}
		if i >= str.Count[num] {
			tmp.PlayData = name
			list = append(list, tmp)
			name = []string{}
			num++
		}
	}
	tmp.PlayData = name
	list = append(list, tmp)
	return list
}

// ConvertCsvはstr配列に設定したデータを出力します
//strは一次配列入力する。
//1次配列をカンマ区切りの一文に変更
func ConvertCsv(str []string, subpass string) string {
	output := ""
	for i := 0; i < len(str); i++ {
		if i == 0 {
			output = "\"" + subpass + str[0] + "\""
		} else {
			output += ",\"" + subpass + str[i] + "\""
		}

	}
	return output
}

//ListPoerser(list []PlayListA, str string, subpass string) string
//ファイル名からCSV形式のstrデータから、必要な分だけのファイル名を取り出す。
func ListPoerser(list []PlayListA, str string, subpass string) string {
	output := ""
	outtmp := ""

	if strings.Index(str, ",") >= 0 {
		tmp := strings.Split(str, ",")
		//log.Println(tmp)

		for i := 0; i < len(tmp); i++ {
			num, err := strconv.Atoi(tmp[i])
			if err != nil {
				break
			}
			editdata, numserch := listserchNum(list, num)
			if numserch == 0 {
				outtmp = "\"" + subpass + "/" + editdata + "\""
			} else {
				outtmp = "\"" + subpass + strconv.Itoa(numserch) + "/" + editdata + "\""
			}
			if output == "" {
				output = outtmp
			} else {
				output += "," + outtmp
			}
		}
	} else {
		num, err := strconv.Atoi(str)
		if err != nil {
			num = 0
		}
		editdata, numserch := listserchNum(list, num)
		if numserch == 0 {
			output = "\"" + subpass + "/" + editdata + "\""
		} else {
			output = "\"" + subpass + strconv.Itoa(numserch) + "/" + editdata + "\""
		}
	}
	return output
}

//listserchNum(list []PlayListA, num int)
//list からnum番号に対応したファイル名を返す。
func listserchNum(list []PlayListA, num int) (string, int) {
	output := ""
	outnum := 0
	count := 0
	countb := 0

	for i := 0; i < len(list); i++ {
		count += len(list[i].PlayData)
		if num < count {
			output = list[i].PlayData[num-countb]
			outnum = i
			break
		}
		countb += count
	}
	return output, outnum
}

//ReadHtml(path string)
//pathに入力したファイルパスから読み取る
//pathはテキストパスでそのテキスト値をもとに戻す
func ReadHtml(path string) string {
	var output string
	fp, err := os.Open(path)
	if err != nil {
		log.Panic(err)
		return ""
	}
	defer fp.Close()
	buf := make([]byte, 1024)
	for {
		n, err := fp.Read(buf)
		if err != nil {
			break
		}
		if n == 0 {
			break
		}
		output += string(buf[:n])
	}
	return output
}

//ConvertDataはstrに入力されたデータから<%%>に囲まれた文字列から
//data[文字列]の入力された値と置き換えて変換値を戻り値にする
func ConvertData(str string, data map[string]string) string {
	//func ConvertData(str string) string {
	tmp := str
	output := str
	for {
		n := strings.Index(tmp, "<%")
		if n >= 0 {
			m := strings.Index(tmp, "%>")
			if m >= 0 {
				dtmp := tmp[n+2 : m]
				output = strings.Replace(output, "<%"+dtmp+"%>", data[dtmp], 1)
				tmp = tmp[m+2:]
			} else {
				break
			}
		} else {
			break
		}
	}
	return output
}
