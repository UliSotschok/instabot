package bot

import (
	"encoding/json"
	"io/ioutil"
)

var (
	// Whether we are in development mode or not
	dev bool

	// Whether we want an email to be sent when the script ends / crashes
	nomail bool

	// Whether we want to launch the unfollow mode
	unfollow bool

	// Acut
	run bool

	// Whether we want to have logging
	logs bool

	// Used to skip following, liking and commenting same user in this session
	noduplicate bool
)

type ActionType string

const (
	like    ActionType = "like"
	comment ActionType = "comment"
	follow  ActionType = "follow"
)

type config struct {
	HashtagConfigs []hashtagConfig `json:"hashtag_configs"`
	UserBlacklist  []string        `json:"user_blacklist"`
	UserWhitelist  []string        `json:"user_whitelist"`
	Scheduling     scheduling      `json:"scheduling"`
	Authentication authentication  `json:"authentication"`
}

type scheduling struct {
	PauseAfterActionInS int `json:"pause_after_action_in_s"`
	BatchSize           int `json:"batch_size"`
	PauseAfterBatchInM  int `json:"pause_after_batch_in_m"`
}

type authentication struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type hashtagConfig struct {
	Hashtag        string   `json:"hashtag"`
	Comments       []string `json:"comments"`
	MinFollowers   int      `json:"min_followers"`
	MaxFollowers   int      `json:"max_followers"`
	WeightToChoose uint     `json:"weight_to_choose"`
}

func readConfig(path string) (config, error) {
	var conf config
	text, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}
	err = json.Unmarshal(text, &conf)
	if err != nil {
		return conf, err
	}
	// TODO: validate
	return conf, nil
}
