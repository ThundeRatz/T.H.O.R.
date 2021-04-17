// Package dclient is a wrapper around DiscordGo library
package dclient

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"go.etcd.io/bbolt"
)

// Context holds extra data to be passed along to command handlers
type Context struct {
	Content         string
	Command         string
	Args            []string
	Settings        map[string]string // Prefix
	AuthorPermLevel string
	Logger          zerolog.Logger

	Session *discordgo.Session
	Message *discordgo.Message
}

// Command holds information about a command
type Command struct {
	Run         func(*Context)
	Name        string
	Aliases     []string
	Enabled     bool
	GuildOnly   bool
	PermLevel   string
	Category    string
	Description string
	Usage       string
}

// Client is the main struct
type Client struct {
	SettingsDB *bbolt.DB
	LevelCache map[string]int // Vem da config
	Config     []string

	logger     zerolog.Logger
	aliases    map[string]string
	commands   map[string]*Command
	permLevels map[string]int
	session    *discordgo.Session
	token      string
}

// Init initializes a client with a provided token and logger
func (c *Client) Init(token string, logger *zerolog.Logger) error {
	c.logger = logger.With().Str("pkg", "dclient").Logger()
	c.token = token
	c.aliases = make(map[string]string)
	c.commands = make(map[string]*Command)
	c.permLevels = map[string]int{
		"User":   0,
		"Gestão": 9,
		"Admin":  10,
		"Owner":  99,
	}

	c.logger.Debug().Msg("Initializing discord client")

	for _, v := range c.commands {
		c.logger.Debug().Str("cmd", v.Name).Msg("Loaded command")
	}

	c.session, _ = discordgo.New(fmt.Sprintf("Bot %s", c.token))
	c.session.AddHandler(c.OnMessageCreate)
	c.session.AddHandler(c.OnReady)

	err := c.session.Open()
	if err != nil {
		c.logger.Error().Err(err).Msg("error opening connection to Discord")
		return err
	}

	return nil
}

// OnMessageCreate is a DiscordGo Event Handler function.  This must be
// registered using the DiscordGo.Session.AddHandler function.  This function
// will receive all Discord messages and parse them for matches to registered
// commands.
func (c *Client) OnMessageCreate(ds *discordgo.Session, mc *discordgo.MessageCreate) {
	if mc.Author.Bot {
		return
	}

	// Get settings
	prefix := "." // Get prefix from settings

	if !strings.HasPrefix(mc.Content, prefix) {
		return
	}

	args := strings.Fields(mc.Content)

	var command string
	command, args = strings.ToLower(strings.TrimPrefix(args[0], prefix)), args[1:]

	cmd, ok := c.commands[command]

	if !ok {
		cmd, ok = c.commands[c.aliases[command]]
	}

	if !ok || !cmd.Enabled {
		ds.ChannelMessageSend(mc.ChannelID, fmt.Sprintf("Não consegui encontrar o comando `%s`. Use `%shelp` para ver os comandos disponíveis", command, prefix))
		return
	}

	if mc.GuildID == "" && cmd.GuildOnly {
		ds.ChannelMessageSend(mc.ChannelID, "Esse comando não está disponível em mensagens privadas. Use-o numa guilda.")
		return
	}

	permLevel := "User"

	// For now, admin IDs are hardcoded
	if mc.Author.ID == "232163710506893312" {
		permLevel = "Owner"
	}

	if mc.Author.ID == "369989751974920201" {
		permLevel = "Admin"
	}

	if c.permLevels[permLevel] < c.permLevels[cmd.PermLevel] {
		ds.ChannelMessageSend(mc.ChannelID, "Você não tem permissão para usar esse comando!")
		return
	}

	c.logger.Debug().
		Str("user", mc.Author.Username).
		Str("user_id", mc.Author.ID).
		Str("command", command).
		Msg("Running command")

	cmd.Run(&Context{
		Content:         strings.TrimSpace(mc.Content),
		Command:         command,
		Args:            args,
		Session:         ds,
		Message:         mc.Message,
		Logger:          c.logger.With().Str("cmd", command).Logger(),
		AuthorPermLevel: permLevel,
	})

	// StopTyping
}

// OnReady is a DiscordGo Event Handler function. This must be
// registered using the DiscordGo.Session.AddHandler function.  This function
// will be called once the client is ready.
func (c *Client) OnReady(ds *discordgo.Session, r *discordgo.Ready) {
	c.logger.Info().Msg("Discord client ready")

	err := ds.UpdateStreamingStatus(0, "WCXV", "https://youtube.com/watch?v=XZHTBcfhZwg")

	if err != nil {
		c.logger.Error().Err(err)
	}
}

// SetMessageReactionAddHandler allows the addition of a custom
// handler for MessageReactions
func (c *Client) SetMessageReactionAddHandler(f func(ds *discordgo.Session, r *discordgo.MessageReactionAdd)) {
	c.session.AddHandler(f)
}

// SetGuildMemberAddHandler allows the addition of a custom
// handler for GuildMemberAdd
func (c *Client) SetGuildMemberAddHandler(f func(ds *discordgo.Session, r *discordgo.GuildMemberAdd)) {
	c.session.AddHandler(f)
}

// AddCommand adds a new command to the client
func (c *Client) AddCommand(cmd *Command) {
	c.logger.Debug().Str("cmd", cmd.Name).Msg("Adding command")

	c.commands[cmd.Name] = cmd

	for _, a := range cmd.Aliases {
		c.aliases[a] = cmd.Name
	}
}

// Help is a built-in help command
func (c *Client) Help(ctx *Context) {
	if len(ctx.Args) == 0 {

		embed := &discordgo.MessageEmbed{
			Color: 0xe800ff,
			Author: &discordgo.MessageEmbedAuthor{
				Name: "Ajuda",
				URL:  "https://thunderatz.org",
			},
			Description: "Aqui você pode ver todos os comandos que eu conheço.",
			Timestamp:   time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "T.H.O.R | ThundeRatz",
			},
		}

		values := ""
		for _, v := range c.commands {
			if c.permLevels[v.PermLevel] <= c.permLevels[ctx.AuthorPermLevel] {
				values += "`" + v.Name + "`, "
			}
		}
		values = strings.TrimSuffix(values, ", ")

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Commands",
			Value: values,
		})

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "\u200b",
			Value: "**Use `.help <comando>` para ajuda mais específica.**",
		})

		c.session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
	} else {
		cmd := ctx.Args[0]

		if cmd, ok := c.commands[cmd]; ok && cmd.Enabled {
			msg := "```asciidoc\n"
			msg += "= " + cmd.Name + " =\n"
			msg += cmd.Description + "\n"
			msg += "usage:: " + cmd.Usage + "\n"
			msg += "aliases:: " + strings.Join(cmd.Aliases, ", ") + "```"

			c.session.ChannelMessageSend(ctx.Message.ChannelID, msg)
		}
	}
}

// AddHelpCmd adds the built-in help command to the list of commands
func (c *Client) AddHelpCmd() {
	c.AddCommand(&Command{
		Run: c.Help,

		Aliases:     []string{"h"},
		Category:    "Geral",
		Description: "Mostra todos os comandos disponíveis, ou detalhes de um comando específico.",
		Usage:       "help [command]",
		Enabled:     true,
		GuildOnly:   false,
		Name:        "help",
		PermLevel:   "User",
	})
}

// Stop stops the bot
func (c *Client) Stop() {
	if c.session != nil {
		c.logger.Info().Msg("Closing discord client")
		c.session.Close()
	}
}

// SendMessage sends a message to the specified channel
func (c *Client) SendMessage(channelID, content string) (*discordgo.Message, error) {
	return c.session.ChannelMessageSend(channelID, content)
}

// SendEmbed sends an embed message to the specified channel
func (c *Client) SendEmbed(channelID string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.session.ChannelMessageSendEmbed(channelID, embed)
}
