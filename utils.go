package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.ToLower(b) == strings.ToLower(a) {
			return true
		}
	}
	return false
}

// GetRole returns a Role using session, guild id and the name of the wanted role.
func GetRole(session *discordgo.Session, guildID, roleName string) *discordgo.Role {
	guild, err := session.Guild(guildID)

	if err != nil {
		log.WithFields(log.Fields{
			"err":     err,
			"guildID": guildID,
		}).Error("Cannot get guild")
		return nil
	}

	var role *discordgo.Role

	for i := range guild.Roles {
		role = guild.Roles[i]
		if strings.ToLower(role.Name) == strings.ToLower(roleName) {
			break
		}
		role = nil
	}

	return role
}

// GetRoles gets an array of Roles based on the role names given in a string array.
func GetRoles(session *discordgo.Session, guildID string, roleNames []string) []*discordgo.Role {
	var roles []*discordgo.Role

	for _, roleName := range roleNames {
		role := GetRole(session, guildID, roleName)
		if role != nil {
			roles = append(roles, role)
		}
	}

	return roles
}

// UserHasRole checks if a Guild Member has one of a Role in an array of Roles.
func UserHasRole(member *discordgo.Member, roles []*discordgo.Role) bool {
	for _, validRole := range roles {
		for _, memberRole := range member.Roles {
			if validRole.ID == memberRole {
				return true
			}
		}
	}
	return false
}
