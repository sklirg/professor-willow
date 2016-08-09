package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

// ValidTeams is the list of valid team names for role grants using the bot.
var ValidTeams = []string{"INSTINCT", "MYSTIC", "VALOR"}

// HandleTeamJoin is the event handler for joining a team.
func HandleTeamJoin(session *discordgo.Session, m *discordgo.MessageCreate) {
	help := fmt.Sprintf("Usage: `!team <team name>`, example: `!team %v`", strings.ToLower(ValidTeams[rand.Intn(len(ValidTeams))]))
	message := strings.TrimSpace(m.Content)
	words := strings.Split(message, " ")

	if len(words) != 2 {
		session.ChannelMessageSend(m.ChannelID, help)
		return
	}

	team := words[1]

	if !stringInSlice(team, ValidTeams) {
		session.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf("%v\nError: '%v' is not in the valid teams list.", help, team))
		return
	}

	channel, err := session.Channel(m.ChannelID)

	if err != nil {
		log.WithFields(log.Fields{
			"err":        err,
			"channel_id": m.ChannelID,
		}).Error("Cannot get channel")
		return
	}

	role := GetRole(session, channel.GuildID, team)

	if role == nil {
		log.WithFields(log.Fields{
			"role":  team,
			"guild": channel.GuildID,
		}).Error("Could not find role in this guild.")
		session.ChannelMessageSend(channel.ID, fmt.Sprintf("Error: Could not find the role '%v'. Does it exist?", team))
		return
	}

	member, err := session.GuildMember(channel.GuildID, m.Author.ID)

	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"guild_id": channel.GuildID,
			"user_id":  m.Author.ID,
		}).Error("Cannot find member.")
		return
	}

	if UserHasRole(member, GetRoles(session, channel.GuildID, ValidTeams)) {
		session.ChannelMessageSend(channel.ID, fmt.Sprintf("Error: <@%v> has already selected a team.", m.Author.ID))
		return
	}

	newRoles := append(member.Roles, role.ID)
	err = session.GuildMemberEdit(channel.GuildID, m.Author.ID, newRoles)

	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"guild": channel.GuildID,
		}).Error("I don't have permission to Manage Roles on this server!")
		session.ChannelMessageSend(channel.ID, fmt.Sprintf("Error: Could not assign the role \"%v\" to <@%v>. Do I have permissions to Manage Roles?", team, m.Author.ID))
		return
	}

	log.WithFields(log.Fields{
		"user": m.Author.Username,
		"role": role.Name,
	}).Info("Added user to team.")
	session.ChannelMessageSend(channel.ID, fmt.Sprintf("Added %v to '%v'.", m.Author.Username, role.Name))
}

func hi(session *discordgo.Session, m *discordgo.MessageCreate) {
	channel, _ := session.Channel(m.ChannelID)
	guild, _ := session.Guild(channel.GuildID)
	user, _ := session.User(m.Author.ID)
	log.WithFields(log.Fields{
		"user":    user.Username,
		"channel": channel.Name,
		"guild":   guild.Name,
		"message": m.Content,
	}).Debug(m.Content)

	if strings.HasPrefix(m.Content, "!role") || strings.HasPrefix(m.Content, "!team") {
		HandleTeamJoin(session, m)
	}
}

func main() {
	botToken := os.Getenv("BOT_ACCESS_TOKEN")

	if botToken == "" {
		log.Fatal("No access token provided. Please populate BOT_ACCESS_TOKEN.")
		return
	}

	session, err := discordgo.New(botToken)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Login err!")
	}

	session.AddHandler(hi)

	err = session.Open()

	log.Info("Connected.")

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("WS err!", err)
	}

	<-make(chan struct{})
	return
}
