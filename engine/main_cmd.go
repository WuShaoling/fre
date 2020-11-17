package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mydocker/api"
	"github.com/urfave/cli"
	"os"
)

const usage = "fre is a simple container runtime implementation for serverless computing."

func main() {
	app := cli.NewApp()
	app.Name = "fre"
	app.Usage = usage

	app.Commands = []cli.Command{
		api.InitCommand,
		api.RunCommand,
		api.ListCommand,
		api.LogCommand,
		api.ExecCommand,
		api.StopCommand,
		api.RemoveCommand,
		api.CommitCommand,
	}

	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})

		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
