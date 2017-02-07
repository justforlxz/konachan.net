package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/widuu/gojson"
	// simplejson "github.com/bitly/go-simplejson"
)

var (
	url = "http://konachan.net/post.json?page="

	programPath = getCurrentDirectory() + "/id.txt"
	ID          string
	state       bool = false
)

func main() {

	file, _ := os.Open(programPath)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println("上次执行到的id为: " + scanner.Text())
		ID = scanner.Text()
	}

	res, _ := http.Get(url + "1")
	byte, _ := ioutil.ReadAll(res.Body)
	// json, _ := simplejson.NewJson(byte)

	// fmt.Println(string(byte))

	c1 := gojson.Json(string(byte))
	fmt.Println(c1)

	i := 0
	for {

		// fmt.Println(c1.GetIndex(i).ToArray())

		// id, _ := json.GetIndex(i).Get("id").String()

		// if id == "" {
		// 	break
		// }

		// fmt.Println(id)

		// if id == ID {
		// 	fmt.Println("数据已最新")
		// 	return
		// }
		// if state == false {
		// 	dstFile, err := os.Create(programPath)
		// 	if err != nil {
		// 		fmt.Println(err.Error())
		// 		return
		// 	}
		// 	defer dstFile.Close()
		// 	dstFile.WriteString(id)
		// 	state = true
		// }
		// data, _ := json.GetIndex(i).Get("file_url").String()
		// if data == "" {
		// 	break
		// }
		i++
		// fmt.Println(data)
	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
