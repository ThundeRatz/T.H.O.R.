// Package dclient is a wrapper around DiscordGo library
package dclient

import (
	"fmt"
	"strings"

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
	Logger     *zerolog.Logger
	SettingsDB *bbolt.DB
	LevelCache map[string]int // Vem da config
	Config     []string

	aliases  map[string]string
	commands map[string]*Command
	session  *discordgo.Session
	token    string
}

// New creates a new client
func New() *Client {
	return &Client{
		aliases:  make(map[string]string),
		commands: make(map[string]*Command),
	}
}

// Init initializes a client with a provided token and logger
func (c *Client) Init(token string, logger *zerolog.Logger) {
	c.Logger = logger
	c.token = token

	c.session, _ = discordgo.New(fmt.Sprintf("Bot %s", c.token))
	c.session.AddHandler(c.OnMessageCreate)

	err := c.session.Open()
	if err != nil && c.Logger != nil {
		c.Logger.Error().Err(err).Msg("error opening connection to Discord")
	}
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

	// Get permLevel for mc c.permLevel(mc)
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

	// Checa se o level é maior que que level do comando

	if c.Logger != nil {
		c.Logger.Debug().
			Str("user", mc.Author.Username).
			Str("user_id", mc.Author.ID).
			Str("cmd", command).
			Msg("Running command")
	}

	cmd.Run(&Context{
		Content: strings.TrimSpace(mc.Content),
		Command: command,
		Args:    args,
		Session: ds,
		Message: mc.Message,
	})

	// StopTyping
}

// AddCommand adds a new command to the client
func (c *Client) AddCommand(cmd *Command) {
	if c.Logger != nil {
		c.Logger.Debug().Str("cmd", cmd.Name).Msg("Adding Command")
	}

	c.commands[cmd.Name] = cmd

	for _, a := range cmd.Aliases {
		c.aliases[a] = cmd.Name
	}
}

func (c *Client) permLevel(*discordgo.Message) int {
	// sort.Slice(c.Config, func(i, j int) bool { c.Config[i] < c.Config[j] })
	return 0
}

// Start runs the bot
func (c *Client) Start() {

}

// Stop stops the bot
func (c *Client) Stop() {
	if c.session != nil {
		c.session.Close()
	}
}
