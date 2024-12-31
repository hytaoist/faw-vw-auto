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
	. "github.com/hytaoist/faw-vw-auto/infrastructure/database"
	"github.com/robfig/cron/v3"
)

var (
	RETURN_STATUS_SUCCEED = "SUCCEED"
	RETURN_STATUS_FAILED  = "FAILED"

	// 请求的默认Header
	defaultHeaders = http.Header{
		"entry":               []string{"vwapp"},
		"Content-Type":        []string{"application/json"},
		"x-microservice-name": []string{"api-gateway"},
		"x-namespace-code":    []string{"production"},
	}

	// 登录的请求body
	signinRequestBody SigninRequestBody
	securityCodeBody  SecurityCodeBody
)

// 通用Resp
type OneAppResp struct {
	ReturnStatus string      `json:"returnStatus"`
	ErrorCode    string      `json:"errorCode"`
	ErrorMessage string      `json:"errorMessage"`
	Data         interface{} `json:"data"`
}

type CheckInDataResp struct {
	// 连续签到天数
	ContinueCheckInDays int `json:"continueCheckInDays"`
	// 是否可以开积分盲盒
	Lottery bool `json:"lottery"`
	// 是否第一次签到
	FirstCheckIn         bool               `json:"firstCheckIn"`
	CheckInOnOtherDevice bool               `json:"checkInOnOtherDevice"`
	CheckInDataList      []*CheckInDataItem `json:"checkInDataList"`
}

// 查询签到信息响应报文
type GetCheckInDataResp struct {
	CheckInToday        bool
	ContinueCheckInDays int
	LotteryCheckInDays  int
}

type CheckInDataItem struct {
	Date    string `json:"date"`
	Checkin bool   `json:"checkin"`
	Score   int    `json:"score"`
	Today   bool   `json:"today"`
}

// 登录时的请求body
type SigninRequestBody struct {
	Password     string `json:"password"`
	GraphCode    string `json:"graphCode"`
	SecurityCode string `json:"securityCode"`
	Mobile       string `json:"mobile"`
}

type SecurityCodeBody struct {
	SecurityCode string `json:"securityCode"`
}

type FAW_VW struct {
	psql *Psql
}

func NewFAW(p *Psql) *FAW_VW {
	b := &FAW_VW{p}
	return b
}

func (fawvw *FAW_VW) LoadConfig(cfg *Config) {
	// signinRequestBody
	signinRequestBody.GraphCode = ""
	signinRequestBody.Mobile = cfg.Mobile
	signinRequestBody.Password = cfg.Password
	signinRequestBody.SecurityCode = cfg.SecurityCode
	// securityCode
	securityCodeBody.SecurityCode = cfg.SecurityCode

	defaultHeaders.Set("did", cfg.Did)
}

func (fawvw *FAW_VW) BackgroundRunning() {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(loc))
	c.AddFunc("10 10 * * *", func() {
		fmt.Println(time.Now(), "一汽大众每日签到")
		fawvw.Running()
	})
	c.Start()
}

// 任务执行过程
func (fawvw *FAW_VW) Running() {
	// 先查出最新的token
	auth, _ := fawvw.psql.FindLatestOne()
	var authorization = ""
	if auth != nil {
		authorization = auth.TokenType + auth.AccessToken
	}

	authorization = fawvw.getValidToken(authorization)

	// 1.查签到信息
	oneappResp, err := fawvw.getCheckInInfo(authorization)

	if err != nil {
		fmt.Println(time.Now(), "查询签到信息异常", err)
		Push(FAWVWGroupName, TitleError, err.Error())
	} else {
		checkinInfo, ok := oneappResp.Data.(*GetCheckInDataResp)
		if !ok {
			fmt.Println(time.Now(), "类型断言失败")
			return
		}
		if checkinInfo.CheckInToday {
			msg := "今天已经签过到了，无需进行签到!"
			fmt.Println(time.Now(), msg)
			Push(FAWVWGroupName, TitleSignin, msg)
			return
		}
	}

	// 2.执行签到
	checkinV1Resp, err := fawvw.checkinV1(authorization)
	if err != nil {
		Push(FAWVWGroupName, TitleSignin, err.Error())
		return
	}
	checkinV1, ok := checkinV1Resp.Data.(*CheckInDataResp)
	if !ok {
		errMSG := "签到返回类型断言失败"
		fmt.Println(time.Now(), errMSG)
		Push(FAWVWGroupName, TitleError, errMSG)
		return
	}

	seccessMSG := fmt.Sprintf("签到成功！连续签到天数：%d, 是否开盲盒：%t", checkinV1.ContinueCheckInDays, checkinV1.Lottery)
	Push(FAWVWGroupName, TitleSignin, seccessMSG)

	// 3.根据签到信息，能开盲盒时自动打开盲盒。
	if checkinV1.Lottery {
		lotteryV1Resp, err := fawvw.lotteryV1(authorization)
		if err != nil {
			errMSG := err.Error()
			fmt.Println(time.Now(), errMSG)
			Push(FAWVWGroupName, TitleError, errMSG)
			return
		}

		if lotteryV1Resp.ReturnStatus == RETURN_STATUS_SUCCEED {
			Push(FAWVWGroupName, TitleSignin, "打开盲盒成功！")
		} else {
			errMSG := fmt.Sprintf("打开盲盒失败！%s", lotteryV1Resp.ErrorMessage)
			fmt.Println(time.Now(), errMSG)
			Push(FAWVWGroupName, TitleError, errMSG)
		}
	}
}

// 获取有效的Token
// 1.查询签到信息接口能正常返回，Token有效
// 2.查询签到信息接口异常返回（未授权等），调用登录接口获取最新Token并入库
func (fawvw *FAW_VW) getValidToken(authorization string) string {
	if authorization != "" {
		// 1.检验当前这个token是否有效（通过查询签到信息接口）
		oneappResp, err := fawvw.getCheckInInfo(authorization)
		if err != nil {
			fmt.Println(time.Now(), "查询签到信息异常", err)
		}
		// ReturnStatus == "SUCCEED"，表示当前的这个请求能正常执行，且有返回结果，authorization有效！
		if oneappResp != nil && oneappResp.ReturnStatus == RETURN_STATUS_SUCCEED {
			return authorization
		}
	}

	// 1.1未授权/Token 过期，需重新获取Token
	// 1.签到获取AccessToken，并入库
	newAuthorization, err := fawvw.signinByPassword()
	if err != nil {
		fmt.Println(time.Now(), "执行登陆获取授权异常", err)
		Push(FAWVWGroupName, TitleError, err.Error())
		return ""
	}

	return newAuthorization
}

func (favw *FAW_VW) signinByPassword() (string, error) {
	//1.
	// 定义目标 URL
	targetURL := "https://oneapp-api.faw-vw.com/account/login/loginByPassword/v1"

	// 创建参数
	params := url.Values{}

	// 定义请求体数据
	bodyData, err := json.Marshal(signinRequestBody)
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
	// 添加默认 headers 到请求中
	for key, values := range defaultHeaders {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

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
	var auth FAWAuth
	oneappResp := &OneAppResp{Data: &auth}
	errJson := json.Unmarshal(body, &oneappResp)
	if errJson != nil {
		return "", err
	}

	// 3.将授权入库
	err = favw.psql.CreateFAW_Auth(resp.Request.Context(), &auth)
	if err != nil {
		return "", err
	}
	var authorization string = ""
	authorization = auth.TokenType + auth.AccessToken
	return authorization, nil
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
		fmt.Println(time.Now(), "签到成功，解析返回结果异常", err)
		return nil, err
	}
	fmt.Println(time.Now(), "一汽大众签到成功！")
	return checkinV1, nil
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
