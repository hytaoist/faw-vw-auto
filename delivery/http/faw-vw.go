package http

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/hytaoist/faw-vw-auto/config"
	// . "github.com/hytaoist/faw-vw-auto/domain"
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

	securityCodeBody SecurityCodeBody
	execFreq         string
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

// 会员积分响应报文
type ManufacturerScoreResp struct {
	Businesstypename string
	Invalidscore     string
	Businessid       int
	// 剩余积分
	Remainscore string
	Frozenscore string
	// 可用积分
	Availablescore    string
	Dealerservicecode string
	DealerServiceName string
	Businessname      string
	Usedscore         string
	Scoretypename     string
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

func (fawvw *FAW_VW) LoadAppConfig(cfg *Config) {
	// securityCode
	securityCodeBody.SecurityCode = cfg.SecurityCode
	defaultHeaders.Set("did", cfg.Did)
	execFreq = cfg.ExecFreq
}

func (fawvw *FAW_VW) BackgroundRunning() {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(loc))
	c.AddFunc(execFreq, func() {
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

	// 3.根据签到信息，能开盲盒时自动打开盲盒。
	if checkinV1.Lottery {
		lotteryV1Resp, err := fawvw.lotteryV1(authorization)
		if err != nil {
			errMSG := err.Error()
			fmt.Println(time.Now(), errMSG)
			Push(FAWVWGroupName, TitleError, errMSG)
			return
		}

		if lotteryV1Resp.ReturnStatus != RETURN_STATUS_SUCCEED {
			errMSG := fmt.Sprintf("打开盲盒失败！%s", lotteryV1Resp.ErrorMessage)
			fmt.Println(time.Now(), errMSG)
			Push(FAWVWGroupName, TitleError, errMSG)
			return
		}
	}

	// 4.查询当日新增积分数以及总积分数，构建整体响应并推送Bark通知。
	dailyIncrease, err := fawvw.getDayIncrease(authorization)
	if err != nil {
		errMSG := err.Error()
		fmt.Println(time.Now(), errMSG)
		Push(FAWVWGroupName, TitleError, errMSG)
	}
	// 4.1 当日新增积分入库
	fawvw.psql.InsertPointRecord(*dailyIncrease)

	manufacturerScore, err := fawvw.getManufacturerScore(authorization)
	if err != nil {
		errMSG := err.Error()
		fmt.Println(time.Now(), errMSG)
		Push(FAWVWGroupName, TitleError, errMSG)
	}

	var successMsg string
	if checkinV1.Lottery {
		successMsg = fmt.Sprintf("成功！连续签到：%d天；积分新增：%d，总可用：%s；已开盲盒", checkinV1.ContinueCheckInDays, *dailyIncrease, manufacturerScore.Availablescore)
	} else {
		successMsg = fmt.Sprintf("成功！连续签到：%d天；积分新增：%d，总可用：%s", checkinV1.ContinueCheckInDays, *dailyIncrease, manufacturerScore.Availablescore)
	}

	Push(FAWVWGroupName, TitleSignin, successMsg)
}

// 获取有效的Token
// 1.查询签到信息接口能正常返回，Token有效
// 2.查询签到信息接口异常返回（未授权等），调用登录接口获取最新Token并入库
func (fawvw *FAW_VW) getValidToken(authorization string) string {
	if authorization != "" {
		// 1.检验当前这个token是否有效
		r, err := fawvw.checkToken(authorization)
		if err != nil {
			fmt.Println(time.Now(), "校验Token接口异常", err)
		}
		if r {
			return authorization
		}
	}

	// 1.1未授权/Token 过期，需重新获取Token
	// 1.签到获取AccessToken，并入库
	newAuthorization, err := fawvw.registeOrLogin()
	if err != nil {
		fmt.Println(time.Now(), "执行登陆获取授权异常", err)
		Push(FAWVWGroupName, TitleError, err.Error())
		return ""
	}

	return newAuthorization
}
