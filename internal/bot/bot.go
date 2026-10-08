package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func New(token string) *discordgo.Session {
	session, err := discordgo.New(token)

	if err != nil {
		log.Fatalf("Error")
	}

	return session
}
