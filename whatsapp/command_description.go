package whatsapp

type commandDescription struct {
	shortDescription string
	longDescription  string
}

var (
	descriptions = map[string]commandDescription{
		"#!loki": commandDescription{
			shortDescription: "adds user to blacklist",
			longDescription: `
Adds user to the blacklist. This user will always get a message by the bot`,
		},

		"#!thor": commandDescription{
			shortDescription: "adds user to whitelist",
			longDescription: `
Adds user to the whitelist. This user will never get a message by the bot`,
		},

		"#!odin": commandDescription{
			shortDescription: "toggles the do not disturb between *RUNNING* and *IDLE*",
			longDescription: `
Toggles the do not disturb between *RUNNING* and *IDLE*`,
		},

		"#!hella": commandDescription{
			shortDescription: "sets contact to a normal contact if either in whitelist or blacklist",
			longDescription: `
Sets contact to a normal contact if either in whitelist or blacklist`,
		},

		"#!freyja": commandDescription{
			shortDescription: "shows help command for users",
			longDescription: `
use cases *#!freyja [COMMAND]* to get help on a command. if no command is passed, it will list 
a short desctipion of all the available commands.
`,
		},
	}
)
