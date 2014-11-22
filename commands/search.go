package commands

import (
    "fmt"

    "github.com/codegangsta/cli"
)

var CmdSearch = cli.Command{
    Name:  "search",
    Usage: "Search for duplicates",
    Description: `Search for duplicates`,
    Before: runSearch,
    Action: func(ctx *cli.Context) {},
    Flags:  []cli.Flag{},
}

func runSearch(ctx *cli.Context) error {
    fmt.Println("It's working!")

    return nil
}
