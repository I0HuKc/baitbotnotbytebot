package baitbot

// Commands
var (
	CommandStart              = CreateCommand("start")
	CommandHelp               = CreateCommand("help")
	CommandAddDecription      = CreateCommand("ad", "-v")
	CommandGetDecription      = CreateCommand("gd")
	CommandGetDescriptionList = CreateCommand("gdl")

	CommandAntre                 = CreateCommand("antre")
	CommandStopAntre             = CreateCommand("sa")
	CommandJoke                  = CreateCommand("joke")
	CommandPing                  = CreateCommand("ping")
	CommandStartChangeDecription = CreateCommand("scd")
)

var helpInfo string = `
/start — перезапуск бота (локальный).

*Групповые команды*
/scd — старт процесса смены описания
/bll — начать буллить человека
/cd — принудительное изменение статуса в группе

*Приватные команды*
/ad -v <value> — добавить статус (for Authors)
/gd -id <recordid> — получить статус (for Admins).
/help — получить эту инструкцию.
`
