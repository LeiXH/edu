package services

import (
	"database/sql"
	"errors"
	"signModule/pkg/models"
	"signModule/pkg/services/mapper"
)

func QueryUserInfo(queryType mapper.UserQueryType, keyword string) (*models.UserInfo, error) {
	info, e := mapper.GetUserInfo(queryType, keyword)
	if e != nil && e == sql.ErrNoRows {
		return nil, errors.New("没有该用户")
	}
	return info, e
}
func DoSIgnForUser(queryType mapper.UserQueryType, keyword string) error {
	err := mapper.SignForUser(queryType, keyword)
	if err != nil {
		return err
	}
	return nil
}
func IsUserSigned(user *models.UserInfo) bool {
	if rc, err := mapper.CountUserSignedRecord(mapper.QueryTypeTelephone, user.Telephone); err != nil || rc == 0 {
		return false
	}
	return true
}

func GetAll(nums int) ([] models.UserInfo, bool) {
	if all, err := mapper.GetAllUser(nums); err == nil {
		return all, true
	}

	return nil, false
}

func InsertData(user *models.UserInfo) error {
	err := mapper.InsertData(user)
	if err != nil {
		return err
	}
	return nil
}

func IsExist(tele string) bool {
	if res, err :=mapper.GetUserByTelephone(tele); err !=nil || res == 0 {
		return true
	}
	return false
}

func GetSignedNum() (int, error) {
	num, err := mapper.GetSignedNum()
	if  err != nil {
		return 0, err
	}
	return num, nil
}
