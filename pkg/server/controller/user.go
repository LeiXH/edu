package controller

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"net/http"
	"signModule/pkg/logger"
	"signModule/pkg/models"
	"strconv"

	"signModule/pkg/server/controller/form"
	"signModule/pkg/server/controller/response"
	"signModule/pkg/services"
	"signModule/pkg/services/mapper"
	"signModule/pkg/services/printer"
)

func GinDoUserSignManually( c echo.Context) error {

	signForm := new(form.SignForm)
	if err := c.Bind(signForm); err != nil {
		return c.JSON(http.StatusBadRequest, response.RenderError(response.INVALID_PARAM, err))
	}

	queryType := mapper.UserQueryType(signForm.SignType)

	userInfo, err := services.QueryUserInfo(queryType, signForm.Keyword)
	if err != nil {
		return c.JSON(http.StatusOK, response.RenderError(response.INVALID_PARAM, err))
	}

	if signed := services.IsUserSigned(userInfo); !signed {
		_ = services.DoSIgnForUser(mapper.QueryTypeTelephone, userInfo.Telephone)
		//printer.PrintUserLabel(userInfo)
	}
	if queryType == mapper.QueryTypeQrCode {
		if webSocketConned.Load() == true {
			userInfoChan <- *userInfo
		}
	}

	return c.JSON(http.StatusOK, response.RenderSuccess(userInfo))
}

func GinJustPrintUserLabel(c echo.Context) error{
	workNum := c.FormValue("num")
	if len(workNum) == 0 {
		return c.JSON(http.StatusOK, response.RenderError(response.INVALID_PARAM, errors.New("no args")))
	}
	userInfo, err := services.QueryUserInfo(mapper.QueryTypeTelephone, workNum)
	if err != nil {
		return c.JSON(http.StatusOK, response.RenderError(response.INVALID_PARAM, err))
	}
	printer.PrintUserLabel(userInfo)
	return c.JSON(http.StatusOK, response.RenderSuccess(""))
}

func GinFaceSign(c echo.Context) error {
	key := c.FormValue("key")
	if len(key) == 0 {
		return c.JSON(http.StatusOK, response.RenderError(response.INVALID_PARAM, errors.New("no args")))
	}
	userInfo, err := services.QueryUserInfo(mapper.QueryTypeQrCode, key)

	if err != nil {
		return c.JSON(http.StatusOK, response.RenderError(response.INVALID_PARAM, err))
	}

	userInfoChan <- *userInfo

	printer.PrintUserLabel(userInfo)
	return c.JSON(http.StatusOK, response.RenderSuccess(userInfo))
}

func Import(c echo.Context) error  {
	data := c.FormValue("data")
	var dataMap []map[string]string

	var user models.UserInfo

	json.Unmarshal([]byte(data),&dataMap)

	for _, one := range dataMap {
		user.Name = one["name"]
		user.Signed, _ = strconv.Atoi(one["signed"])
		user.Conference,_ = strconv.Atoi(one["conference"])
		user.Company = one["company"]
		user.Telephone = one["telephone"]

		user.Meeting = ""
		if one["meeting6"] ==  "是"  {
			user.Meeting += "6"
		}
		if one["meeting5"] ==  "是"  {
			user.Meeting += "5"
		}
		if one["meeting4"] ==  "是"  {
			user.Meeting += "4"
		}
		if one["meeting3"] ==  "是"  {
			user.Meeting += "3"
		}
		if one["meeting2"] ==  "是"  {
			user.Meeting += "2"
		}
		if one["meeting2"] ==  "是"  {
			user.Meeting += "1"
		}
		user.Mark = one["mark"]
		re :=services.IsExist(user.Telephone)
		if re == true {
			err := services.InsertData(&user)
			if err != nil {
				logger.Infof("insert error %s", err)
			}
		}
	}
	return c.JSON(http.StatusOK, response.RenderSuccess(""))
}

func All(c echo.Context) error  {

	nums, _ := strconv.Atoi(c.FormValue("num"))
	users, ok :=services.GetAll(nums)

	if ok {
		return c.JSON(http.StatusOK, response.RenderSuccess(users))
	}
	return c.JSON(http.StatusOK, response.RenderSuccess(""))
}

func GetSignedNum(c echo.Context) error {
	nums, err := services.GetSignedNum()
	if err != nil{
		return c.JSON(http.StatusOK, response.RenderError(response.INVALID_PARAM, err))
	}
	return c.JSON(http.StatusOK, response.RenderSuccess(nums))
}

