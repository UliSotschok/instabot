package bot

import (
	"encoding/json"
	"io/ioutil"
)

var (
	// Whether we are in development mode or not
	dev bool
)

type ActionType string

const (
	like    ActionType = "like"
	comment ActionType = "comment"
	follow  ActionType = "follow"
)

type config struct {
	HashtagConfigs []hashtagConfig `json:"hashtag_configs"` // list of hashtags in which images are searched
	UserBlacklist  []string        `json:"user_blacklist"`  // users in this list will be ignored. No likes, comments, follow
	UserWhitelist  []string        `json:"user_whitelist"`  // users in this list won't be automatically unfollowed
	Scheduling     scheduling      `json:"scheduling"`      // pauses to prevent getting banned
	Authentication authentication  `json:"authentication"`  // instagram account credentials
}

type scheduling struct {
	PauseAfterActionInS int `json:"pause_after_action_in_s"` // After each action (like, comment, follow) wait for seconds: randInt(PauseAfterActionInS) + PauseAfterActionInS
	BatchSize           int `json:"batch_size"`              // A batch is a series of actions. BatchSize defines how many actions are one batch. Batches are used to make longer pauses after some time
	PauseAfterBatchInM  int `json:"pause_after_batch_in_m"`  // After each batch wait for minutes: randInt(PauseAfterBatchInM) + PauseAfterBatchInM
}

type authentication struct {
	Username string `json:"username"` // instagram Username
	Password string `json:"password"` // instagram Password
}

type hashtagConfig struct {
	Hashtag        string   `json:"hashtag"`          // Hashtag name to search for images
	Comments       []string `json:"comments"`         // list of Comments from which one is randomly chosen every time something is commented
	MinFollowers   int      `json:"min_followers"`    // Min number of followers an account needs to have, make any action
	MaxFollowers   int      `json:"max_followers"`    // Max number of followers an account needs to have, make any action. Can be used to ignore big account who won't follow back anyway
	WeightToChoose uint     `json:"weight_to_choose"` // When a hashtag is chosen, this is the weight, that defines the probability of getting chosen. A weight of 2 is twice as often drawn as one with 1. Weights are relative. (100, 100) = (1, 1)
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
