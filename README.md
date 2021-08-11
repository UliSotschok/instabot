[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0) [![Made with: Golang](https://img.shields.io/badge/Made%20with-Golang-brightgreen.svg)](https://golang.org/)

[![Run on Repl.it](https://repl.it/badge/github/tducasse/go-instabot)](https://repl.it/github/tducasse/go-instabot)

### Not actively maintained, feel free to fork 👍

# What is go-instabot?

The easiest way to boost your Instagram account and get likes and followers.

Go-instabot automates **following** users, **liking** pictures, **commenting**, and **unfollowing** people that don't follow you back on Instagram.

It uses the unofficial but excellent Go Instagram API, [goinsta](https://github.com/ahmdrz/goinsta) (v2).


### Concept
The idea behind the script is that when you like, follow, or comment something, it will draw the user's attention back to your own account. There's a hidden convention in Instagram, that will make people follow you back, as a way of saying "thank you" I guess.

Moreover, you may have noticed that when you follow someone, Instagram tells you about 'similar people to follow'. The more active you are on Instagram, the more likely you are to be in this section.


# How does it work?
## Algorithm
- There is a config file, where you can enter hashtags, and for each hashtag, how many likes, comments, and follows you want.
- The script will then fetch pictures from the 'explore by tags' page on Instagram, using the API.
- It will decide (based on your settings) if it has to follow, like or comment.
- At the end of the routine, an email can be sent, with a report on everything that's been done.

Additionally, there is a retry mechanism in the eventuality that Instagram is too slow to answer (it can happen sometimes), and the script will wait for some time before trying again.

## Bans
The script is coded so that your Instagram account will not get banned ; it waits between every call to simulate human behavior.

The first time it logs you in, it will store the session object in a file, encrypted with AES (thanks to [goinsta](https://github.com/ahmdrz/goinsta)). Every next launch will try to recover the session instead of logging in again. This is a trick to avoid getting banned for connecting on too many devices (as suggested by [anastalaz](https://github.com/tducasse/go-instabot/issues/1)).


# How to use
## Installation

1. [Install Go](https://golang.org/doc/install) on your system.

2. Download and install go-instabot, by executing this command in your terminal / cmd :

   `go get github.com/tducasse/go-instabot`

## Configuration
### Config.json
Go to the project folder :

`cd [YOUR_GO_PATH]/src/github.com/tducasse/go-instabot`

There, in the 'dist/' folder, you will find a sample 'config.json', that you have to copy to the 'config/' folder :

```json
{
    "hashtag_configs": [
        {
            "hashtag": "dog",
            "actions": {
                "like": {
                    "min_followers": 5,
                    "max_followers": 200,
                    "weight_to_choose": 20
                },
                "comment": {
                    "min_followers": 5,
                    "max_followers": 200,
                    "weight_to_choose": 1
                },
                "follow": {
                    "min_followers": 5,
                    "max_followers": 200,
                    "weight_to_choose": 5
                }
            },
            "comments": [
                "😍😍😍",
                "😻😻😻",
                "wow nice pic",
                "Lovely 😍",
                "Wonderful 😍😍👏👏"
            ],
            "weight_to_choose": 1
        }
    ],
        "user_blacklist": [],
        "user_whitelist": [],
        "scheduling": {
        "pause_after_action_in_s": 20,
            "batch_size": 20,
            "pause_after_batch_in_m": 360
    },
    "authentication": {
        "username": "YahooGMBH",
        "password": "$M\u0026l9K06IyVi71^\u0026iEh$m"
    }
}
```


Ich habe mir nochmal eine neue Config und Algorithmus überlegt. Kannst ja mal sagen, was noch geändert werden sollte.
Grundfaktoren:
- Der Bot soll endlos laufen. 
- Die Aktionen sollen mehr zufällig sein. 
- Die länge der Pausen sollen nicht einheitlich sein. 
- Es gibt auch sehr lange Pausen nach einer gewissen Zeit.

Algorithmus:

- Wiederhole für immer:
  - Wähle einen hashtag 
    - Zwischen allen hashtags wird gewichtet gezogen. 
    - Das heißt, ein Hashtag mit Gewicht 2 wird doppelt so oft gezogen wie einer mit Gewicht 1.
    - `weight_to_choose` im config file
  - Wähle eine Aktion
    - Zwischen allen Aktionen wird gewichtet gezogen, wie beim hashtag
    - `weight_to_choose` im config file
  - Finde ein Bild zum Hashtag
    - Schließe eigene Bilder aus.
    - Nehme nur Bilder von Usern mit: `min_followers` ≤ Follower ≤ `max_followers`
    - Ignoriere Bilder von Usern, die in `user_blacklist` stehen.
  - Führe die Aktion auf das gefundene Bild aus
    - Bei comment und follow, wird das Bild davor noch geliked.
    - Nach jeder Aktion wird X Sekunden gewartet: randomZahl(zwischen 0 bis `pause_after_action_in_s`) + `pause_after_action_in_s`
  - Wurden `batch_size` Aktionen durchgeführt, wird gewartet
    - Minuten: randomZahl(zwischen 0 bis `pause_after_batch_in_m`) + `pause_after_batch_in_m`

Die Neue config würde so aussehen. Wenn was unklar ist, kannst nochmal fragen. Weißt du eigentlich wie das json Format aufgebaut ist?

```json
{
    "hashtag_configs": [
        {
            "hashtag": "dog",
            "actions": {
                "like": {
                    "min_followers": 5,
                    "max_followers": 200,
                    "weight_to_choose": 20
                },
                "comment": {
                    "min_followers": 5,
                    "max_followers": 200,
                    "weight_to_choose": 1
                },
                "follow": {
                    "min_followers": 5,
                    "max_followers": 200,
                    "weight_to_choose": 5
                }
            },
            "comments": [
                "wow nice pic",
                "Lovely 😍"
            ],
            "weight_to_choose": 1
        }
    ],
    "user_blacklist": [],
    "user_whitelist": [],
    "scheduling": {
        "pause_after_action_in_s": 20,
        "batch_size": 20,
        "pause_after_batch_in_m": 360
    },
    "authentication": {
        "username": "HIDDEN",
        "password": "HIDDEN"
    }
}
```


# TODO

### Note on the emails

I use Gmail to send and receive the emails. If you want to use Gmail, there's something important to do first (from the [Google accounts](https://support.google.com/accounts/answer/6010255) website) :
```
Change your settings to allow less secure apps to access your account.
We don't recommend this option because it might make it easier for someone to break into your account.
If you want to allow access anyway, follow these steps:
    - Go to the "Less secure apps" section in My Account.
    - Next to "Access for less secure apps," select Turn on.
```
(If you can't find where it is exactly, I think [this link](https://myaccount.google.com/security) should work)

As this procedure might not be safe, I recommend not doing it on your main Gmail account, and maybe create another account on the side. Or try to find a less secure webmail provider!

## How to run
This is it!
Since you used the `go get` command, you now have the `go-instabot` executable available from anywhere* in your system. Just launch it in a terminal :

`go-instabot -run`

**\*** : *You will need to have a folder named 'config' (with a 'config.json' file) in the directory where you launch it.*

### Options
**-run** : This is the main option. Use it to actually launch the script.

**-h** : Use this option to display the list of options.

**-dev** : Use this option to use the script in development mode : nothing will be done for real. You will need to put a config file in a 'local' folder.

**-logs** : Use this option to enable the logfile. The script will continue writing everything on the screen, but it will also write it in a .log file.

**-nomail** : Use this option to disable the email notifications.

**-sync** : Use this option to unfollow users that don't follow you back. Don't worry, the script will ask before actually doing it, so you can use it just to check the number!

**-noduplicate** : Use this to skip following, liking and commenting same user in this session!

### Tips
- If you want to launch a long session, and you're afraid of closing the terminal, I recommend using the command __screen__.
- If you have a Raspberry Pi, a web server, or anything similar, you can run the script on it (again, use screen).
- To maximize your chances of getting new followers, don't spam! If you follow too many people, you will become irrelevant.

  Also, try to use hashtags related to your own account : if you are a portrait photographer and you suddenly start following a thousand #cats related accounts, I doubt it will bring you back a thousand new followers...
  
Good luck getting new followers!

### ⚠️ Reporting issues/PRs/license
This is _very_ loosely maintained, as in, I'll _probably_ try and fix things if everything is broken, but I'm no longer working on it. Feel free to fork it though!

