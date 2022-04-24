package todo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

func GetFromFile(r io.Reader) ([]string, error) {
	var b strings.Builder
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return strings.Split(b.String(), "\n"), nil
}

func GetAll(r io.Reader) ([]*Todo, error) {
	lines, err := GetFromFile(r)
	if err != nil {
		return nil, err
	}

	todos := make([]*Todo, 0)
	for _, l := range lines {
		t, err := Parse(l)
		if err != nil {
			log.Println(err)
			return todos, err
		}
		todos = append(todos, t)

	}
	return todos, nil
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

func FindLineByText(lines []string, text string) int {
	for i, l := range lines {
		if strings.Contains(strings.ToLower(l), strings.ToLower(text)) {
			return i
		}
	}
	return -1
}

func SkipFirst(r io.Reader, text string) ([]string, error) {
	var b strings.Builder
	scanner := bufio.NewScanner(r)
	found := false
	for scanner.Scan() {
		s := scanner.Text()
		if !found && strings.Contains(strings.ToLower(s), strings.ToLower(text)) {
			// Skip the first occurrence of this string
			found = true
			continue
		}
		b.WriteString(s + "\n")
	}

	if !found {
		return nil, errors.New("text was not found")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	if len(b.String()) > 0 {
		return strings.Split(b.String(), "\n"), nil
	}

	return nil, errors.New("couldn't write lines")
}

func WriteAll(w io.Writer, lines []string) error {
	for _, fl := range lines {
		if _, err := fmt.Fprintln(w, fl); err != nil {
			return err
		}
	}
	return nil
}

func DeleteFirst(file io.Reader, text string) error {
	lines, err := SkipFirst(file, text)
	if len(lines) < 1 || err != nil {
		fmt.Println("couldn't parse todo file")
		log.Fatal("couldn't parse todo file")
	}

	if err != nil {
		return err
	}

	// Write lines to new tmp file skipping the first occurrence
	tmp, err := os.Create("copy.tmp")
	if err != nil {
		panic(err)
	}
	defer tmp.Close()

	if err := WriteAll(tmp, lines); err != nil {
		log.Fatal(err)
	}

	// Replace old file with new file
	if osf, ok := file.(*os.File); ok {
		if err := os.Rename("copy.tmp", osf.Name()); err != nil {
			return err
		}
	}

	if err := tmp.Close(); err != nil {
		log.Fatal(err)
	}

	return nil
}
