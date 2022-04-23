package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	todos "github.com/go-do/todo"
	"github.com/urfave/cli/v2"
)

func init() {
	f, err := os.OpenFile(time.Now().Format(todos.YYYYMMDD)+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("couldn't open log file")
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(f)
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
					t, err := todos.Parse(c.Args().First())
					log.Println(t.Original)
					if err != nil {
						return err
					}

					todos.AddToFile(t)
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
						if len(value) <= 0 {
							fmt.Println("you have to provide a value when passing in a tag")
							return errors.New("you have to provide a value when passing in a tag")
						}

						switch strings.ToLower(tag) {
						case strings.ToLower(todos.Project.String()):
							todos.PrintByTag(todos.Project, value)
						case strings.ToLower(todos.Context.String()):
							todos.PrintByTag(todos.Context, value)
						case strings.ToLower(todos.KeyValue.String()):
							todos.PrintByKVTag(value)
						default:
							return errors.New("viable tag values are one of project, context or keyvalue")
						}
					} else {
						todos.PrintAll()
					}
					return nil
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "Delete a todo",
				Action: func(c *cli.Context) error {
					if c.Args().Len() > 0 {
						f, err := os.Open("todos-copy.txt")
						if err != nil {
							fmt.Println("couldn't open file")
						}
						todos.DeleteFirst(f, c.Args().First())
					} else {
						fmt.Println("Please, provide a description for the todo to be deleted.")
						log.Println("Couldn't parse todo.")
					}
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
