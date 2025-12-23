package main

import (
	"github.com/odysseia-greek/agora/plato/logging"
)

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=DAREIOS&x=none&v=4&h=4&w=80&we=false
	logging.System(`
 ___     ____  ____     ___  ____  ___   _____
|   \   /    ||    \   /  _]|    |/   \ / ___/
|    \ |  o  ||  D  ) /  [_  |  ||     (   \_ 
|  D  ||     ||    / |    _] |  ||  O  |\__  |
|     ||  _  ||    \ |   [_  |  ||     |/  \ |
|     ||  |  ||  .  \|     | |  ||     |\    |
|_____||__|__||__|\_||_____||____|\___/  \___|
`)
	logging.System("\"οἱ βασιλέως λόγοι πρὸς τοὺς Ἕλληνας οὐκ ἀληθεῖς εἰσίν.\"")
	logging.System("The king's words to the Greeks are not to be trusted.")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")
}
