package infocmdb

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"os"
)

type logOutputSplitter struct{}

func (splitter *logOutputSplitter) Write(msg []byte) (n int, err error) {
	if bytes.HasPrefix(msg, []byte("[TRACE]")) ||
		bytes.HasPrefix(msg, []byte("[DEBUG]")) ||
		bytes.HasPrefix(msg, []byte("[INFO]")) {
		return os.Stdout.Write(msg)
	}
	return os.Stderr.Write(msg)
}

func init() {
	// Log to stdout and stderr depending on log level:
	// Any message on stderr is interpreted as workflow failure
	log.SetOutput(&logOutputSplitter{})
	// Time is omitted in the log message, because it is already shown in the workflow log in a separate column
	log.SetFormatter(&easy.Formatter{
		LogFormat: "[%lvl%] %msg%\n",
	})
	log.SetLevel(log.InfoLevel)
	if os.Getenv("WORKFLOW_DEBUGGING") == "true" {
		log.SetLevel(log.DebugLevel)
	}
}
