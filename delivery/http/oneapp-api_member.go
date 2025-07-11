package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

/*
*
查会员当前的积分数
*
*/
func (fawvw *FAW_VW) getManufacturerScore(authorization string) (*ManufacturerScoreResp, error) {

	targetURL := "https://oneapp-api.faw-vw.com/profile/member/getManufacturerScore/v1"
	params := url.Values{}
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
	req.Header.Add("authorization", authorization)
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
	// 2.使用decoder来解析
	reader := strings.NewReader(string(body))
	fmt.Println(time.Now(), "请求响应Body：", string(body))

	var facturerScores []ManufacturerScoreResp
	oneappResp := &OneAppResp{Data: &facturerScores}
	decoder := json.NewDecoder(reader)
	errJson := decoder.Decode(oneappResp)
	if errJson != nil {
		fmt.Println(time.Now(), "解析失败", errJson)
		return nil, err
	}

	if len(facturerScores) > 0 {
		return &facturerScores[0], nil
	} else {
		return nil, nil
	}
}

/*
*
当日新增数量
*
*/
func (fawvw *FAW_VW) getDayIncrease(authorization string) (*int, error) {

	targetURL := "https://oneapp-api.faw-vw.com/profile/member/getDayIncrease/v1"
	params := url.Values{}
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
	req.Header.Add("authorization", authorization)
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
	var increase int
	// 2.使用decoder来解析
	reader := strings.NewReader(string(body))
	fmt.Println(time.Now(), "请求响应Body：", string(body))
	oneappResp := &OneAppResp{Data: &increase}
	decoder := json.NewDecoder(reader)
	errJson := decoder.Decode(oneappResp)
	if errJson != nil {
		fmt.Println(time.Now(), "解析失败", errJson)
		return nil, err
	}

	return &increase, nil
}
