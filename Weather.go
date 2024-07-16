package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bregydoc/gtranslate"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я бот, который предоставляет информацию о погоде. Используй команды /today <город> для погоды на сегодня и /week <город> для прогноза на неделю.")
				bot.Send(msg)
			case "help":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Используй команды /today <город> для погоды на сегодня и /week <город> для прогноза на неделю.")
				bot.Send(msg)
			default:
				if strings.HasPrefix(update.Message.Text, "/today") {
					location := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/today"))
					weatherText, err := getTodayWeather(location)
					if err != nil {
						log.Printf("Ошибка при получении погоды: %v", err)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении погоды.")
						bot.Send(msg)
						continue
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, weatherText)
					bot.Send(msg)
				} else if strings.HasPrefix(update.Message.Text, "/week") {
					location := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/week"))
					forecastText, err := getWeekWeather(location)
					if err != nil {
						log.Printf("Ошибка при получении прогноза: %v", err)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении прогноза.")
						bot.Send(msg)
						continue
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, forecastText)
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда. Попробуй /start или /help.")
					bot.Send(msg)
				}
			}
		}
	}
}

func translateToEnglish(text string) (string, error) {
	translated, err := gtranslate.TranslateWithParams(
		text,
		gtranslate.TranslationParams{
			From: "ru",
			To:   "en",
		},
	)
	if err != nil {
		return "", err
	}
	return translated, nil
}

func getTodayWeather(location string) (string, error) {
	if location == "" {
		return "", fmt.Errorf("название города не указано в запросе")
	}

	translatedLocation, err := translateToEnglish(location)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://www.meteoservice.ru/weather/today/%s", strings.ReplaceAll(strings.ToLower(translatedLocation), " ", "-"))
	log.Printf("Отправка запроса на URL: %s", url)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("ошибка запроса: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	// Извлекаем заголовок
	header := strings.TrimSpace(doc.Find("h5").First().Text())

	var weatherDetails []string
	doc.Find(".row.small-collapse.medium-uncollapse.align-middle").Each(func(i int, s *goquery.Selection) {
		time := s.Find(".smedium-1.column.time.text-center.medium-text-left .value").Text()
		weatherCondition := s.Find(".column.text-center.medium-text-left.weather .column.show-for-smedium.text-left").Text()
		temperature := s.Find(".small-2.smedium-1.columns.temperature.text-center .value").Text()

		time = fmt.Sprintf("%s:00", time)

		timeEmoji := getTimeEmoji(time)
		weatherEmoji := getWeatherEmoji(weatherCondition)

		weatherDetails = append(weatherDetails, fmt.Sprintf("%s Время: %s", timeEmoji, time))
		weatherDetails = append(weatherDetails, fmt.Sprintf("%s Погода: %s", weatherEmoji, weatherCondition))
		weatherDetails = append(weatherDetails, fmt.Sprintf("🌡️ Температура: %s", temperature))
		weatherDetails = append(weatherDetails, "---------------------")
	})

	headerText := fmt.Sprintf("%s сегодня по часам\n\n", header)
	response := headerText + strings.Join(weatherDetails, "\n")
	return response, nil
}

func getWeekWeather(location string) (string, error) {
	if location == "" {
		return "", fmt.Errorf("название города не указано в запросе")
	}

	translatedLocation, err := translateToEnglish(location)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://www.meteoservice.ru/weather/week/%s", strings.ReplaceAll(strings.ToLower(translatedLocation), " ", "-"))
	log.Printf("Отправка запроса на URL: %s", url)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("ошибка запроса: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	// Извлекаем заголовок
	header := strings.TrimSpace(doc.Find("h1.text-center.medium-text-left").First().Text())

	var forecast strings.Builder

	doc.Find(".forecast-week-overview .column.text-center").Each(func(i int, s *goquery.Selection) {
		day := s.Find(".weekday").Text()
		maxTemp := s.Find("span.value[title='Макс.']").Text()
		minTemp := s.Find("span.value[title='Мин.']").Text()

		if day != "" && (maxTemp != "" || minTemp != "") {
			forecast.WriteString(fmt.Sprintf("%s: Макс: %s, Мин: %s\n", day, maxTemp, minTemp))
		}
	})

	if forecast.Len() == 0 {
		return "", fmt.Errorf("не удалось получить прогноз погоды. Попробуйте еще раз с другим названием.")
	}

	headerText := fmt.Sprintf("%s\n\n", header)
	return headerText + forecast.String(), nil
}

// getTimeEmoji returns an emoji based on the hour of the day
func getTimeEmoji(time string) string {
	// Split the time string and get the hour part
	hourStr := strings.Split(time, ":")[0]
	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		// Handle conversion error
		fmt.Println("Ошибка преобразования времени:", err)
		return ""
	}

	// Determine emoji based on the hour
	switch {
	case hour >= 0 && hour < 4:
		return "🌙"
	case hour >= 4 && hour < 11:
		return "🌅"
	case hour >= 11 && hour < 16:
		return "🌞"
	case hour >= 16 && hour < 19:
		return "🌇"
	case hour >= 19:
		return "🌌"
	default:
		return ""
	}
}

func getWeatherEmoji(condition string) string {
	condition = strings.ToLower(condition)
	switch {
	case strings.Contains(condition, "ясно"):
		return "☀️"
	case strings.Contains(condition, "облачно"):
		return "🌥"
	case strings.Contains(condition, "дождь"):
		return "🌧"
	case strings.Contains(condition, "снег"):
		return "❄️"
	case strings.Contains(condition, "гроза"):
		return "🌩"
	default:
		return "🌤"
	}
}
