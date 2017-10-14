package setup

import (
	"github.com/orange-lightsaber/pretty-safe-backup/settings"
)

func Build() {
	answerSet := mainMenu(&settings.Setup{})
	settings.NewConfig(answerSet)
}
