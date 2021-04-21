package notify

import (
	"errors"

	"github.com/TimothyYe/godns"
	"github.com/bwmarrin/discordgo"
)

type DiscordNotify struct {
	conf *godns.Settings
}

func NewDiscordNotify(conf *godns.Settings) INotify {
	return &DiscordNotify{conf: conf}
}

func (n *DiscordNotify) Send(domain, currentIP string) error {

	if n.conf.Notify.Discord.BotApiToken == "" {
		return errors.New("bot api token cannot be empty")
	}

	if n.conf.Notify.Discord.Channel == "" {
		return errors.New("channel id cannot be empty")
	}

	tpl := n.conf.Notify.Discord.MsgTemplate
	if tpl == "" {
		tpl = "Your IP address for {{.Domain}} has been updated to {{ .CurrentIP }} "
	}
	msg := buildTemplate(currentIP, domain, tpl)

	//Create discordgo client
	d, err := discordgo.New("Bot " + n.conf.Notify.Discord.BotApiToken)
	if err != nil {
		return errors.New("error creating discord bot")
	}
	//Open socket connection
	err = d.Open()
	if err != nil {
		return errors.New("error opening connection,")
	}
	//Send message
	_, err = d.ChannelMessageSend(n.conf.Notify.Discord.Channel, msg)
	if err != nil {
		return errors.New("error sending message")
	}
	//Close socket connection
	err = d.Close()
	if err != nil {
		return errors.New("error closing discord connection")
	}
	return nil
}
