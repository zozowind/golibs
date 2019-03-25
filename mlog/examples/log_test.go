package exmaples

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

var logger *Logger

func init() {
	setting1 := &LogSetting{
		Write: zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true, TimeFormat: time.RFC3339},
		Level: zerolog.DebugLevel,
	}
	setting2 := &LogSetting{
		Write: os.Stdout,
		Level: zerolog.InfoLevel,
	}
	logger = New([]*LogSetting{setting1, setting2})
}

func TestLog(t *testing.T) {
	var err error
	defer func() {
		if nil != err {
			logger.ErrorKv(map[string]interface{}{
				"requestid": "1111",
			}, "test a error %s", err.Error())
		}
	}()
	err = fmt.Errorf("error test")
	return
}
