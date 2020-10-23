package main

import (
	"fmt"
	"log"

	"github.com/nicklaw5/helix"
	"github.com/spf13/viper"
)

func main() {
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

	/*
		resp, err := client.GetUsers(&helix.UsersParams{
			Logins: []string{"summit1g", "lirik"},
		})
		if err != nil {
			log.Fatal("Error getting users:", err.Error())
		}

		fmt.Printf("Status code: %d\n", resp.StatusCode)
		fmt.Println("Error message: " + resp.ErrorMessage)
		fmt.Printf("Rate limit: %d\n", resp.GetRateLimit())
		fmt.Printf("Rate limit remaining: %d\n", resp.GetRateLimitRemaining())
		fmt.Printf("Rate limit reset: %d\n\n", resp.GetRateLimitReset())

		for _, user := range resp.Data.Users {
			fmt.Printf("ID: %s Name: %s\n", user.ID, user.DisplayName)
		}
	*/

	resp, err := client.SearchChannels(&helix.SearchChannelsParams{
		Channel: "bfroggio",
		First:   1,
	})
	if err != nil {
		log.Fatal("Error searching streams:", err.Error())
	}

	for _, channel := range resp.Data.Channels {
		fmt.Printf("Channel: %s; Live: %t; Language: %s\n", channel.DisplayName, channel.IsLive, channel.Language)
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
