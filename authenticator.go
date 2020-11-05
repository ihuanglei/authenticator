package main

import (
	"fmt"
	"os"

	"net/http"
	_ "net/http/pprof"

	"github.com/ihuanglei/authenticator/models"
	"github.com/ihuanglei/authenticator/pkg/build"
	"github.com/ihuanglei/authenticator/pkg/config"
	"github.com/ihuanglei/authenticator/pkg/logger"
	"github.com/ihuanglei/authenticator/pkg/web"

	"github.com/urfave/cli"
)

func main() {

	// pprof 监控
	go func() {
		http.ListenAndServe("0.0.0.0:16060", nil)
	}()

	app := cli.NewApp()
	app.Name = "authenticator"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf,c",
			Value: "auth.simple.yml",
			Usage: "configuration `file`",
		},
	}
	app.Action = run
	err := app.Run(os.Args)
	if err != nil {
		logger.Panic("Startup error!!!", err)
	}
}

func run(cli *cli.Context) {
	logger.Infof(fmt.Sprintf("Start authenticator build time %v", build.BuildTime))
	config, err := config.Load(cli.String("c"))
	if err != nil {
		logger.Fatal("Load Config error!!!", err)
		return
	}
	logger.SetLevel(config.Log)
	if err := models.Init(config); err != nil {
		logger.Fatal("Connect Database error!!!", err)
		return
	}
	web.Run(config)
}
