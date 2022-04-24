package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	todos "github.com/go-do/todo"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

func confirmDeletion(fname, result string) error {
	pr := promptui.Prompt{
		Label:     fmt.Sprintf("Delete %q", result),
		IsConfirm: true,
	}

	if _, err := pr.Run(); err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}
	fmt.Printf("You deleted %q\n", result)

	// File is somehow getting "consumed" when reading it to build the
	// selection list, so we have to re-open it here
	f, err := os.Open(fname)
	if err != nil {
		fmt.Println("couldn't open file")
		return err
	}
	defer f.Close()

	todos.DeleteFirst(f, result)
	return nil
}

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
						if len(c.Args().First()) < 1 {
							fmt.Println("Todo description cannot be empty.")
							log.Fatal("passed empty description")
						}

						f, err := os.Open("todos-copy.txt")
						if err != nil {
							fmt.Println("couldn't open file")
						}
						defer f.Close()
						err = todos.DeleteFirst(f, c.Args().First())
						if err != nil {
							fmt.Printf("Failed to delete todo with text: %q", c.Args().First())
						}
						fmt.Printf("You deleted %q", c.Args().First())
					} else {
						fmt.Println("Please, provide a description for the todo to be deleted.")
						log.Println("Couldn't parse todo: empty description")
					}
					return nil
				},
			},
			{
				Name:    "delete-by-select",
				Aliases: []string{"ds"},
				Usage:   "Select a todo to delete by listing all todos",
				Action: func(c *cli.Context) error {
					fname := "todos-copy.txt"
					f, err := os.Open(fname)
					if err != nil {
						fmt.Println("couldn't open file")
					}
					defer f.Close()

					lines, _ := todos.GetFromFile(f)

					prompt := promptui.Select{
						Size:  len(lines),
						Label: "Select Todo",
						Items: lines,
					}
					_, result, _ := prompt.Run()
					confirmDeletion(fname, result)
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
