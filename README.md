# Telegram Weather Bot

## Overview
This project is a Telegram bot that provides weather information for various cities. Users can get today's weather or a weekly forecast by using specific commands. The bot fetches weather data from the [Ventusky](https://my.ventusky.com/about/)
 website, translates city names from Russian to English, and formats the response with emojis for better readability.

## Features
- Provides today's weather using the /today <city> command.
- Provides a weekly weather forecast using the /week <city> command.
- Uses emojis to represent different weather conditions for better user experience.

## Dependencies
- [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery) for HTML parsing.
- [bregydoc/gtranslate](https://github.com/bregydoc/gtranslate) for translating city names.
- [gopkg.in/telegram-bot-api.v4](https://gopkg.in/telegram-bot-api.v4/) for interacting with the Telegram Bot API.

## Installation

1. Clone the repository:
```bash
git clone https://github.com/dadilll/WeatherTelegramBot.git
cd WeatherTelegramBot
```
2. Install the required Go packages:
```bash
go get -u github.com/PuerkitoBio/goquery
go get -u github.com/bregydoc/gtranslate
go get -u gopkg.in/telegram-bot-api.v4
```

3. Run the program:
```bash
go build -o WeatherTelegramBot
./WeatherTelegramBot
```

## Usage

### Commands
- /start - Initializes the bot and provides a welcome message.
- /help - Provides instructions on how to use the bot.
- /today <city> - Fetches and returns the weather for the specified city today.
- /week <city> - Fetches and returns the weather forecast for the specified city for the upcoming week.

### Example
1. Get today's weather:

```md
*Погода сегодня:*
  Состояние: 🌞 Ясно
  Температура: 25°C 🌡️
  Скорость ветра: 5 м/с 💨
  Осадки: 0 mm 🌂
  Атмосферное давление: 1015 hPa 🌬️
  Видимость: 10 км 👁️
  Облачность: 20 % 🌤️
  Влажность: 45 % 🌳
```

2. Get today's weather:

```md
*Погода сегодня:*
<b>День: Понедельник</b>
<b>Погода:</b> ☀️ Ясно
<b>Температура:</b> 25°C

<b>День: Вторник</b>
<b>Погода:</b> 🌦️ Дождь
<b>Температура:</b> 22°C
```