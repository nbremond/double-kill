package main

import (
    "github.com/codegangsta/cli"

    "github.com/nbremond/double-kill/commands"
)

const APP_VER = "0.0.1"

func init() {
}

func main() {
    app := cli.NewApp()
    app.Name = "double-kill"
    app.Usage = "double-kill"
    app.Version = APP_VER
    app.Commands = []cli.Command{
        commands.CmdSearch,
    }
    app.Flags = append(app.Flags, []cli.Flag{}...)
    app.RunAndExitOnError()
}
