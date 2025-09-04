package manager

import (
	"context"
	"fmt"
	"time"

	model "schedule/model/bot"
	"schedule/util"

	"github.com/go-telegram/bot"
)

func GetAndFireAlerts(ctx context.Context, b *bot.Bot) bool {
	alertsToFire := getAlertsToFire()
	for _, alert := range alertsToFire {
		fireAlert(ctx, b, alert)
	}

	return true
}

func GetAlertCompetitionViewsForUser(userId int) []CompetitionView {
	var alerts []model.Alert

	dbGorm.
		Joins("join fixture on alert.fixture_id = fixture.id").
		Preload("Fixture").
		Preload("User").
		Preload("Fixture.HomeTeam").
		Preload("Fixture.AwayTeam").
		Preload("Fixture.Competition").
		Preload("Fixture.Competition.Country").
		Where("user_id = ? and is_fired = 0 AND fixture.date < ?", userId, time.Now().AddDate(0, 0, DefaultDaysInFuture)).
		Order("fixture.date ASC").
		Find(&alerts)

	return CreateCompetitionFixtureViewFromAlers(alerts)
}

func fireAlert(ctx context.Context, b *bot.Bot, alert model.Alert) {
	success, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: alert.User.ID,
		Text: fmt.Sprintf(
			"⏰ %s %s %s %s %d %s",
			alert.Fixture.HomeTeam.Name,
			util.Translate(alert.User.Locale, "Vs"),
			alert.Fixture.AwayTeam.Name,
			util.Translate(alert.User.Locale, "Starts"),
			alert.TimeBefore/60,
			util.Translate(alert.User.Locale, "Minutes"),
		),
	})

	if success != nil && err == nil {
		dbGorm.Model(&alert).Update("is_fired", 1)
	}
}

func getAlertsToFire() []model.Alert {
	var alerts []model.Alert

	dbGorm.
		Preload("User").
		Joins("join fixture on alert.fixture_id = fixture.id").
		Where("is_fired = ? AND DATE_ADD(fixture.date, INTERVAL - alert.time_before SECOND) <= NOW()", 0).
		Preload("User").
		Preload("Fixture").
		Preload("Fixture.HomeTeam").
		Preload("Fixture.AwayTeam").
		Find(&alerts)

	return alerts
}

func createOrDeleteAlertForFixture(user *model.User, fixtureId int) {
	var existingAlert model.Alert
	result := dbGorm.Where("user_id = ? AND fixture_id = ? AND time_before = ?", 
		user.ID, fixtureId, user.AlertOffset).First(&existingAlert)
	
	if result.Error != nil {
		// Alert doesn't exist, create it
		alert := model.Alert{
			UserID:     uint(user.ID),
			FixtureID:  uint(fixtureId),
			TimeBefore: user.AlertOffset,
		}
		dbGorm.Create(&alert)
	} else {
		// Alert exists, delete it
		dbGorm.Delete(&existingAlert)
	}
}
