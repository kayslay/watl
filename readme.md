# WATL 

WATL is short for **W**hats**A**pp to **T**e**L**egram. The main aim of WATL is to redirect messages from whatsapp due to friends that are not willing to move.

I'll like to give a shout-out to the [contributors](https://github.com/Rhymen/go-whatsapp/graphs/contributors) of [https://github.com/Rhymen/go-whatsapp/graphs/contributors](https://github.com/Rhymen/go-whatsapp/graphs/contributors). THey did all the work, I am just building on top of what they did.

For now I can't do that because I have some lazy excuses. So for now, what I created is something to reply my messages in do not disturb mode.

The users phone must be active to use it. Just like the way your phone has to be active for whatsapp web. Don't blame me, as I said earlier, I built on what was given to me. Now, if you want to use it without your phone being online, **BUILD IT YOURSELF** 😁.

## How to use

Now if u want to use the one that is hosted up in the cloud above us, visit this link [watl](http://ec2-34-226-217-107.compute-1.amazonaws.com:8000/public/).

If you are using project locally, make sure Go is installed on your laptop. If Go is not installed, it will never work, **NEVER**.

Now, if you have Go installed, run `make watl` to start the web server. Now that the server has started, visit [http://localhost:8000/public/](http://localhost:8000/public/).

Now that the page has loaded, you must be wandering what the fuck is this? Well it's the fucking lack of good frontend skills. So if u have a cool design you can add it.

The bot has 2 states that are of interest to you which are `IDLE` and `RUNNING`. Whe the state is `RUNNING` the bot is in a *DO NOT DISTURB MODE*, and every message sent will be marked as read and will be replied to by the bot

**NOTE**: if you are a fan of whatsapp web, this is not for you. For some reason whatsapp limits the number of active devices to 1, so I can't make it more. One has to be logged out for one to be logged in.

**NOTE**: The barcode expires after 30 sec or so (can't say). if the scan is not successful, reload or click on the page to get a new barcode.

## Commands

For some reason which I don't understand I'm making use of names from Thor and Norse Mythology. The commands don't have any relation to the attributes of this names.

The command are passed through the chat input or status input.

**#!freyja**: act as the help command. Shows info about the commands.

For list of other commands type `#!freyja`. I will list the remaining commands when I'm ready to.
Since `#!freyja` shows al the commands it will be redudant adding them here 😏.

## Motivation

Apart from doing it for Knowledge sake and not having a side project to work on, I did it for **MONEY, POWER AND GLORY**.

Now if you face any bug 🐞 or problem (**WHICH YOU WILL**) contact me on Twitter: [@kayslaycode](https://twitter.com/Kayslaycode) if you can't meet me physically.

