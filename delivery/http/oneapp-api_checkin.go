package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// 签到
func (fawvw *FAW_VW) checkinV1(authorization string) (*OneAppResp, error) {
	//1.
	// 定义目标 URL
	targetURL := "https://oneapp-api.faw-vw.com/profile/checkin/v1"

	// 定义请求体数据
	bodyData, err := json.Marshal(securityCodeBody)
	if err != nil {
		fmt.Println(time.Now(), "登录请求签到解析异常请排查", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(bodyData))
	if err != nil {
		fmt.Println(time.Now(), "创建请求异常", err)
		return nil, err
	}

	// 设置请求头
	// 添加默认 headers 到请求中
	for key, values := range defaultHeaders {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}
	// 设置授权
	req.Header.Set("authorization", authorization)
	// 添加参数到 URL
	// req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(time.Now(), "签到请求异常", err)
		return nil, err
	}
	defer resp.Body.Close()
	// 2.解析结果
	body, err := io.ReadAll(resp.Body)

	var checkinData CheckInDataResp
	checkinV1 := &OneAppResp{Data: &checkinData}
	errJson := json.Unmarshal(body, &checkinV1)
	if errJson != nil {
		fmt.Println(time.Now(), "成功，解析返回结果异常", err)
		return nil, err
	}
	fmt.Println(time.Now(), "一汽大众签到成功！")
	return checkinV1, nil
}

// 查询签到信息
func (favw *FAW_VW) getCheckInInfo(authorization string) (*OneAppResp, error) {
	//1.
	// 定义目标 URL
	targetURL := "https://oneapp-api.faw-vw.com/profile/checkin/data/v1"

	// 创建参数
	params := url.Values{}

	// 创建请求
	req, err := http.NewRequest("GET", targetURL, nil)

	if err != nil {
		fmt.Println(time.Now(), "创建请求异常", err)
		return nil, err
	}

	// 设置请求头
	// 添加默认 headers 到请求中
	for key, values := range defaultHeaders {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	// 设置授权
	req.Header.Add("authorization", authorization)
	// 添加参数到 URL
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(time.Now(), "请求异常", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 请求未授权
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println(time.Now(), "请求未授权")
		return nil, errors.New(string(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	var checkinInfo GetCheckInDataResp
	// fmt.Println(time.Now(), "请求响应Body：", body)
	// 1.json解析方式
	// oneappResp := &OneAppResp{Data: &checkinInfo}
	// errJson := json.Unmarshal(body, oneappResp)

	// 2.使用decoder来解析
	reader := strings.NewReader(string(body))
	fmt.Println(time.Now(), "请求响应Body：", string(body))
	oneappResp := &OneAppResp{Data: &checkinInfo}
	decoder := json.NewDecoder(reader)
	errJson := decoder.Decode(oneappResp)
	if errJson != nil {
		fmt.Println(time.Now(), "解析查询签到信息失败", errJson)
		return nil, err
	}

	return oneappResp, nil
}

/*
*
打开盲盒
*/
func (fawvw *FAW_VW) lotteryV1(authorization string) (*OneAppResp, error) {
	targetURL := "https://oneapp-api.faw-vw.com/profile/lottery/v1"

	// 创建参数
	params := url.Values{}
	// 定义请求体数据
	// 无Body
	// 创建请求
	req, err := http.NewRequest(http.MethodPost, targetURL, nil)
	if err != nil {
		fmt.Println(time.Now(), "创建请求异常", err)
		return nil, err
	}

	// 设置请求头
	// 添加默认 headers 到请求中
	for key, values := range defaultHeaders {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}
	// 设置授权
	req.Header.Set("authorization", authorization)
	// 添加参数到 URL
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(time.Now(), "签到请求异常", err)
		return nil, err
	}
	defer resp.Body.Close()
	// 2.解析结果
	body, err := io.ReadAll(resp.Body)
	lotteryV1 := &OneAppResp{}
	errJson := json.Unmarshal(body, &lotteryV1)
	if errJson != nil {
		fmt.Println(time.Now(), "打开盲盒解析返回结果异常", err)
		return nil, err
	}
	return lotteryV1, nil
}
