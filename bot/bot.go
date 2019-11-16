package bot

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/bwmarrin/discordgo"
	"github.com/haydenmcfarland/discord_chan/config"
	log "github.com/sirupsen/logrus"
)

//BotID contains the ID of the bot
var BotID string
var goBot *discordgo.Session

//Start runs the go bot
func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	u, err := goBot.User("@me")

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	log.Info("bot successfully initialized")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return
	}

	if strings.Contains(m.Content, config.BotPrefix) {
		commandHandler(s, m)
		return
	}
}

func sendMessage(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	_, _ = s.ChannelMessageSend(m.ChannelID, msg)
}

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	r := regexp.MustCompile(`!(\w+)\s*(.*)?`)
	g := r.FindStringSubmatch(m.Content)

	if len(g) < 1 {
		log.Warn("Empty message was propagated to the command handler")
		return
	}

	response := ""
	var err error

	switch g[1] {
	case "now":
		response = fmt.Sprintf("The current time is: %s", time.Now().Format("3:04:05 PM"))
	case "ping":
		response, err = foaasHandler(m.Author.Username, "Everyone")
	case "roll":
		response = fmt.Sprintf("%s rolled a %s", m.Author, strconv.Itoa(rand.Intn(100)))
	case "weather":
		if len(g) > 2 {
			response, err = weatherStackHandler(g[2])
		}
	case "echo":
		if len(g) > 2 {
			response = g[2]
		}
	case "woof":
		response, err = woofHandler()
	case "gif":
		if len(g) > 2 {
			response, err = giphyHandler(g[2])
		}
	}

	if err != nil {
		log.Error(err.Error())
		sendMessage(s, m, "Something went wrong with your request. Sucks to suck.")
		return
	}

	if response != "" {
		sendMessage(s, m, response)
	}
}

//getRequest take 1 URL argument and 2 position argumetns params and headers
//defaults to application/json if header position argument is not provided
func getRequest(URL string, args ...map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	argLen := len(args)

	if argLen > 0 {
		q := req.URL.Query()
		for k, v := range args[0] {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	if argLen > 1 {
		for k, v := range args[1] {
			req.Header.Add(k, v)
		}
	} else {
		req.Header.Add("Accept", "application/json")
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func jsonParser(data []byte, args ...string) ([]byte, error) {
	result, _, _, err := jsonparser.Get(data, args...)
	if err != nil {
		log.WithFields(log.Fields{
			"data":  string(data),
			"args":  fmt.Sprintf("%s", args),
			"error": err.Error(),
		}).Error("json parsing error")
		return nil, err
	}
	return result, nil
}

func giphyHandler(s string) (string, error) {
	params := map[string]string{
		"key":   config.GiphyAPIKey,
		"q":     s,
		"limit": "1",
	}

	body, err := getRequest(config.GiphyURL, params)
	if err != nil {
		return "", err
	}

	res, err := jsonParser(body, "data", "[0]", "url")
	if err != nil {
		return "", err
	}

	quotedURL := string(res)
	url := strings.Replace(quotedURL, "\\/", "/", -1)
	return url, nil
}

func woofHandler() (string, error) {
	body, err := getRequest(config.WoofURL)

	if err != nil {
		return "", err
	}

	url, err := jsonParser(body, "url")

	if err != nil {
		return "", err
	}

	return string(url), nil
}

func foaasHandler(user string, from string) (string, error) {
	url := fmt.Sprintf("%s/madison/%s/%s", config.FoaasURL, user, from)
	body, err := getRequest(url)

	if err != nil {
		return "", err
	}

	res, err := jsonParser(body, "message")
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func weatherStackHandler(s string) (string, error) {
	var location string
	if s != "" {
		location = s
	} else {
		location = config.WeatherStackDefaultLocation
	}

	params := map[string]string{
		"access_key": config.WeatherstackAPIKey,
		"query":      location,
	}
	body, err := getRequest(config.WeatherStackURL, params)
	if err != nil {
		return "", err
	}

	rawTemp, err := jsonParser(body, "current", "temperature")
	rawRegion, err := jsonParser(body, "location", "region")
	rawDescription, err := jsonParser(body, "current", "weather_descriptions", "[0]")
	rawName, err := jsonParser(body, "location", "name")
	rawHumidity, err := jsonParser(body, "current", "humidity")

	if err != nil {
		return "", err
	}

	region, err := strconv.Unquote(`"` + string(rawRegion) + `"`)
	temp, err := strconv.ParseFloat(string(rawTemp), 32)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"It is %s in %s, %s with a temperature of %.0f °F and humidity at %s percent.",
		rawDescription,
		rawName,
		region,
		1.8*temp+32,
		rawHumidity,
	), nil
}
