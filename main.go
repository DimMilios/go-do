package main

import (
	"log"
	"os"

	"github.com/go-do/todo"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "create",
				Aliases:  []string{"c"},
				Usage:    "`Todo` value based on todo.txt format",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			if c.String("create") != "" {
				todo, err := todo.Parse(c.String("create"))
				log.Println(todo)
				if err != nil {
					return err
				}

			} else {
				log.Println("Couldn't parse todo.")
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
