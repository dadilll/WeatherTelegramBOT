package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bregydoc/gtranslate"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type WeatherData struct {
	Condition     string
	Temperature   string
	WindSpeed     string
	Precipitation string
	Pressure      string
	Visibility    string
	Cloudiness    string
}

func main() {
	bot, err := tgbotapi.NewBotAPI("6699865318:AAFZeuStbL37m07Qmod0iguI9H1jZlIYYU8")
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
				handleWeatherCommands(bot, update.Message)
			}
		}
	}
}

func handleWeatherCommands(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if strings.HasPrefix(message.Text, "/today") {
		location := strings.TrimSpace(strings.TrimPrefix(message.Text, "/today"))
		weatherText, err := getTodayWeather(location)
		if err != nil {
			log.Printf("Ошибка при получении погоды: %v", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при получении погоды.")
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, weatherText)
		msg.ParseMode = tgbotapi.ModeHTML // Use HTML mode for rich formatting
		bot.Send(msg)
	} else if strings.HasPrefix(message.Text, "/week") {
		location := strings.TrimSpace(strings.TrimPrefix(message.Text, "/week"))
		forecastText, err := getWeekWeather(location)
		if err != nil {
			log.Printf("Ошибка при получении прогноза: %v", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при получении прогноза.")
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, forecastText)
		msg.ParseMode = tgbotapi.ModeHTML // Use HTML mode for rich formatting
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда. Попробуй /start или /help.")
		bot.Send(msg)
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
		return "", errors.New("название города не указано в запросе")
	}

	translatedLocation, err := translateToEnglish(location)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://www.ventusky.com/ru/%s", strings.ReplaceAll(strings.ToLower(translatedLocation), " ", "-"))
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

	var weatherData WeatherData

	// Extract weather data
	weatherData.Condition = strings.TrimSpace(doc.Find(".info_tables .big_ico_cell img").AttrOr("title", ""))
	weatherData.Temperature = strings.TrimSpace(doc.Find(".info_tables .temperature").Text())
	weatherData.WindSpeed = strings.TrimSpace(doc.Find(".info_tables .wind_ico").Text())
	weatherData.Precipitation = strings.TrimSpace(doc.Find(".info_tables tr:contains('Осадки') b").Text())
	weatherData.Pressure = strings.TrimSpace(doc.Find(".info_tables tr:contains('Атмосферное давление') b").Text())
	weatherData.Visibility = strings.TrimSpace(doc.Find(".info_tables tr:contains('Видимость') b").Text())
	weatherData.Cloudiness = strings.TrimSpace(doc.Find(".info_tables tr:contains('Облачность') b").Text())

	// Format output as Markdown with emojis
	var weatherDetails strings.Builder
	weatherDetails.WriteString("*Погода сегодня:*\n")
	weatherDetails.WriteString(fmt.Sprintf("- *Состояние:* %s %s\n", getWeatherEmoji(weatherData.Condition), weatherData.Condition))
	weatherDetails.WriteString(fmt.Sprintf("- *Температура:* 🌡️ %s\n", weatherData.Temperature))
	weatherDetails.WriteString(fmt.Sprintf("- *Скорость ветра:* 💨 %s\n", weatherData.WindSpeed))
	weatherDetails.WriteString(fmt.Sprintf("- *Осадки:* %s %s\n", getEmojiForPrecipitation(weatherData.Precipitation), weatherData.Precipitation))
	weatherDetails.WriteString(fmt.Sprintf("- *Атмосферное давление:* 🌬️ %s\n", weatherData.Pressure))
	weatherDetails.WriteString(fmt.Sprintf("- *Видимость:* 👁️ %s\n", weatherData.Visibility))
	weatherDetails.WriteString(fmt.Sprintf("- *Облачность:* %s %s\n", getEmojiForCloudiness(weatherData.Cloudiness), weatherData.Cloudiness))

	return weatherDetails.String(), nil
}

func getWeekWeather(location string) (string, error) {
	if location == "" {
		return "", fmt.Errorf("название города не указано в запросе")
	}

	translatedLocation, err := translateToEnglish(location)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://www.ventusky.com/ru/%s", strings.ReplaceAll(strings.ToLower(translatedLocation), " ", "-"))
	log.Printf("Отправка запроса на URL: %s", url)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("ошибка запроса: %d %s", res.StatusCode, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	var weatherDetails strings.Builder

	doc.Find("ul.menu li").Each(func(i int, s *goquery.Selection) {
		day := strings.TrimSpace(s.Find("a").Contents().First().Text())
		condition, _ := s.Find("img").Attr("title")
		temperature := strings.TrimSpace(s.Find("a").Contents().Last().Text())

		if day != "" && condition != "" && temperature != "" {
			weatherDetails.WriteString(fmt.Sprintf("<b>День: %s</b>\n", day))
			weatherDetails.WriteString(fmt.Sprintf("<b>Погода:</b> %s %s\n", getWeatherEmoji(condition), condition))
			weatherDetails.WriteString(fmt.Sprintf("<b>Температура:</b> %s\n\n", temperature))
		}
	})

	if weatherDetails.Len() == 0 {
		return "", fmt.Errorf("не удалось получить прогноз погоды. Попробуйте еще раз с другим названием.")
	}

	return weatherDetails.String(), nil
}

func getWeatherEmoji(condition string) string {
	condition = strings.ToLower(condition)
	switch {
	case strings.Contains(condition, "ясно"):
		return "🌞"
	case strings.Contains(condition, "облачно"):
		return "☁️"
	case strings.Contains(condition, "дождь"):
		return "🌧️"
	case strings.Contains(condition, "снег"):
		return "❄️"
	case strings.Contains(condition, "гроза"):
		return "⛈️"
	case strings.Contains(condition, "туман"):
		return "🌫️"
	case strings.Contains(condition, "чистое небо"):
		return "☀️"
	case strings.Contains(condition, "смешанный с дождевыми дождями"):
		return "🌦️"
	default:
		return "🌤️" // Для всех остальных случаев
	}
}

// Функция getEmojiForPrecipitation возвращает эмодзи для осадков
func getEmojiForPrecipitation(precipitation string) string {
	precipitation = strings.ToLower(precipitation)
	precipitationValue, err := strconv.Atoi(strings.TrimSuffix(precipitation, " mm"))
	if err != nil {
		return "" // Вернуть пустую строку, если не удается преобразовать в число
	}

	switch {
	case precipitationValue > 20:
		return "🌧️" // Сильный дождь
	case precipitationValue > 5:
		return "🌦️" // Легкий дождь
	case precipitationValue > 0:
		return "🌂" // Небольшие осадки
	default:
		return "" // Без осадков
	}
}

// Функция getEmojiForCloudiness возвращает эмодзи для облачности на основе процента облачности
func getEmojiForCloudiness(cloudiness string) string {
	cloudiness = strings.ToLower(cloudiness)
	cloudinessValue, err := strconv.Atoi(strings.TrimSuffix(cloudiness, " %"))
	if err != nil {
		return "" // Вернуть пустую строку, если не удается преобразовать в число
	}

	switch {
	case cloudinessValue > 75:
		return "☁️" // Пасмурно
	case cloudinessValue > 50:
		return "🌥️" // Облачно
	case cloudinessValue > 25:
		return "🌤️" // Переменная облачность
	default:
		return "☀️" // Ясно
	}
}
