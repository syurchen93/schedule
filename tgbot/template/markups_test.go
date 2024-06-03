package template

import (
	"testing"
	"time"

	model "schedule/model/bot"
	"schedule/tgbot/manager"
	"schedule/util"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/syurchen93/api-football-client/common"
)

func TestTranslateKeyboardForUser(t *testing.T) {
	mockTranslator := new(MockTranslator)

	mockTranslator.On("Tr", "en", "SettingsCountry").Return("Countries")
	mockTranslator.On("Tr", "en", "SettingsCompetition").Return("Competitions")
	mockTranslator.On("Tr", "en", "SettingsUser").Return("User")
	mockTranslator.On("Tr", "en", "SettingsAlert").Return("Alerts")
	mockTranslator.On("Tr", "en", "ToSchedule").Return("ToSchedule")

	util.SetTranslator(mockTranslator)

	assert.Equal(
		t,
		&models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Countries", CallbackData: "settings_country"},
					{Text: "Competitions", CallbackData: "settings_competition"},
				},
				{
					{Text: "Alerts", CallbackData: "settings_alert"},
					{Text: "User", CallbackData: "settings_user"},
				},
				{
					ButtonSchedule,
				},
			},
		},
		TranslateKeyboardForUser(
			model.User{
				Locale: "en",
			},
			KeyboardSettingsGeneral,
		),
	)
}

func Test_generateFixtureButtonText(t *testing.T) {
	user := model.User{
		EnableSpoilers: false,
	}
	fixture := manager.FixtureView{
		HomeTeamName: "Borussia Dortmund with extra words",
		AwayTeamName: "Real Madrid with extra words",
		Date:         time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC),
		IsToggled:    false,
		Status:       common.NotStarted,
	}

	assert.Equal(
		t,
		"🔔 Fri 1.01 12:00 Borussia Dortmund with extra words  Real Madrid with extra words",
		generateFixtureButtonText(user, fixture),
	)

	fixture.AwayTeamCode = "RMA"
	fixture.HomeTeamCode = "BVB"

	assert.Equal(
		t,
		"🔔 Fri 1.01 12:00 BVB  RMA",
		generateFixtureButtonText(user, fixture),
	)
}

type MockTranslator struct {
	mock.Mock
}

func (m *MockTranslator) Tr(locale string, key string, args ...interface{}) string {
	mockArgs := m.Called(locale, key)
	return mockArgs.String(0)
}
