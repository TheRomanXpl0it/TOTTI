package config

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"sub/db"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogLevel string `yaml:"log_level"`

	Team          int    `yaml:"team"`
	NumberOfTeams int    `yaml:"number_of_teams"`
	TeamToken     string `yaml:"team_token"`

	TeamFormat string   `yaml:"team_format"`
	TeamIp     string   `yaml:"team_ip"`
	NopTeam    string   `yaml:"nop_team"`
	Teams      []string `yaml:"teams"`

	RoundDuration int    `yaml:"round_duration"`
	FlagAlive     int    `yaml:"flag_alive"`
	FlagFormat    string `yaml:"flag_format"`
	FlagRegex     *regexp.Regexp

	FlagIdUrl string `yaml:"flagid_url"`

	SubProtocol       string `yaml:"sub_protocol"`
	SubLimit          int    `yaml:"sub_limit"`
	SubInterval       int    `yaml:"sub_interval"`
	SubMaxPayloadSize int    `yaml:"sub_max_payload_size"`
	SubUrl            string `yaml:"sub_url"`
	StartRound        string `yaml:"start_round"`
	FirstRound        int64

	Database string `yaml:"database"`
	DB       *db.DB
}

func LoadConfig(path string) (Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	dec := yaml.NewDecoder(file)
	if err = dec.Decode(&config); err != nil {
		return Config{}, err
	}

	config.TeamIp = fmt.Sprintf(config.TeamFormat, config.Team)
	config.NopTeam = fmt.Sprintf(config.TeamFormat, 0)
	config.FlagAlive *= config.RoundDuration
	config.FlagRegex = regexp.MustCompile(config.FlagFormat)

	for i := range config.NumberOfTeams + 1 {
		if i == 0 || i == config.Team {
			continue
		}
		config.Teams = append(config.Teams, fmt.Sprintf(config.TeamFormat, i))
	}

	startRound, err := time.Parse(time.DateTime, config.StartRound)
	if err != nil {
		return Config{}, err
	}
	config.FirstRound = startRound.Unix()

	return config, nil
}
