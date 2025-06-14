package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func handleStart(b *gotgbot.Bot, ctx *ext.Context) error {
	text := "Привет! Этот бот позволяет распознать текст в фотографиях.\n\nДля распознавания текста используется tesseract.\n\nПо умолчанию распознается только русский язык, чтобы выбрать другой используйте /lang\n\n*Присылайте по одной фотографии.*"
	_, err := ctx.EffectiveMessage.Reply(b, text, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
	return err
}

func handleLang(b *gotgbot.Bot, ctx *ext.Context) error {
	kb := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				gotgbot.InlineKeyboardButton{Text: "Русский", CallbackData: "lang:rus"},
				gotgbot.InlineKeyboardButton{Text: "Английский", CallbackData: "lang:eng"},
			},
			{gotgbot.InlineKeyboardButton{Text: "Русский + Английский", CallbackData: "lang:rus+eng"}},
			{gotgbot.InlineKeyboardButton{Text: "Английский + Русский", CallbackData: "lang:eng+rus"}},
		},
	}
	text := "Выберите язык для распознавания текста:"
	_, err := ctx.EffectiveMessage.Reply(b, text, &gotgbot.SendMessageOpts{ReplyMarkup: kb})
	return err
}

func handlePhoto(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	userID := msg.From.Id
	lang := "rus"
	if l, ok := userLangs.Load(userID); ok {
		lang = l.(string)
	}

	best := msg.Photo[len(msg.Photo)-1]
	file, err := b.GetFile(best.FileId, nil)
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}
	downloadURL := file.URL(b, nil)
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	text, err := runTesseract(resp.Body, lang)
	if err != nil {
		return fmt.Errorf("failed to recognize text: %w", err)
	}

	_, err = msg.Reply(b, text, nil)
	return err
}

func runTesseract(imageData io.Reader, lang string) (string, error) {
	cmd := exec.Command("tesseract", "-", "stdout", "-l", lang)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	defer stdin.Close()

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if _, err := io.Copy(stdin, imageData); err != nil {
		return "", err
	}
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	return stdout.String(), nil
}

func handleLanguageSelection(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.CallbackQuery
	userID := cb.From.Id
	lang := cb.Data[len("lang:"):]

	userLangs.Store(userID, lang)
	text := fmt.Sprintf("Язык установлен: %s", lang)

	_, _, err := cb.Message.EditText(b, text, nil)
	if err != nil {
		return fmt.Errorf("failed to edit message: %w", err)
	}

	_, err = cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: text})
	return err
}
