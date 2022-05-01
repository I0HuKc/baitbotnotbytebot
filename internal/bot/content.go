package baitbot

// Actions
var (
	ActionSendJoke = CreateAction("SEND_JOKE")
)

// Commands
var (
	CommandStart              = CreateCommand("start")
	CommandAddDecription      = CreateCommand("ad", "-v")
	CommandGetDecription      = CreateCommand("gd")
	CommandGetDescriptionList = CreateCommand("gdl")
	CommandSendJoke           = CreateCommand("sjoke")

	CommandAntre                 = CreateCommand("antre")
	CommandJoke                  = CreateCommand("joke")
	CommandPing                  = CreateCommand("ping")
	CommandStartChangeDecription = CreateCommand("scd")
)
