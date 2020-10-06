package main

import (
	"fmt"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var HandlerList []string

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  TBToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	Config(DeleteDice, tb.OnDice)
	Config(DeleteUserJoined, tb.OnUserJoined)
	Config(DeleteUserLeft, tb.OnUserLeft)
	Config(DeleteDNewGroupTitle, tb.OnNewGroupTitle)
	Config(DeleteNewGroupPhoto, tb.OnNewGroupPhoto)
	Config(DeleteGroupPhotoDeleted, tb.OnGroupPhotoDeleted)
	Config(DeleteOnPinned, tb.OnPinned)

	log.Println(fmt.Sprintf("I will delete %s message.", HandlerList))

	for i := 0; i < len(HandlerList); i++ {
		OnEvent := HandlerList[i]
		b.Handle(OnEvent, func(m *tb.Message) {
			b.Delete(m)
			myRights, _ := b.ChatMemberOf(m.Chat, b.Me)
			if !myRights.Rights.CanDeleteMessages {
				b.Send(m.Chat, "爷没有删除消息权限，告辞！")
				b.Leave(m.Chat)
			}
		})
	}

	b.Start()
}

func Config(check bool, handler string) {
	if check {
		HandlerList = append(HandlerList, handler)
	}
}
