package main

import (
	"errors"
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

	for _, channel := range viper.GetStringSlice("priority_streamers") {
		if isChannelRaidCandidate(client, channel, bannedGames) {
			log.Println("Potentially raid: " + channel)
		}
	}
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

func isChannelRaidCandidate(client *helix.Client, channel string, bannedGames *helix.GamesResponse) bool {
	channelStatus, err := getChannelStatus(client, channel)
	if err != nil {
		log.Fatal("Error getting channel info:", err.Error())
	}

	// Is channel live?
	if !channelStatus.IsLive {
		return false
	}

	/*
		// TODO Is channel in preferred language?
		if !channelStatus.IsLive {
			return false
		}
	*/

	// TODO Is channel streaming a banned game?
	for _, game := range bannedGames.Data.Games {
		if channelStatus.GameID == game.ID {
			return false
		}
	}

	return true
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
