package main

import (
	"bufio"
	"bytes"
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

	simplejson "github.com/bitly/go-simplejson"
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
		byte, err := ioutil.ReadAll(res.Body)
		json, err := simplejson.NewJson(byte)

		if err != nil {
			fmt.Println("我也不知道发生了什么错误，反正没请求到数据...")
			break
		}

		if bytes.Equal(byte, p) {
			fmt.Println("已经到网站最后一页，首次爬虫执行完毕，开始下载...")
			getPic()
			return
		}

		i := 0

		for {

			has_children, err := json.GetIndex(i).Get("has_children").Bool()

			if err != nil {
				fmt.Println("获取R18标志失败...")
				break
			}

			if !has_children {
				fmt.Println("跳过R18...")
				break
			}

			id, idErr := json.GetIndex(i).Get("id").Int64()

			if idErr != nil {
				fmt.Println("获取Id失败...")
				break
			}

			if id == 0 {
				break
			}

			if id == ID {
				fmt.Println("数据已最新")
				getPic()
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
			data, dataErr := json.GetIndex(i).Get("file_url").String()

			if dataErr != nil {
				fmt.Println("获取data失败...")
				break
			}

			if data == "" {
				break
			}
			data = "http:" + data
			urlFile.WriteString(data + "\n")
			fmt.Println("正在记录文件: " + data)
			i++
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
