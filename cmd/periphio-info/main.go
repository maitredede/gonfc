package main

import (
	goflag "flag"

	"github.com/maitredede/gonfc/cmd/common"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"periph.io/x/host/v3"
	_ "periph.io/x/host/v3/ftdi"
	"periph.io/x/host/v3/rpi"
)

var logger *zap.SugaredLogger

func main() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	log := common.InitLogger(true)
	defer log.Sync()
	redir := zap.RedirectStdLog(log)
	defer redir()

	logger = log.Sugar()

	state, err := host.Init()
	if err != nil {
		logger.Fatal(err)
	}
	if rpi.Present() {
		logger.Infof("This is a raspberry pi board")
	} else {
		logger.Infof("Not a raspberry pi board")
	}
	logger.Infof("%+v", state)
}
