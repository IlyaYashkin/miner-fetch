package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot"
	"os"
	"sync"
)

type TelegramSender struct {
	chats map[string]int64
	mu    sync.Mutex
	b     *bot.Bot
}

func NewTelegramSender() *TelegramSender {
	return &TelegramSender{
		chats: make(map[string]int64),
	}
}

func (t *TelegramSender) Init(b *bot.Bot) error {
	t.b = b

	err := t.loadChatsFromFile()

	return err
}

func (t *TelegramSender) SaveChatId(user string, chatId int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	savedChatId, ok := t.chats[user]
	if ok && savedChatId == chatId {
		return nil
	}

	t.chats[user] = chatId

	err := t.saveChatsToFile()

	return err
}

func (t *TelegramSender) GetChatIds() map[string]int64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.chats
}

func (t *TelegramSender) SendMessage(ctx context.Context, chatID int64, nodeName string, message string) error {
	params := bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf(" 🔌 %s\n\n%s", nodeName, message),
	}

	_, err := t.b.SendMessage(ctx, &params)

	return err
}

func (t *TelegramSender) saveChatsToFile() error {
	data, err := json.Marshal(t.chats)
	if err != nil {
		return err
	}
	return os.WriteFile("mf-chat-ids.json", data, 0644)
}

func (t *TelegramSender) loadChatsFromFile() error {
	_, err := os.Stat("mf-chat-ids.json")
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	data, err := os.ReadFile("mf-chat-ids.json")
	if err != nil {
		return err
	}

	var chats map[string]int64

	if err := json.Unmarshal(data, &chats); err != nil {
		return err
	}

	t.chats = chats

	return nil
}
