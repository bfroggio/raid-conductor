package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/nicklaw5/helix"
	"github.com/spf13/viper"
)

func main() {
	client, err := configureClient()
	if err != nil {
		log.Fatal("Error configuring client:", err.Error())
	}

	// Get banned game IDs once at the top of main to preserve API calls and decrease latency
	bannedGames, err := getBannedGameIDs(client)
	if err != nil {
		log.Fatal("Error getting banned game IDs:", err.Error())
	}

	// Check list of priority streamers
	checkStreamers(client, viper.GetStringSlice("priority_streamers"), bannedGames)

	// Check list of priority streamers
	checkStreamers(client, viper.GetStringSlice("backup_streamers"), bannedGames)
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

func configureClient() (*helix.Client, error) {
	err := readConfigFile()
	if err != nil {
		log.Fatal("Could not read config file:", err.Error())
	}

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
		log.Fatal("Error getting channel info:", err.Error())
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

	// TODO Is channel streaming a banned game?
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

	return helix.Channel{}, errors.New("No channels found that matched: " + channel)
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
