package todo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const defaultFileName = "todos.txt"

func fileOrDefault(name string) string {
	fname := defaultFileName
	if len(name) > 0 {
		fname = name
	}
	return fname
}

func AddToFile(todo *Todo) {
	fname := fileOrDefault("")
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, todo.Original)
	if err != nil {
		log.Println(err)
	}
}

func PrintAll() {
	fname := fileOrDefault("")
	f, err := os.Open(fname)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func Contains(todos []*Todo, todo *Todo) bool {
	for _, t := range todos {
		if t == todo {
			return true
		}
	}
	return false
}

func todosFromFileByValue(f *os.File, value string) []*Todo {
	var b strings.Builder
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		b.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(b.String(), "\n")
	todos := make([]*Todo, 0)
	for _, l := range lines {
		if !strings.Contains(strings.ToLower(l), value) {
			// todo doesn't contain this tag value
			continue
		}

		t, err := Parse(l)
		if err != nil {
			log.Println(err)
		}
		todos = append(todos, t)

	}
	return todos
}

func PrintByTag(tag TagType, value string) {
	fname := fileOrDefault("")
	f, err := os.Open(fname)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	formattedVal := strings.ToLower(value)
	todos := todosFromFileByValue(f, formattedVal)

	filteredTodos := make([]*Todo, 0)
	for _, todo := range todos {
		for _, tg := range todo.Description.Tags {
			if tg.TagType == tag && strings.Contains(strings.ToLower(tg.Value), formattedVal) && !Contains(filteredTodos, todo) {
				filteredTodos = append(filteredTodos, todo)
			}
		}
	}

	for _, f := range filteredTodos {
		fmt.Println(f.Original)
	}
}

func PrintByKVTag(key string) {
	fname := fileOrDefault("")
	f, err := os.Open(fname)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	formattedKey := strings.ToLower(key)
	todos := todosFromFileByValue(f, formattedKey)

	filteredTodos := make([]*Todo, 0)
	for _, todo := range todos {
		for _, tg := range todo.Description.Tags {
			if tg.TagType == KeyValue && strings.Contains(strings.ToLower(*tg.Key), formattedKey) && !Contains(filteredTodos, todo) {
				filteredTodos = append(filteredTodos, todo)
			}
		}
	}

	for _, f := range filteredTodos {
		fmt.Println(f.Original)
	}
}

func FindByDescrText(todos []Todo, text string) *Todo {
	for _, t := range todos {
		if strings.Contains(strings.ToLower(t.Description.Text), strings.ToLower(text)) {
			return &t
		}
	}
	return nil
}
