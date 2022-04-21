package todo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const defaultFileName = "todos.txt"

func AddToFile(todo *Todo, fileName *string) {
	fname := defaultFileName
	if fileName != nil {
		fname = *fileName
	}

	f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, todo.Format())
	if err != nil {
		log.Println(err)
	}
}

func PrintAll(fileName *string) {
	fname := defaultFileName
	if fileName != nil {
		fname = *fileName
	}

	f, err := os.Open(fname)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println("Todo: ", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func PrintByTag(tag TagType, value string) {
	fname := defaultFileName
	f, err := os.Open(fname)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	var b strings.Builder
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		b.WriteString(scanner.Text())
		// fmt.Println("Todo: ", scanner.Text())
	}

	lines := strings.Split(b.String(), "\n")
	todos := make([]*Todo, 0)
	for _, l := range lines {
		t, err := Parse(l)
		if err != nil {
			log.Println(err)
		}
		todos = append(todos, t)
	}

	fmt.Println(todos)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
