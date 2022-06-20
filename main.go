package main

import (
	"log"
	"time"

	tb "gopkg.in/telebot.v3"
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

	log.Printf("I will delete %s message.", HandlerList)

	DeleteAndBan := func(c tb.Context, ban ...bool) (err error) {
		myRights, err := b.ChatMemberOf(c.Chat(), b.Me)
		if err != nil {
			return
		}
		if !myRights.Rights.CanDeleteMessages || !myRights.Rights.CanRestrictMembers {
			_ = c.Send("爷权限不足，告辞！")
			return b.Leave(c.Chat())
		}
		if ban != nil && ban[0] && c.Sender().ID == 777000 {
			return
		}

		if c.Message().SenderChat != nil {
			if c.Message().SenderChat.ID == c.Chat().ID {
				return
			}
			if ban != nil && ban[0] {
				_ = b.BanSenderChat(c.Chat(), c.Sender())
			}
		}

		return c.Delete()
	}

	for i := 0; i < len(HandlerList); i++ {
		OnEvent := HandlerList[i]
		b.Handle(OnEvent, func(c tb.Context) error {
			return DeleteAndBan(c)
		})
	}

	stop := make(chan struct{})
	stopConfirm := make(chan struct{})

	go func() {
		b.Poller.Poll(b, b.Updates, stop)
		close(stopConfirm)
	}()

	for {
		upd := <-b.Updates
		S, V, W := false, false, false
		if upd.Message != nil {
			S = upd.Message.SenderChat != nil
			V = upd.Message.Via != nil
			if V {
				_, W = ViaWL[upd.Message.Via.Username]
			}
		}
		switch {
		case S && DeleteChannel:
			_ = DeleteAndBan(b.NewContext(upd), true)
		case V && W && DeleteVia:
			_ = DeleteAndBan(b.NewContext(upd))
		default:
			b.ProcessUpdate(upd)
		}
	}
}

func Config(check bool, handler string) {
	if check {
		HandlerList = append(HandlerList, handler)
	}
}
