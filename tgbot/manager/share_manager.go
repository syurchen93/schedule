package manager

import (
	"schedule/db"
	"schedule/model/bot"
)

func GetShareSubscriptionsForUser(userId int) ([]bot.UserShare, error) {
	var userShares []bot.UserShare
	err := db.InitDbOrPanic().Where("target_user_id = ?", userId).
		Preload("SourceUser").
		Find(&userShares).Error
	if err != nil {
		return nil, err
	}
	return userShares, nil
}

func SubUserByTargetUsername(user *bot.User, targetUsername string) error {
	var targetUser bot.User
	err := db.InitDbOrPanic().Where("username = ?", targetUsername).First(&targetUser).Error
	if err != nil {
		return err
	}

	userShare := bot.UserShare{
		SourceUserId: uint(user.ID),
		TargetUserId: targetUser.ID,
	}
	err = db.InitDbOrPanic().Create(&userShare).Error
	if err != nil {
		return err
	}
	return nil
}
