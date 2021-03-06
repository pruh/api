package config

import (
	"encoding/json"
	"errors"
	"github.com/golang/glog"
	"net"
	"os"
	"strconv"
)

// Configuration contrains configuration parameters.
type Configuration struct {
	Port             *string
	TelegramBoToken  *string
	DefaultChatID    *int
	APIV1Credentials *map[string]string
	LocalNets        []*net.IPNet

	MongoUsername *string
	MongoPassword *string
}

// NewFromEnv creates new configuration from environment variables.
func NewFromEnv() (*Configuration, error) {
	port := ptr(os.Getenv("PORT"))
	botToken := ptr(os.Getenv("TELEGRAM_BOT_TOKEN"))
	chatID := ptrOrNil(os.LookupEnv("TELEGRAM_DEFAULT_CHAT_ID"))
	apiCreds := ptrOrNil(os.LookupEnv("API_V1_CREDS"))

	mongoUsername := ptrOrNil(os.LookupEnv("MONGO_INITDB_ROOT_USERNAME"))
	mongoPassword := ptrOrNil(os.LookupEnv("MONGO_INITDB_ROOT_PASSWORD"))

	return NewFromParams(port, botToken, chatID, apiCreds, mongoUsername, mongoPassword)
}

// NewFromParams creates new configuration from arguments.
func NewFromParams(port *string, boToken *string, defaultChatID *string,
	apiV1Credentials *string, mongoUsername *string, mongoPassword *string) (*Configuration, error) {
	var conf Configuration
	if port == nil || *port == "" {
		return nil, errors.New("port should not be empty")
	}
	conf.Port = port

	if boToken == nil || *boToken == "" {
		return nil, errors.New("telegram bot token should not be empty")
	}
	conf.TelegramBoToken = boToken

	if defaultChatID != nil && *defaultChatID != "" {
		if chatID, err := strconv.Atoi(*defaultChatID); err == nil {
			conf.DefaultChatID = &chatID
		} else {
			return nil, err
		}
	}

	if apiV1Credentials != nil && *apiV1Credentials != "" {
		conf.APIV1Credentials = &map[string]string{}
		err := json.Unmarshal([]byte(*apiV1Credentials), conf.APIV1Credentials)
		if err != nil {
			return nil, err
		}
	}

	conf.LocalNets = getLocalIPNets()
	if mongoUsername != nil && mongoPassword != nil {
		conf.MongoUsername = mongoUsername
		conf.MongoPassword = mongoPassword
	}

	return &conf, nil
}

func ptr(str string) *string {
	return &str
}

func ptrOrNil(str string, valueSet bool) *string {
	if valueSet {
		return &str
	}

	return nil
}

func getLocalIPNets() []*net.IPNet {
	var localIPNets []*net.IPNet
	cidrs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, cidr := range cidrs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			glog.Infof("Cannot parse CIDR %s", cidr)
			continue
		}
		localIPNets = append(localIPNets, ipnet)
	}
	return localIPNets
}
