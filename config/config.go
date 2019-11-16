package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var (
	//Token is the OAUTH token used by the bot
	Token string

	//BotPrefix I honestly am not sure yet
	BotPrefix string

	//WeatherstackAPIKey I honestly am not sure yet
	WeatherstackAPIKey string

	//Private variables
	config *configStruct

	//WeatherStackURL is the api endpoint for weather information
	WeatherStackURL = "http://api.weatherstack.com/current"

	//WeatherStackDefaultLocation is the default location for weather information
	WeatherStackDefaultLocation = "Los Angeles"

	//WoofURL is the api endpoint for dog pictures
	WoofURL = "https://random.dog/woof.json"

	//FoaasURL is the api endpoint for Foaas
	FoaasURL = "https://foaas.com"

	//GiphyURL is the api endpoint for giphy
	GiphyURL = "http://api.giphy.com/v1/gifs/search"

	//GiphyAPIKey is the api key for giphy
	GiphyAPIKey string
)

// used for local config
type configStruct struct {
	Token              string `json:"Token"`
	BotPrefix          string `json:"BotPrefix"`
	WeatherstackAPIKey string `json:"WeatherstackAPIKey"`
	GiphyAPIKey        string `json:"GiphyAPIKey"`
}

// EnvMissingError defines
type EnvMissingError struct {
	msg string
}

func (error *EnvMissingError) Error() string {
	return fmt.Sprintf("environment variable: '%s' missing", error.msg)
}

// EnvMissing defines
func EnvMissing(s string) error {
	return &EnvMissingError{s}
}

//GetenvStr stub
func GetenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, EnvMissing(key)
	}
	return v, nil
}

//GetenvBool stub
func GetenvBool(key string) (bool, error) {
	s, err := GetenvStr(key)
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}

//ReadConfig is used to read the config.json
func ReadConfig() error {
	local, err := GetenvBool("LOCAL")
	if err != nil {
		return err
	}

	if local {
		file, err := ioutil.ReadFile("./config.json")

		if err != nil {
			return err
		}

		err = json.Unmarshal(file, &config)

		if err != nil {
			return err
		}

		Token = config.Token
		BotPrefix = config.BotPrefix
		WeatherstackAPIKey = config.WeatherstackAPIKey
		GiphyAPIKey = config.GiphyAPIKey
		return nil
	}

	token, err := GetenvStr("TOKEN")
	botPrefix, err := GetenvStr("BOT_PREFIX")
	weatherstackAPIKey, err := GetenvStr("WEATHERSTACK_API_KEY")
	giphyAPIKey, err := GetenvStr("GIPHY_API_KEY")

	if err != nil {
		return err
	}

	Token = token
	BotPrefix = botPrefix
	WeatherstackAPIKey = weatherstackAPIKey
	GiphyAPIKey = giphyAPIKey

	return nil
}
