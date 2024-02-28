package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Rep struct {
	Code int64 `json:"code"`
}
type Client struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func main() {
	var name string
	var serverUrt string
	var del string

	flag.StringVar(&name, "name", "客户端", "监控客户端名字")
	flag.StringVar(&serverUrt, "url", "0.0.0.0:80", "服务端地址")
	flag.StringVar(&del, "del", "false", "删除客户端")
	flag.Parse()

	id, err := machineid.ID()
	if err != nil {
		log.Println("无法获取机器ID:", err)
		return
	}
	clientBody := Client{
		Name: name,
		Id:   id,
	}
	bodyJson, err := json.Marshal(clientBody)
	if err != nil {
		log.Println("解析失败", err)
		return
	}

	if del == "true" {
		// 删除操作
		resp, err := sendRequest(fmt.Sprintf("%s/api/delete", serverUrt), bodyJson)
		if err != nil {
			log.Println("删除失败", err)
			return
		} else {
			BodyJson(resp)
		}
	} else if del == "false" {
		// 注册
		resp, err := sendRequest(fmt.Sprintf("%s/api/register", serverUrt), bodyJson)

		if err != nil {
			log.Println("注册失败", err)
		} else {
			BodyJson(resp)
		}

		// 探活
		for {
			resp, err := sendRequest(fmt.Sprintf("%s/api/alive", serverUrt), bodyJson)
			if err != nil {
				log.Println("注册失败", err)
			} else {
				BodyJson(resp)
			}
			time.Sleep(30 * time.Second)
		}
	}

}
func BodyJson(r *http.Response) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("读取返回结果失败", err)
		return
	}
	var b Rep
	if err = json.Unmarshal(body, &b); err != nil {
		log.Println("解析失败", err)
		return
	}
	if b.Code == 1000 {
		log.Println("请求成功🏅️")
	}
}

func sendRequest(url string, payload []byte) (*http.Response, error) {
	// 创建一个新的请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
