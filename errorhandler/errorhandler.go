package errorhandler

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func LogErrorAndExit(err error) {
	log.Error(err)
	os.Exit(1)
}
