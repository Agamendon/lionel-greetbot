package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload" // this will register environment variables from .env
)

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	log.Println("New event triggered")
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		log.Println("Failed to find guild: ", err)
		return
	}

	sysid := guild.SystemChannelID
	_, err = s.ChannelMessageSend(sysid, "New Lionelian here!")
	if err != nil {
		log.Println("Failed to send new member message: ", err)
	}
}

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN not found in .env, exiting")
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Failed to create session: ", err, ", exiting")
	}

	discord.AddHandler(guildMemberAdd)

	discord.Identify.Intents = discordgo.IntentsGuildMembers

	err = discord.Open()
	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}

	log.Println("App is running. Press CTRL-C or send SIGINT to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()

}
