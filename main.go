package main

import (
	"errors"
	"log"
	"os"
	"strings"

	todo "github.com/go-do/todo"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	var tag, value string

	app := &cli.App{
		Commands: []*cli.Command{{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "`Todo` value based on todo.txt format",
			Action: func(c *cli.Context) error {
				if c.Args().Len() > 0 {
					t, err := todo.Parse(c.Args().First())
					log.Println(t.Format())
					if err != nil {
						return err
					}

					todo.AddToFile(t, nil)
				} else {
					log.Println("Couldn't parse todo.")
				}
				return nil
			},
		},
			{
				Name:  "show",
				Usage: "Show all saved todos",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "tag", Aliases: []string{"t"}, Destination: &tag},
					&cli.StringFlag{Name: "value", Aliases: []string{"v"}, Destination: &value},
				},
				Action: func(c *cli.Context) error {
					if len(tag) > 0 {
						log.Println(tag)
						if len(value) == 0 {
							return errors.New("you have to provide a value when passing in a tag")
						}
						log.Println(value)

						switch strings.ToLower(tag) {
						case strings.ToLower(todo.Project.String()):
							todo.PrintByTag(todo.Project, value)
						case strings.ToLower(todo.Context.String()):
						case strings.ToLower(todo.KeyValue.String()):
						default:
							return errors.New("viable tag values are one of project, context or keyvalue")
						}
					}
					todo.PrintAll(nil)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
