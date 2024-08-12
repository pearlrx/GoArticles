package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Log.SetLevel(logrus.InfoLevel)

	Log.SetOutput(os.Stdout)
}
