package messages

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/maps"
	"strings"
	"time"
)

type RoutingData struct {
	Routing string `json:"routing"`
}

type TimeData struct {
	Routing string       `json:"routing"`
	Time    HoursMinutes `json:"time"`
}

type HoursMinutes time.Time

func (h *HoursMinutes) UnmarshalJSON(bytes []byte) (err error) {
	value := strings.Trim(string(bytes), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	result, err := time.Parse("15:04", value)
	*h = HoursMinutes(result)
	return err
}

func (h HoursMinutes) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(h).Format("15:04") + `"`), nil
}

func SelectTime(chatID int64, routing string, data map[string]interface{}) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, "Выберите время")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(func() (buttons [][]tgbotapi.InlineKeyboardButton) {
		var (
			now              = time.Now()
			minHour, maxHour = 10, 22
		)

		t := time.Date(now.Year(), now.Month(), now.Day(), minHour, 0, 0, 0, now.Location())

		if data == nil {
			data = make(map[string]interface{})
		}

		buttonData := func(time time.Time) string {
			dataToMerge := data
			maps.Copy(dataToMerge, map[string]interface{}{
				"routing": routing,
				"time":    time.Format("15:04"),
			})

			marshalledData, _ := json.Marshal(dataToMerge)
			return string(marshalledData)
		}

		for t.Hour() != maxHour {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(t.Format("15:04"), buttonData(t)),
				tgbotapi.NewInlineKeyboardButtonData(t.Add(time.Minute*30).Format("15:04"), buttonData(t.Add(time.Minute*30))),
				tgbotapi.NewInlineKeyboardButtonData(t.Add(time.Hour).Format("15:04"), buttonData(t.Add(time.Hour))),
				tgbotapi.NewInlineKeyboardButtonData(t.Add(time.Hour+time.Minute*30).Format("15:04"), buttonData(t.Add(time.Hour+time.Minute*30))),
			))

			t = t.Add(time.Hour * 2)
		}

		return buttons
	}()...)

	return &msg
}
