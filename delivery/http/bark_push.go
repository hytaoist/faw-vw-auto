package http

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

const (
	// Group
	FAWVWGroupName = "一汽大众"
	ErrorGroupName = "服务异常"
	TitleError     = "异常"
	// Title - 一汽大众
	TitleSignin = "签到"
)

var (
	PushServerURL string
)

func SetPushServerURL(url string) {
	PushServerURL = url
}

/*
推送消息
*/
func Push(group string, title string, body string) {
	if PushServerURL == "" {
		fmt.Println(time.Now(), "Bark推送服务器未配置")
		return
	}
	//1.
	// 定义请求体数据
	bodyString := fmt.Sprintf(`{"group": "%s", "title": "%s", "body": "%s"}`, group, title, body)
	bodyData := []byte(bodyString)

	// 创建请求
	req, err := http.NewRequest("POST", PushServerURL, bytes.NewBuffer(bodyData))
	if err != nil {
		fmt.Println(time.Now(), "创建通知异常", err)
		return
	}

	// 设置请求头
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(time.Now(), "发送通知异常", err)
		return
	}
	defer resp.Body.Close()
}

/*
服务异常时推送
*/
func ErrorPush(err error) {
	Push(ErrorGroupName, TitleError, err.Error())
}
