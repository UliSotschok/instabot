package bot

import (
	"errors"
	"github.com/UliSotschok/goinsta"
	"log"
)

// login will try to reload a previous session, and will create a new one if it can't
func (bot *MyInstabot) login() error {
	err := bot.reloadSession()
	if err != nil {
		return bot.createAndSaveSession()
	}
	return nil
}

// reloadSession will attempt to recover a previous session
func (bot *MyInstabot) reloadSession() error {

	insta, err := goinsta.Import("config/goinsta-session")
	if err != nil {
		log.Printf("msg='session recovery failed' err='%v'\n", err)
		return errors.New("Couldn't recover the session")
	}

	if insta != nil {
		bot.instaApi = insta
	}

	log.Println("Successfully logged in")
	return nil

}

// Logins and saves the session
func (bot *MyInstabot) createAndSaveSession() error {
	bot.instaApi = goinsta.New(bot.config.Authentication.Username, bot.config.Authentication.Password)
	err := bot.instaApi.Login()
	if err != nil {
		log.Printf("msg='Login failed.' err=%v\n", err)
		return err
	}

	err = bot.instaApi.Export("config/goinsta-session")
	if err != nil {
		return err
	}
	log.Println("Created and saved the session")
	return nil
}
