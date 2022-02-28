package main

import (
	"fmt"
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

	log.Println(fmt.Sprintf("I will delete %s message.", HandlerList))

	dab := func(c tb.Context, ban bool) error {
		myRights, _ := b.ChatMemberOf(c.Chat(), b.Me)
		if !myRights.Rights.CanDeleteMessages || !myRights.Rights.CanRestrictMembers {
			_ = c.Send("爷权限不足，告辞！")
			return b.Leave(c.Chat())
		}
		if ban && c.Sender().ID == 777000 {
			return nil
		}
		err := c.Delete()
		if ban {
			return b.BanSenderChat(c.Chat(), c.Sender())
		}
		return err
	}

	for i := 0; i < len(HandlerList); i++ {
		OnEvent := HandlerList[i]
		b.Handle(OnEvent, func(c tb.Context) error {
			return dab(c, false)
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
		M, V, S, W := upd.Message != nil, false, false, true
		if M {
			S, V = upd.Message.SenderChat != nil, upd.Message.Via != nil
			if V {
				if _, ok := ViaWL[upd.Message.Via.Username]; ok {
					W = false
				}
			}
		}
		switch {
		case M && S && DeleteChannel:
			_ = dab(b.NewContext(upd), true)
		case M && V && DeleteVia && W:
			_ = dab(b.NewContext(upd), false)
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
