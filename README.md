# Telegram OCR Bot

Телеграм-бот на Go для распознавания текста на изображениях с помощью Tesseract OCR. Поддерживает выбор языка (русский / английский).

## Возможности

- Распознавание текста с изображений (фото).
- Выбор языка распознавания: `русский`, `английский`, `русский + английский`, `английский + русский`.
- Использует Tesseract.
- Собирается и запускается в Docker.

## Команды бота

- `/start` — приветственное сообщение.
- `/lang` — выбор языка распознавания.

## Используемые технологии

- Go
- [gotgbot](https://github.com/PaulSonOfLars/gotgbot) — Telegram Bot API клиент
- Tesseract OCR
- Docker

## Запуск без Docker

1. Установи [Tesseract](https://tesseract-ocr.github.io/tessdoc/Installation.html).
2. Установи Go 1.24.4+.
3. Создай `.env` файл:
4. Запусти через `make run`

## Сборка и запуск через Docker

```bash
docker build -t telegram-ocr-bot .
docker run -e TELEGRAM_TOKEN=your_bot_token_here telegram-ocr-bot
```
