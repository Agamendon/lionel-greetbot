package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload" // this will register environment variables from .env
)

var messi_img []string = []string{
	"https://imagevars.gulfnews.com/2019/03/18/190318_Messi_resources1_16a30b3710b_large.jpg",
	"https://pbs.twimg.com/media/C2YBGSNUoAAMwQE?format=jpg",
	"https://www.the-sun.com/wp-content/uploads/sites/6/2022/08/NINTCHDBPICT000750922632.jpg?strip=all&quality=100&w=1920&h=1080&crop=1",
	"https://a4.espncdn.com/combiner/i?img=%2Fphoto%2F2016%2F0706%2Fr100563_1296x729_16%2D9.jpg",
	"https://www.thesun.co.uk/wp-content/uploads/2022/08/NINTCHDBPICT000750782122.jpg",
	"https://pbs.twimg.com/media/E8b1ubkWEAkLegz?format=jpg",
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	guildid := m.GuildID
	log.Printf("GuildMemberAdd triggered, GuildID %s\n", guildid)

	guild, err := s.Guild(guildid)
	if err != nil {
		log.Println("Failed to find guild: ", err)
		return
	}

	sysid := guild.SystemChannelID
	message := discordgo.MessageSend{
		Content: fmt.Sprintf("%s joined! Time to greet them with a Messi pic", m.Mention()),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Do it™",
						Style:    discordgo.PrimaryButton,
						Disabled: false,
						CustomID: "did_it", // doesn't matter in this implementation
					},
				},
			},
		},
	}
	_, err = s.ChannelMessageSendComplex(sysid, &message)
	if err != nil {
		log.Println("Failed to send new member message: ", err)
	}

}

func interactionRespond(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Responding to button push, GuidID %s\n", i.GuildID)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s and Lionel Messi are happy to see you!", i.Member.Mention()),
			Embeds: []*discordgo.MessageEmbed{
				{
					Image: &discordgo.MessageEmbedImage{
						URL: messi_img[rand.Intn(len(messi_img))],
					},
				},
			},
		},
	})
	if err != nil {
		log.Println("Failed to send image: ", err)
	}
}

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Println("TOKEN not found in .env, trying to extract environment variable")
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Failed to create session: ", err, ", exiting")
	}

	discord.AddHandler(guildMemberAdd)
	discord.AddHandler(interactionRespond)

	discord.Identify.Intents = discordgo.IntentsGuildMembers

	err = discord.Open()
	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}
	defer discord.Close()

	log.Println("App is running. Press CTRL-C or send SIGINT to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}
