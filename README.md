[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0) [![Made with: Golang](https://img.shields.io/badge/Made%20with-Golang-brightgreen.svg)](https://golang.org/)


# What is **instabot**?

**instabot** automates **following** users, **liking** pictures, **commenting** to get more followers.

It uses the unofficial but excellent Go Instagram API, [goinsta](https://github.com/ahmdrz/goinsta) (v2).


### Concept
The idea behind the script is that when you like, follow, or comment something, it will draw the user's attention back to your own account. There's a hidden convention in Instagram, that will make people follow you back, as a way of saying "thank you" I guess.

Moreover, you may have noticed that when you follow someone, Instagram tells you about 'similar people to follow'. The more active you are on Instagram, the more likely you are to be in this section.


# How does it work?

## Algorithm

- run forever
  - Choose a hashtag from the config
    - Hashtag is chosen randomly with weighted drawing.
    - A hashtag with weight 2 is chosen twice as often a one with weight 1.
    - `weight_to_choose` in config file
  - Find a picture
    - Ignore own images
    - Choose only pictures from users with: `min_followers` ≤ followers ≤ `max_followers`
    - ignore users from `user_blacklist`.
  - Like, comment and follow
    - after each action wait for seconds: randomNumber(between 0 and `pause_after_action_in_s`) + `pause_after_action_in_s`
  - After `batch_size` actions, wait for a longer time
    - minutes: randomNumber(between 0 and `pause_after_batch_in_m`) + `pause_after_batch_in_m`

## Bans

To reduce the risk of getting banned, pauses are inserted. You can configure the length. The longer the pause, the lower the risk of getting banned.

### Disclaimer

Although I did my best to add options that reduce the risk of getting banned, I take no liability if anything happens to your Account.
Use this code at your own risk.

The first time it logs you in, it will store the session object in a file, encrypted with AES (thanks to [goinsta](https://github.com/ahmdrz/goinsta)). Every next launch will try to recover the session instead of logging in again. This is a trick to avoid getting banned for connecting on too many devices (as suggested by [anastalaz](https://github.com/tducasse/go-instabot/issues/1)).


# How to use

## for non programmers: As binary

1. Download the latest binary from the [releases page](https://github.com/UliSotschok/instabot/releases)
2. In the folder where you downloaded the binary create a folder `config` and create a file `config.json`
3. [Configure](#configuration) it. You can use [this example config](config/config_example.json) as starting point.
4. [Run](#run) the binary and enjoy getting followers😍/banned😈

## for programmers: from source code

1. `git clone https://github.com/UliSotschok/instabot/`
2. `cd instabot`
3. `go get`
4. 

## Configuration

- Example config file: [config_example.json](config/config_example.json)
- Explanation of all attributes in [config.go](bot/config.go)

## Run

### Options

**-h** : Use this option to display the list of options.

**-dev** : Use this option to use the script in development mode : nothing will be done for real.

### Tips
- If you want to launch a long session, and you're afraid of closing the terminal, I recommend using the command __screen__.
- If you have a Raspberry Pi, a web server, or anything similar, you can run the script on it (again, use screen).
- To maximize your chances of getting new followers, don't spam! If you follow too many people, you will become irrelevant.

  Also, try to use hashtags related to your own account : if you are a portrait photographer and you suddenly start following a thousand #cats related accounts, I doubt it will bring you back a thousand new followers...
  
Good luck getting new followers!
