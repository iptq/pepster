package pepster

import (
	"log"
	"os"
)

type Reporter struct {
	pepster *Pepster
}

func NewLogger(pepster *Pepster) *log.Logger {
	return log.New(Reporter{
		pepster: pepster,
	}, "", log.Ldate|log.Ltime)
}

func (logger Reporter) Write(p []byte) (n int, err error) {
	os.Stdout.Write(p)
	msg, err := logger.pepster.dg.ChannelMessageSend(logger.pepster.conf.LogChannel, string(p))
	if err != nil {
		return -1, err
	}
	return len(msg.Content), nil
}
