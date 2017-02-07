package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	// simplejson "github.com/bitly/go-simplejson"
)

var (
	url = "http://konachan.net/post.json?page="

	programPath = getCurrentDirectory() + "/id.txt"
	ID          int64
	state       bool = false
)

func main() {

	file, _ := os.Open(programPath)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println("上次执行到的id为: " + scanner.Text())
		ID, _ = strconv.ParseInt(scanner.Text(), 10, 64)
	}

	res, _ := http.Get(url + "1")
	byte, _ := ioutil.ReadAll(res.Body)
	json, _ := simplejson.NewJson(byte)

	// fmt.Println(string(byte))

	i := 0
	for {

		id, _ := json.GetIndex(i).Get("id").Int64()
		if id == 0 {
			break
		}

		if id == ID {
			fmt.Println("数据已最新")
			return
		}

		if state == false {
			dstFile, err := os.Create(programPath)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer dstFile.Close()
			dstFile.WriteString(strconv.FormatInt(id, 10))
			state = true
		}
		data, _ := json.GetIndex(i).Get("file_url").String()
		if data == "" {
			break
		}
		i++
		fmt.Println("http:" + data)
	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
