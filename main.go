package main

import (
    "github.com/codegangsta/cli"

    "github.com/nbremond/double-kill/commands"
    "github.com/nbremond/double-kill/modules/settings"
)

func init() {
    settings.Version = "0.0.0.1"
}

func main() {
    app := cli.NewApp()
    app.Name = "double-kill"
    app.Usage = "double-kill"
    app.Version = settings.Version
    app.Commands = []cli.Command{
        commands.CmdSearch,
    }
    app.Flags = append(app.Flags, []cli.Flag{}...)
    app.RunAndExitOnError()
}
