package main

import (
	"fmt"
	"strings"
	"testing"
)

// Примерная функция перевода местоположений
func tsranslateToEnglish(location string) (string, error) {
	translations := map[string]string{
		"Волгоград":    "Volgograd",
		"Арбат":        "Arbat",
		"Екатеринбург": "Yekaterinburg",
		"Самара":       "Samara",
	}
	if translated, ok := translations[location]; ok {
		return translated, nil
	}
	return "", fmt.Errorf("не удалось перевести местоположение: %s", location)
}

// Моковая функция для получения погоды на сегодня
func WOKgetTodayWeather(location string) (string, error) {
	return fmt.Sprintf("Погода в %s сегодня по часам\n\nТемпература: 20°C", location), nil
}

// Моковая функция для получения погоды на неделю
func WOKgetWeekWeather(location string) (string, error) {
	return fmt.Sprintf("Прогноз на неделю для %s\n\nТемпература: 20°C", location), nil
}

func TestGetTodayWeather(t *testing.T) {
	locations := []string{"Волгоград", "Арбат", "Екатеринбург"}

	for _, location := range locations {
		t.Run(location, func(t *testing.T) {
			translatedLocation, err := translateToEnglish(location)
			if err != nil {
				t.Fatalf("Ошибка при переводе местоположения: %v", err)
			}

			url := fmt.Sprintf("https://www.meteoservice.ru/weather/today/%s", strings.ReplaceAll(strings.ToLower(translatedLocation), " ", "-"))
			t.Logf("'%s' переведено на английский как '%s'", location, translatedLocation)
			t.Logf("Запрос на сайт: %s", url)

			weatherText, err := getTodayWeather(location)
			if err != nil {
				t.Fatalf("Ошибка при получении погоды: %v", err)
			}

			if weatherText == "" {
				t.Error("Ожидаемый прогноз не должен быть пустым")
			}

			if !strings.Contains(weatherText, "Температура") {
				t.Errorf("Температура не найдена в выводе для %s: %s", location, weatherText)
			}

			t.Logf("Прогноз на сегодня для %s:\n%s", location, weatherText)
		})
	}
}

func TestGetWeekWeather(t *testing.T) {
	location := "Самара"

	translatedLocation, err := translateToEnglish(location)
	if err != nil {
		t.Fatalf("Ошибка при переводе местоположения: %v", err)
	}

	weatherText, err := getWeekWeather(translatedLocation)
	if err != nil {
		t.Fatalf("Ошибка при получении погоды: %v", err)
	}

	if weatherText == "" {
		t.Error("Ожидаемый прогноз не должен быть пустым")
	}

	t.Logf("Прогноз на неделю для %s:\n%s", location, weatherText)
}
