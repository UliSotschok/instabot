package bot

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"
	"unsafe"

	"github.com/ahmdrz/goinsta/v2"
	wr "github.com/mroth/weightedrand"
)

// MyInstabot is a wrapper around everything
type MyInstabot struct {
	instaApi       *goinsta.Instagram
	counterActions int
	config         config
	actionLogs     []actionLog
	followingCache []string
}

// TODO: replace with DB
type actionLog struct {
	hashtag string
	action  ActionType
	time    time.Time
}

func RunBots() {
	var bot MyInstabot

	// Gets the command line options
	parseOptions()
	// Gets the config
	conf, err := readConfig("config/config.json")
	if err != nil {
		log.Fatalf("Could not read 'config.json'. error='%s'\n", err.Error())
		return
	}
	bot.config = conf

	// Tries to login
	err = bot.login()
	if err != nil {
		log.Fatalf("Login failed. error='%s'\n", err.Error())
		return
	}

	if unfollow {
		bot.syncFollowers()
	} else if run {
		// Loop through tags ; follows, likes, and comments, according to the config file
		bot.mainLoop()
	}
}

func (bot *MyInstabot) syncFollowers() {
	following := bot.instaApi.Account.Following()
	followers := bot.instaApi.Account.Followers()

	var followerUsers []goinsta.User
	var followingUsers []goinsta.User

	for following.Next() {
		for _, user := range following.Users {
			followingUsers = append(followingUsers, user)
			bot.followingCache = append(bot.followingCache, user.Username)
		}
	}
	for followers.Next() {
		for _, user := range followers.Users {
			followerUsers = append(followerUsers, user)
		}
	}

	var users []goinsta.User
	for _, user := range followingUsers {
		// Skip whitelisted users.
		if containsString(bot.config.UserWhitelist, user.Username) {
			continue
		}

		if !containsUser(followerUsers, user) {
			users = append(users, user)
		}
	}
	if len(users) == 0 {
		return
	}
	fmt.Printf("\n%d users are not following you back!\n", len(users))

	answerUnfollowAll := getInput("Unfollow everyone[A] or decide for each user[R]?")

	for _, user := range users {
		if answerUnfollowAll == "R" {
			answerUserUnfollow := getInputf("Unfollow %s ? [yN]", user.Username)
			if answerUserUnfollow == "y" {
				bot.config.UserBlacklist = append(bot.config.UserBlacklist, user.Username)
				if !dev {
					err := user.Unfollow()
					if err != nil {
						log.Printf("Error unfollowing user. error='%s'\n", err.Error())
					}
				}
				doPauseAfterAction(bot.config.Scheduling)
			} else {
				bot.config.UserWhitelist = append(bot.config.UserWhitelist, user.Username)
			}
		}
	}
}

func (bot *MyInstabot) mainLoop() {
	bot.counterActions = 0
	for true {
		conf := bot.pickRandomTag()
		img, user, err := bot.findImage(conf)
		if err != nil {
			log.Printf("Error finding image: error='%s'", err.Error())
			continue
		}
		bot.executeActions(img, user, conf)
		if bot.counterActions > bot.config.Scheduling.BatchSize {
			doPauseAfterBatch(bot.config.Scheduling)
			bot.counterActions = 0
		}
	}
}

func (bot *MyInstabot) pickRandomTag() hashtagConfig {
	// TODO: cache chooser?
	choices := make([]wr.Choice, len(bot.config.HashtagConfigs))
	for i, elem := range bot.config.HashtagConfigs {
		choices[i] = wr.NewChoice(elem, elem.WeightToChoose)
	}
	chooser, _ := wr.NewChooser(choices...)
	return chooser.Pick().(hashtagConfig)
}

func (bot *MyInstabot) findImage(conf hashtagConfig) (goinsta.Item, *goinsta.User, error) {
	log.Printf("Fetching the list of images for #%s\n", conf.Hashtag)
	var img goinsta.Item
	var user *goinsta.User

	var images *goinsta.FeedTag
	err := retry(4, 20*time.Second, func() (err error) {
		images, err = bot.instaApi.Feed.Tags(conf.Hashtag)
		return
	})
	if err != nil {
		return img, user, err
	}

	for _, image := range images.Images {
		// skip own images
		if image.User.Username == bot.config.Authentication.Username {
			continue
		}
		img = image

		err := retry(10, 20*time.Second, func() (err error) {
			user, err = bot.instaApi.Profiles.ByName(image.User.Username)
			return
		})
		if err != nil {
			return img, user, err
		}

		// skip users based on follower limits
		if user.FollowerCount < conf.MinFollowers || user.FollowerCount > conf.MaxFollowers {
			continue
		}

		// skip users we already follow
		if bot.checkIfFollowing(user) {
			continue
		}

		return img, user, nil // we found an image and did our action(s)
	}
	log.Printf("Warning: No fitting image found for hashtag='%s'\n", conf.Hashtag)
	return img, user, errors.New("no image found")
}

func (bot MyInstabot) checkIfFollowing(user *goinsta.User) bool {
	for _, followingUser := range bot.followingCache {
		if followingUser == user.Username {
			return true
		}
	}
	return false
}

func (bot *MyInstabot) executeActions(image goinsta.Item, user *goinsta.User, conf hashtagConfig) {
	bot.executeAction(like, image, user, conf)
	bot.executeAction(comment, image, user, conf)
	bot.executeAction(follow, image, user, conf)
}

func (bot *MyInstabot) executeAction(actionType ActionType, image goinsta.Item, user *goinsta.User, conf hashtagConfig) {
	switch actionType {
	case like:
		bot.likeImage(image)
	case comment:
		bot.commentImage(conf.Comments, image)
	case follow:
		bot.followUser(user)
	}
	bot.counterActions++
	bot.actionLogs = append(bot.actionLogs, actionLog{
		hashtag: conf.Hashtag,
		action:  actionType,
		time:    time.Now(),
	})
	doPauseAfterAction(bot.config.Scheduling)
}

// Likes an image, if not liked already
func (bot *MyInstabot) likeImage(image goinsta.Item) {
	log.Println("Liking the picture")
	if !image.HasLiked {
		if !dev {
			err := image.Like()
			if err != nil {
				log.Printf("Error liking image: error='%v'\n", err)
			} else {
				log.Println("Liked")
			}
		}
	} else {
		log.Println("Image already liked")
	}
}

// Comments an image
func (bot *MyInstabot) commentImage(comments []string, image goinsta.Item) {
	rand.Seed(time.Now().Unix())
	text := comments[rand.Intn(len(comments))]
	if !dev {
		var err error = nil
		comments := image.Comments
		if comments == nil {
			// monkey patching
			// we need to do that because https://github.com/ahmdrz/goinsta/pull/299 is not in goinsta/v2
			// I know, it's ugly
			newComments := goinsta.Comments{}
			rs := reflect.ValueOf(&newComments).Elem()
			rf := rs.FieldByName("item")
			rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
			item := reflect.New(reflect.TypeOf(image))
			item.Elem().Set(reflect.ValueOf(image))
			rf.Set(item)
			err = newComments.Add(text)
			// end hack!
		} else {
			err = comments.Add(text)
		}
		if err != nil {
			log.Printf("Error commenting image: error='%v'\n", err)
		} else {
			log.Println("Commented " + text)
		}
	}

}

// Follows a user, if not following already
func (bot *MyInstabot) followUser(user *goinsta.User) {
	log.Printf("Following %s\n", user.Username)
	err := user.FriendShip()
	check(err)
	// If not following already
	if !user.Friendship.Following {
		if !dev {
			err = user.Follow()
			if err != nil {
				log.Printf("Error following user: error='%v'\n", err)
			} else {
				log.Println("Followed " + user.Username)
				bot.followingCache = append(bot.followingCache, user.Username)
			}
		}
	} else {
		log.Println("Already following " + user.Username)
	}
}
