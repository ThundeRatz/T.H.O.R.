// Package discord is the discord module for T.H.O.R.
package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"

	"go.thunderatz.org/thor/core/types"
	"go.thunderatz.org/thor/pkg/dclient"
)

var (
	msgCh     types.CoreMsgCh
	serviceId string
)

const (
	emojiReactRoleMember = "😎"
	emojiReactRoleBixe   = "👶"
)

// Service is the main discord service struct
type Service struct {
	AlertChannel string
	Token        string

	client *dclient.Client
	logger zerolog.Logger
}

// Init initializes discord service
func (ds *Service) Init(_logger *zerolog.Logger, _ch types.CoreMsgCh) error {
	ds.logger = _logger.With().Str("serv", "discord").Logger()
	ds.client = &dclient.Client{}
	serviceId = "discord"
	msgCh = _ch

	err := ds.client.Init(ds.Token, &ds.logger)

	if err != nil {
		return err
	}

	ds.client.SetMessageReactionAddHandler(ds.OnMessageReactionAdd)
	ds.client.SetGuildMemberAddHandler(ds.OnMemberAdd)

	ds.client.AddCommand(pingCmd)
	ds.client.AddCommand(infoCmd)
	ds.client.AddCommand(githubCmd)
	ds.client.AddCommand(configCmd)
	ds.client.AddHelpCmd()

	return nil
}

func (ds *Service) OnMessageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	m, err := s.ChannelMessage(r.ChannelID, r.MessageID)

	if err != nil || m == nil || m.Author.ID != s.State.User.ID {
		return
	}

	user, err := s.User(r.UserID)
	if err != nil || user == nil || user.Bot {
		return
	}

	// Ignore messages that are not embeds with a command in the footer
	if len(m.Embeds) != 1 || m.Embeds[0].Footer == nil || m.Embeds[0].Footer.Text == "" {
		return
	}

	args := strings.Split(m.Embeds[0].Footer.Text, ":")
	// Ensure valid footer command
	if len(args) != 2 {
		return
	}

	switch args[0] {
	case "newMember":
		allRoles, err := s.GuildRoles(r.GuildID)

		if err != nil {
			return
		}

		if r.Emoji.ID == emojiReactRoleMember {
			s.GuildMemberRoleAdd(r.GuildID, args[1], getRoleFromSliceByName(allRoles, "Membro").ID)
		}

		if r.Emoji.ID == emojiReactRoleBixe {
			s.GuildMemberRoleAdd(r.GuildID, args[1], getRoleFromSliceByName(allRoles, "Bixe").ID)
		}

		s.MessageReactionsRemoveAll(r.ChannelID, r.MessageID)
	}
}

func (ds *Service) OnMemberAdd(s *discordgo.Session, r *discordgo.GuildMemberAdd) {
	replyCh := make(types.CoreReplyCh)

	msgCh <- types.CoreMsg{
		Type:  types.KVConfigGetMsg,
		Reply: replyCh,
		From:  serviceId,
		Args: types.KVConfigGetArgs{
			Key: "member-add-ch",
		},
	}

	reply := <-replyCh
	replyChannel := reply.Reply.(*types.KVConfigGetReply).Value

	if !reply.Success {
		return
	}

	embed := getBaseEmbed()
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: "newMember:" + r.User.ID,
	}
	embed.Title = "Nova entrada no servidor!"
	embed.Description = fmt.Sprintf(`Reaja para permitir o uso do server

%v: Membro
%v: Bixe`, emojiReactRoleMember, emojiReactRoleBixe)

	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Username",
			Value:  r.User.Username,
			Inline: true,
		},
	}
	embed.Image = &discordgo.MessageEmbedImage{URL: r.User.AvatarURL("128")}

	msg, err := s.ChannelMessageSendEmbed(replyChannel, embed)

	if err != nil {
		return
	}

	s.MessageReactionAdd(msg.ChannelID, msg.ID, emojiReactRoleMember)
	s.MessageReactionAdd(msg.ChannelID, msg.ID, emojiReactRoleBixe)
}

// Stop stops the service
func (ds *Service) Stop() {
	ds.client.Stop()
}

// SendAlert sends a message to the specified alert channel
func (ds *Service) SendAlert(content string) error {
	_, err := ds.client.SendMessage(ds.AlertChannel, content)

	if err != nil {
		ds.logger.Error().
			Err(err).
			Str("channel", ds.AlertChannel).
			Msg("Couldn't send message to channel")
	}

	return err
}

func getBaseEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:     0xe800ff,
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "T.H.O.R. | ThundeRatz",
			IconURL: "https://static.thunderatz.org/ThorJoinha.png",
		},
	}
}

// Helper Functions
func getRoleFromSliceByName(slice []*discordgo.Role, name string) *discordgo.Role {
	for _, v := range slice {
		if v.Name == name {
			return v
		}
	}

	return nil
}
