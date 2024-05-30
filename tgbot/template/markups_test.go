package template

import (
	"testing"

	model "schedule/model/bot"
	"schedule/util"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type MockTranslator struct {
	mock.Mock
}

func (m *MockTranslator) Tr(locale string, key string, args ...interface{}) string {
	mockArgs := m.Called(locale, key)
	return mockArgs.String(0)
}
