package main

import (
	"fmt"
	"testing"
)

func TestGetTodayWeather(t *testing.T) {
	// Устанавливаем место для теста
	location := "Волгоград"

	// Переводим место на английский
	translatedLocation, err := translateToEnglish(location)
	if err != nil {
		t.Fatalf("Ошибка при переводе местоположения: %v", err)
	}

	// Формируем URL запроса
	url := fmt.Sprintf("https://www.meteoservice.ru/weather/today/%s", translatedLocation)

	// Выводим URL запроса для визуальной проверки
	t.Logf("'%s' переведено на английский как '%s'", location, translatedLocation)
	t.Logf("Запрос на сайт: %s", url)

	// Выполняем запрос и проверяем результат
	weatherText, err := getTodayWeather(translatedLocation)
	if err != nil {
		t.Fatalf("Ошибка при получении погоды: %v", err)
	}

	if weatherText == "" {
		t.Error("Ожидаемый прогноз не должен быть пустым")
	}

	// Вывод полученного прогноза для визуальной проверки
	t.Logf("Прогноз на сегодня для %s:\n%s", location, weatherText)
}

func TestGetWeekWeather(t *testing.T) {
	// Устанавливаем место для теста
	location := "Самара"

	// Переводим место на английский
	translatedLocation, err := translateToEnglish(location)
	if err != nil {
		t.Fatalf("Ошибка при переводе местоположения: %v", err)
	}

	// Выполняем запрос и проверяем результат
	weatherText, err := getWeekWeather(translatedLocation)
	if err != nil {
		t.Fatalf("Ошибка при получении погоды: %v", err)
	}

	if weatherText == "" {
		t.Error("Ожидаемый прогноз не должен быть пустым")
	}

	// Вывод полученного прогноза для визуальной проверки
	t.Logf("Прогноз на неделю для %s:\n%s", location, weatherText)
}
