package logger

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func SetupDeviceLogger(level, path string) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		log.Fatalf("unable to open log file %s, error: %v", path, err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)

	log.Infof("Path: %s", path)
	loglevel, err := log.ParseLevel(level)
	if err != nil {
		log.Errorf("unable to parse log level %s, error: %v", level, err)
		return
	}
	log.SetLevel(loglevel)
	log.SetFormatter(&easy.Formatter{TimestampFormat: "2006-01-02 15:04:05.00000", LogFormat: "%time% - [%lvl%] - %msg%\n"})
	log.Infof("Successfully set up loging. Level %s", loglevel)
}

func Setup(level string) {
	log.SetFormatter(&log.JSONFormatter{})
	l, err := log.ParseLevel(level)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(l)
}
