package template

import (
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	model "schedule/model/bot"
	"schedule/util"
)

func TestTranslateKeyboardForUser(t *testing.T) {
	mockTranslator := new(MockTranslator)

	mockTranslator.On("Tr", "en", "SettingsCountry").Return("Countries")
	mockTranslator.On("Tr", "en", "SettingsCompetition").Return("Competitions")
	mockTranslator.On("Tr", "en", "SettingsAlert").Return("Alerts")

	util.SetTranslator(mockTranslator)

	assert.Equal(
		t,
		&models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Countries", CallbackData: "settings_country"},
					{Text: "Competitions", CallbackData: "settings_competition"},
					{Text: "Alerts", CallbackData: "settings_alert"},
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
