package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const API string = "http://fanyi.youdao.com/openapi.do?keyfrom=go-translator&key=307165215&type=data&doctype=json&version=1.1&q="

type Web struct {
	Value []string `json:"value"`
	Key   string   `json:"key"`
}

type Basic struct {
	Phonetic string   `json:"phonetic"`
	Explains []string `json:"explains"`
}

type Translation struct {
	Translation []string `json:"translation"`
	Basic       Basic    `json:"basic"`
	Query       string   `json:"query"`
	ErrorCode   float64  `json:"errorCode"`
	Web         []Web    `json:"web"`
}

func main() {
	keys := make(map[int]bool)

	for k, _ := range os.Args {
		keys[k] = true
	}

	if _, ok := keys[1]; !ok {
		fmt.Println("USEAGE: ")
		fmt.Println("\t- $ goyd 词/句 [是否读出来]")
		fmt.Println("\t- $ goyd golang 1 /* 试一下 */")
		return
	}

	input := os.Args[1]
	client := http.Client{
		Timeout: time.Duration(time.Second * 5),
	}
	resp, err := client.Get(API + input)

	if err != nil {
		fmt.Println("出错啦：网络不稳定啊少年，-1s")
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("出错啦：有道翻译好像出问题了，-1s")
	}

	var j Translation
	err = json.Unmarshal(data, &j)

	if err != nil {
		fmt.Println("出错啦：难道有道已经停用了，-1s")
	}

	if code := j.ErrorCode; code > 0 {
		//errorCode：
		//　0 - 正常
		//　20 - 要翻译的文本过长
		//　30 - 无法进行有效的翻译
		//　40 - 不支持的语言类型
		//　50 - 无效的key
		//　60 - 无词典结果，仅在获取词典结果生效
		switch code {
		case 20:
			fmt.Println("出错啦：要翻译的文本过长")
		case 30:
			fmt.Println("出错啦：无法进行有效的翻译")
		case 40:
			fmt.Println("出错啦：不支持的语言类型")
		case 50:
			fmt.Println("出错啦：无效的key")
		case 60:
			fmt.Println("出错啦：无词典结果，仅在获取词典结果生效")
		}

		return
	}

	fmt.Printf("翻译：\t%s\n", strings.Join(j.Translation[:], " / "))

	if phonetic := j.Basic.Phonetic; len(phonetic) > 0 {
		fmt.Printf("发音：\t%s\n", phonetic)
	}

	if len(j.Basic.Explains) > 0 {
		fmt.Printf("释义：\n\t- %s\n", strings.Join(j.Basic.Explains[:], "\n\t- "))
	}

	if len(j.Web) > 0 {
		fmt.Println("例子：")
		for _, v := range j.Web {
			fmt.Printf("\t%s：\n\t\t- %s\n", v.Key, strings.Join(v.Value[:], "\n\t\t- "))
		}
	}

	if _, ok := keys[2]; ok {
		cmd := exec.Command("say", input)
		if err = cmd.Run(); err != nil {
			fmt.Println("出错啦：不好意思，我不知道怎么读")
		}
	}
}
