package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetTodayWeather(t *testing.T) {
	locations := []string{"Красноярск", "Волгоград"}

	for _, location := range locations {
		t.Run(location, func(t *testing.T) {
			translatedLocation, err := translateToEnglish(location)
			if err != nil {
				t.Fatalf("Ошибка при переводе местоположения: %v", err)
			}

			url := fmt.Sprintf("https://www.ventusky.com/ru/%s", strings.ReplaceAll(strings.ToLower(translatedLocation), " ", "-"))
			t.Logf("'%s' переведено на английский как '%s'", location, translatedLocation)
			t.Logf("Запрос на сайт: %s", url)

			weatherText, err := getTodayWeather(location)
			if err != nil {
				t.Fatalf("Ошибка при получении погоды: %v", err)
			}

			if weatherText == "" {
				t.Error("Ожидаемый прогноз не должен быть пустым")
			}

			t.Logf("Прогноз на сегодня для %s:\n%s", location, weatherText)
		})
	}
}

func TestGetWeekWeather(t *testing.T) {
	locations := []string{"Самара", "Екатеринбург"}

	for _, location := range locations {
		t.Run(location, func(t *testing.T) {
			translatedLocation, err := translateToEnglish(location)
			if err != nil {
				t.Fatalf("Ошибка при переводе местоположения: %v", err)
			}
			url := fmt.Sprintf("https://www.ventusky.com/ru/%s", strings.ReplaceAll(strings.ToLower(translatedLocation), " ", "-"))
			t.Logf("'%s' переведено на английский как '%s'", location, translatedLocation)
			t.Logf("Запрос на сайт: %s", url)

			weatherText, err := getWeekWeather(location)
			if err != nil {
				t.Fatalf("Ошибка при получении прогноза: %v", err)
			}

			if weatherText == "" {
				t.Error("Ожидаемый прогноз не должен быть пустым")
			}

			t.Logf("Прогноз на сегодня для %s:\n%s", location, weatherText)
		})
	}
}
