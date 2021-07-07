package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/nicklaw5/helix"
	"github.com/spf13/viper"
)

var chatClient = &twitch.Client{}

func main() {
	rand.Seed(time.Now().Unix())

	err := readConfigFile()
	if err != nil {
		log.Fatal("Could not read config file:", err.Error())
	}

	go func() {
		err := configureChatClient()
		if err != nil {
			log.Fatal("Error configuring chat client:", err.Error())
		}
	}()

	searchClient, err := configureSearchClient()
	if err != nil {
		log.Fatal("Error configuring search client:", err.Error())
	}

	// Get banned game IDs once at the top of main to preserve API calls and decrease latency
	bannedGames, err := getBannedGameIDs(searchClient)
	if err != nil {
		log.Fatal("Error getting banned game IDs:", err.Error())
	}

	// Check list of priority streamers
	checkStreamers(searchClient, viper.GetStringSlice("priority_streamers"), bannedGames)

	// TODO Trigger a raid
	chatClient.Say(viper.GetString("twitch_username"), "I'm awake!")

	// Check list of priority streamers
	checkStreamers(searchClient, viper.GetStringSlice("backup_streamers"), bannedGames)

	// TODO Trigger a raid here too
}

func readConfigFile() error {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("toml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

func configureChatClient() error {
	chatClient = twitch.NewClient(viper.GetString("twitch_bot_username"), viper.GetString("twitch_bot_secret"))
	chatClient.Join(viper.GetString("twitch_username"))
	return chatClient.Connect()
}

func configureSearchClient() (*helix.Client, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     viper.GetString("twitch_client_id"),
		ClientSecret: viper.GetString("twitch_client_secret"),
	})
	if err != nil {
		log.Fatal("Error with Twitch auth:", err.Error())
	}

	tokenResponse, err := client.RequestAppAccessToken([]string{"user:read:email"})
	if err != nil {
		log.Fatal("Error requesting app access token:", err.Error())
	}

	// Set the access token on the client
	client.SetAppAccessToken(tokenResponse.Data.AccessToken)

	return client, nil

}

func getBannedGameIDs(client *helix.Client) (*helix.GamesResponse, error) {
	resp, err := client.GetGames(&helix.GamesParams{
		Names: viper.GetStringSlice("games_blacklist"),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func checkStreamers(client *helix.Client, allStreamers []string, bannedGames *helix.GamesResponse) {
	for _, channel := range allStreamers {
		candidate, channelStatus := isChannelRaidCandidate(client, channel, bannedGames)
		if candidate {
			currentGame, err := getGameNameByID(client, channelStatus.GameID)
			if err != nil {
				// TODO Handle the error
				log.Println("Error getting game for ID: " + channelStatus.GameID)
			}

			// TODO Alert in a more useful way (pop up window?)
			fmt.Println("https://twitch.tv/" + channel + " is streaming " + currentGame)
		}
	}
}

func isChannelRaidCandidate(client *helix.Client, channel string, bannedGames *helix.GamesResponse) (bool, helix.Channel) {
	channelStatus, err := getChannelStatus(client, channel)
	if err != nil {
		log.Fatal("Error getting channel info: ", err.Error())
	}

	// Is channel live?
	if !channelStatus.IsLive {
		return false, channelStatus
	}

	/*
		// TODO Is channel in preferred language?
		if !channelStatus.IsLive {
			return false, channelStatus
		}
	*/

	// Is channel streaming a banned game?
	for _, game := range bannedGames.Data.Games {
		if channelStatus.GameID == game.ID {
			return false, channelStatus
		}
	}

	return true, channelStatus
}

func getChannelStatus(client *helix.Client, channel string) (helix.Channel, error) {
	resp, err := client.SearchChannels(&helix.SearchChannelsParams{
		Channel: channel,
		First:   1, // We only want the topmost result (ideally a perfect match)
	})
	if err != nil {
		log.Fatal("Error searching streams:", err.Error())
	}

	if len(resp.Data.Channels) > 0 {
		return resp.Data.Channels[0], nil
	}

	return helix.Channel{}, nil
}

func getGameNameByID(client *helix.Client, gameID string) (string, error) {
	resp, err := client.GetGames(&helix.GamesParams{
		IDs: []string{gameID},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Data.Games) > 0 {
		return resp.Data.Games[0].Name, nil
	}

	return "", errors.New("Could not find name for game with ID " + gameID)
}
