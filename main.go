package main

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v3"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  TBToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	b.Handle(tb.OnDice, func(m *tb.Message) {
		b.Delete(m)
	})

	b.Start()
}
