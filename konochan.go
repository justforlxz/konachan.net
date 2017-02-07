package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
)

var (
	url = "http://konachan.net/post.json?page="
)

func main() {

	res, _ := http.Get(url + "1")
	byte, _ := ioutil.ReadAll(res.Body)
	json, _ := simplejson.NewJson(byte)

	i := 0
	for {
		data, _ := json.GetIndex(i).Get("file_url").String()
		if data == "" {
			break
		}
		i++
		fmt.Println(data)
	}
}
