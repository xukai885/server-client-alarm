package tosend

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"server-client-alarm/settings"
)

// 发送webhook
func ToWebHook(msg string) {
	webhook_url := settings.Conf.Webhook.Url
	token := settings.Conf.Webhook.Token
	// 设置请求体
	msgs := struct {
		Message string `json:"message"`
		Title   string `json:"title"`
		To      string `json:"to"`
	}{Message: msg,
		Title: settings.Conf.Webhook.Title,
		To:    settings.Conf.Webhook.To,
	}

	// 将结构体转换为 JSON 字符串
	payloadJSON, err := json.Marshal(msgs)
	if err != nil {
		log.Println("序列化失败", err)
		return
	}
	// 发送 POST 请求
	resp, err := sendRequest(webhook_url, token, payloadJSON)
	if err != nil {
		log.Println("请求webhook出错", err)
		return
	}
	defer resp.Body.Close()
}

func sendRequest(url, token string, payload []byte) (*http.Response, error) {
	// 创建一个新的请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", token)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
