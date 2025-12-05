package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	. "github.com/hytaoist/faw-vw-auto/config"
	. "github.com/hytaoist/faw-vw-auto/domain"
)

// Web应用通用的响应体
type WebAppResp struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// 登录响应报文
type RegisterOrLoginResp struct {
	Token    string `json:"token"`
	OtdToken string `json:"otdToken"`
	Account  string `json:"account"`
	Msg      string `json:"msg"`
	Extmsg   string `json:"extmsg"`
	// ErrorCode                string `json:"errorCode"`
	LoginStatus              string `json:"loginStatus"`
	ErrorVerificationCodeNum int    `json:"errorVerificationCodeNum"`
	ErrorPasswordNum         int    `json:"errorPasswordNum"`
}

type RegisterOrLoginWebReq struct {
	Account                 string `json:"account"`
	Password                string `json:"password"`
	Did                     string `json:"did"`
	GraphVerificationCode   string `json:"graphVerificationCode"`
	GraphVerificationCodeId string `json:"graphVerificationCodeId"`
	VerificationCode        string `json:"verificationCode"`
}

var (
	registerOrLoginWebReq RegisterOrLoginWebReq
)

func (fawvw *FAW_VW) LoadWebAPIConfig(cfg *Config) {
	registerOrLoginWebReq.Account = cfg.Mobile
	// 注意这个Web登录请求使用的password非密码明文，是一串加密后的字符，抓接口获取
	registerOrLoginWebReq.Password = cfg.Password
	registerOrLoginWebReq.Did = cfg.WebDid
	registerOrLoginWebReq.GraphVerificationCode = ""
	registerOrLoginWebReq.GraphVerificationCodeId = ""
	registerOrLoginWebReq.VerificationCode = ""
}

func (favw *FAW_VW) checkToken(token string) (bool, error) {
	targetURL := "https://vw.faw-vw.com/api/business/cpoint/checkToken"

	req, err := http.NewRequest("POST", targetURL, nil)
	if err != nil {
		fmt.Println(time.Now(), "创建请求异常", err)
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(time.Now(), "请求异常", err)
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(time.Now(), "请求失败", resp.Status)
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	var r bool
	webAppResp := &WebAppResp{Data: &r}
	errJson := json.Unmarshal(body, &webAppResp)
	if errJson != nil {
		return false, err
	}
	return r, nil
}

func (favw *FAW_VW) registeOrLogin() (string, error) {
	//1.
	// 定义目标 URL
	targetURL := "https://vw.faw-vw.com/api/business/cpoint/registeOrLogin"

	// 创建参数
	params := url.Values{}

	// 定义请求体数据
	bodyData, err := json.Marshal(registerOrLoginWebReq)
	if err != nil {
		fmt.Println(time.Now(), "登录请求body解析异常请排查", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(bodyData))
	if err != nil {
		fmt.Println(time.Now(), "创建请求异常", err)
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 添加参数到 URL
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(time.Now(), "请求异常", err)
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println(time.Now(), "请求失败", resp.Status)
		return "", err
	}
	defer resp.Body.Close()
	// 2.解析结果
	body, err := io.ReadAll(resp.Body)
	var registerOrLoginResp RegisterOrLoginResp
	webAppResp := &WebAppResp{Data: &registerOrLoginResp}
	errJson := json.Unmarshal(body, &webAppResp)
	// 2.1 具有图形验证码，不能校验通过时，需要给用户发送需要更新“password”参数通知
	if webAppResp.Code == 500 && webAppResp.Msg == "验证失败" {
		fmt.Println(time.Now(), "登录验证失败", err)
		return "", errors.New(webAppResp.Msg)
	}

	if errJson != nil {
		return "", err
	}

	// 3.拼接参数，将授权入库
	tokenType := "Bearer "
	accessToken := strings.Replace(registerOrLoginResp.Token, tokenType, "", 1)

	err = favw.psql.CreateFAW_Auth(&FAWAuth{AccessToken: accessToken, TokenType: tokenType, ExpiresIn: "100"})
	if err != nil {
		return "", err
	}
	return (tokenType + accessToken), nil
}
