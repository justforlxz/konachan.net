package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	url = "http://konachan.net/post.json?page="

	programPath  = getCurrentDirectory() + "/id.txt"
	ID           int64
	state        bool = false
	has_children bool
	id           int64
	data         string
	ok           bool
)

type image struct {
	ID          int64  `json:"id"`
	FileURL     string `json:"file_url"`
	HasChildren bool   `json:"has_children"`
}

func main() {

	file, _ := os.Open(programPath)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println("上次执行到的id为: " + scanner.Text())
		ID, _ = strconv.ParseInt(scanner.Text(), 10, 64)
	}

	var index int64 = 1

	err := os.Remove(getCurrentDirectory() + "/file.txt")

	if err != nil {
		fmt.Println(err.Error())
	}

	urlFile, err := os.OpenFile(getCurrentDirectory()+"/file.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	p := []byte("[]")
	for {

		res, err := http.Get(url + strconv.FormatInt(index, 10))

		if err != nil {
			fmt.Println("Get error:", err)
			return
		}

		byte, err := ioutil.ReadAll(res.Body)

		if err != nil {
			fmt.Println("Reading error:", err)
			return
		}

		imageArray := make([]image, 0)

		if err = json.Unmarshal(byte, &imageArray); err != nil {
			fmt.Println("Unmarshal error: ", err)
		}

		if bytes.Equal(byte, p) {
			fmt.Println("已经到网站最后一页，首次爬虫执行完毕，开始下载...")
			getPic()
			return
		}

		for _, source := range imageArray {

			if source.ID == ID {
				fmt.Println("数据已最新")
				getPic()
				return
			}

			if source.HasChildren {
				fmt.Println("正在跳过R18...")
				continue
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

			data = "http:" + data
			urlFile.WriteString(data + "\n")
			fmt.Println("正在记录文件: " + data)
		}
		index++
	}
	defer urlFile.Close()
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func getPic() {
	cmd := exec.Command("wget", "-i", getCurrentDirectory()+"/file.txt")
	//显示运行的命令
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	cmd.Wait()
}
