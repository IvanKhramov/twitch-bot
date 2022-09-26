/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/spf13/cobra"
)

type templ struct {
	Avg_rating     float64
	Total_votes    int
	Top_vote       int
	Count_top_vote int
}

var (
	usersVotes   = make(map[string]int)
	votesCount   = make(map[int]int)
	ratingStatus = 0
)

func checkAccess(msg twitch.PrivateMessage) bool {
	status := false
	for _, name := range cfg.AllowedUserList {
		if msg.User.Name == strings.ToLower(name) {
			status = true
		}
	}
	return status
}

func checkCommand(msg twitch.PrivateMessage, client *twitch.Client) {
	switch strings.ToLower(msg.Message) {
	case "!startrating":
		startRating(client)
	case "!endrating":
		endRating(client)
	}

}

func startRating(client *twitch.Client) {
	client.Say(cfg.ChannelName, cfg.StartText)
	ratingStatus = 1
}

func endRating(client *twitch.Client) {
	var totalSum int
	var t templ
	var msg bytes.Buffer
	tmp := template.New("simple")
	for _, value := range usersVotes {
		totalSum += value
		votesCount[value] += 1
	}
	for vote, count := range votesCount {
		if count > t.Count_top_vote {
			t.Top_vote, t.Count_top_vote = vote, count
		}
	}
	t.Total_votes = len(usersVotes)
	t.Avg_rating = float64(totalSum) / float64(t.Total_votes)
	if math.IsNaN(t.Avg_rating) {
		t.Avg_rating = 0.0
	}
	tmp, err := tmp.Parse(cfg.EndText)
	if err != nil {
		log.Fatal(err)
	}

	err = tmp.Execute(&msg, t)
	if err != nil {
		log.Fatal(err)
	}
	client.Say(cfg.ChannelName, msg.String())
	ratingStatus = 0
	usersVotes = make(map[string]int)
	votesCount = make(map[int]int)
}

func checkRatingMsg(message twitch.PrivateMessage) {
	msg, _ := strconv.Atoi(message.Message)
	if msg > 0 && msg < 11 {
		usersVotes[message.User.Name] = msg
	}
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "This command will start bot",
	Long:  `This command will start bot at selected channel`,
	Run: func(cmd *cobra.Command, args []string) {

		client := twitch.NewClient(cfg.ChannelName, cfg.Token)

		client.OnPrivateMessage(func(message twitch.PrivateMessage) {
			if ratingStatus == 1 {
				checkRatingMsg(message)
			}

			if checkAccess(message) {
				checkCommand(message, client)
			}

		})

		client.Join(cfg.ChannelName)
		fmt.Println("Connected to channel -", cfg.ChannelName)
		fmt.Println("Allowed users to run command -", cfg.AllowedUserList)
		err := client.Connect()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&Config, "config", "c", "config.yml", "Path to config file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
