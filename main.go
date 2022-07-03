package main

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvFile()

	TOKEN := os.Getenv("DISCORD_BOT_TOKEN")
	if TOKEN == "" {
		log.Fatal("Error loading DISCORD_BOT_TOKEN")
	}

	discord, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal(err.Error())
	}

	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		log.Fatal("Error opening connection, ", err.Error())
		return
	}

	log.Print("Logged in as " + discord.State.User.Username)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Println("Press CTRL-C to exit.")
	<-sc

	discord.Close()
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	r := regexp.MustCompile(`[\w\W]+\.(com|net|org|jp|dev|info|xyz|tokyo|me|link|club|click|space|cc|in|tv|style|work|)`)

	if message.Author.Bot {
		return
	}

  /*
  if message.Content == "ping" {
    session.ChannelMessageSend(message.ChannelID, "pong!")
  }
  */

	domains := r.FindAllString(message.Content, -1)

	if len(domains) > 0 {
    whois(session, message, domains)
	}

}

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
