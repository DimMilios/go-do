package todo

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

const YYYYMMDD = "2006-01-02"

type TokenType int

const (
	DONE_CHAR TokenType = iota
	LEFT_PAREN
	RIGHT_PAREN
	PLUS
	AT
	COLON
	DASH
	STRING
)

func (t TokenType) String() string {
	return [...]string{"x", "(", ")", "+", "@", ":", "-", "STRING"}[t]
}

type Token struct {
	tokenType TokenType
	value     string
}

func (t Token) String() string {
	return fmt.Sprintf("{ tokenType: %s, value: %s }", t.tokenType, t.value)
}

type TagType int8

const (
	Project TagType = iota
	Context
	KeyValue
)

func (t TagType) String() string {
	return [...]string{"Project", "Context", "KeyValue"}[t]
}

type Tag struct {
	// Project, context or key-value.
	TagType TagType
	Value   string

	// A Key exists if a tag is a Key-value.
	Key *string
}

func (t Tag) String() string {
	if t.Key != nil {
		return fmt.Sprintf("{ tagType: %v, value: %s, key: %s }", t.TagType, t.Value, *t.Key)
	}
	return fmt.Sprintf("{ tagType: %v, value: %s, key: %v }", t.TagType, t.Value, t.Key)
}

type Description struct {
	// Text content for the description.
	Text string
	// List of description's Tags.
	Tags []Tag
}

func (d Description) String() string {
	return fmt.Sprintf("{ text: %s, tags: %v }", d.Text, d.Tags)
}

type Todo struct {
	// Mandatory: Description + tags section of the todo.
	Description Description

	// Auto-generated: the original todo.txt string of this todo
	Original string

	// Optional: Todo is complete
	Done bool
	// Optional: The todo's Priority is defined as a capital letter (A-Z)
	// enclosed in parentheses, e.g., (A)
	Priority *string

	// Auto-generated: Date the todo was created at (YYYY-MM-DD)
	CreationDate time.Time
	// Optional: Date the todo was completed (YYYY-MM-DD).
	// Its existence is dependent on creationDate.
	CompletionDate *time.Time
}

func (t Todo) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "{ description: %v, ", t.Description)
	fmt.Fprintf(&b, "creationDate: %s, ", t.CreationDate.Format(YYYYMMDD))
	fmt.Fprintf(&b, "done: %v, ", t.Done)

	if t.Priority != nil {
		fmt.Fprintf(&b, "priority: %s, ", *t.Priority)
	}
	if t.CompletionDate != nil {
		fmt.Fprintf(&b, "completionDate: %s", t.CompletionDate.Format(YYYYMMDD))
	}

	b.WriteString(" }")
	return b.String()
}

func (t Todo) Format() string {
	var b strings.Builder

	if t.Done {
		fmt.Fprintf(&b, "x ")
	}

	if t.Priority != nil {
		fmt.Fprintf(&b, "(%s) ", *t.Priority)
	}

	if t.CompletionDate != nil {
		fmt.Fprintf(&b, "%s ", t.CompletionDate.Format(YYYYMMDD))
	}
	fmt.Fprintf(&b, "%s ", t.CreationDate.Format(YYYYMMDD))

	fmt.Fprintf(&b, "%s ", t.Description.Text)

	for _, t := range t.Description.Tags {
		switch t.TagType {
		case Context:
			fmt.Fprintf(&b, "@%s ", t.Value)
		case Project:
			fmt.Fprintf(&b, "+%s ", t.Value)
		case KeyValue:
			fmt.Fprintf(&b, "%s:%s ", *t.Key, t.Value)
		}
	}

	return strings.TrimRight(b.String(), " ")
}

func projectLiteral(current *int, input string) Token {
	var content string

	var startIndex = *current
	for ; !isAtEnd(*current, input); *current++ {
		if isWhiteSpace(*current, input) {
			break
		}
	}

	content += input[startIndex+1 : *current]
	log.Printf("Project Literal==== Content: %s, current: %v\n", content, *current)
	return Token{tokenType: PLUS, value: content}
}

func contextLiteral(current *int, input string) Token {
	var content string

	var startIndex = *current
	for ; !isAtEnd(*current, input); *current++ {
		if isWhiteSpace(*current, input) {
			break
		}
	}

	content += input[startIndex+1 : *current]
	log.Printf("Context Literal==== Content: %s, current: %v\n", content, *current)
	return Token{tokenType: AT, value: content}
}

func isAtEnd(current int, input string) bool {
	return current >= len(input)
}

func isWhiteSpace(current int, input string) bool {
	matched, err := regexp.MatchString(`(\s+)`, string(input[current]))
	if err != nil {
		log.Printf("regex for string %s with index %v\n", input, current)
	}
	return matched
}

func isCapitalLetter(current int, input string) bool {
	matched, err := regexp.MatchString(`[A-Z]{1}`, string(input[current]))
	if err != nil {
		log.Printf("regex for string %s with index %v\n", input, current)
	}
	return matched
}

func keyValueLiteral(current *int, input string) Token {
	colonPos := *current

	keyBegin := colonPos - 1
	for keyBegin >= 0 && !isWhiteSpace(keyBegin, input) {
		keyBegin--
	}
	keyBegin++

	i := colonPos + 1
	for !isAtEnd(i, input) {
		if isWhiteSpace(i, input) {
			break
		}
		i++
	}
	value := input[colonPos:i]
	*current += len(value)
	log.Printf("KV Token: { key: %s, value: %s }\n", input[keyBegin:colonPos], value)

	return Token{tokenType: COLON, value: input[keyBegin:colonPos] + value}
}

func handlePriority(current *int, input string) (*Token, error) {
	if string(input[*current+2]) != RIGHT_PAREN.String() || !isCapitalLetter(*current+1, input) {
		return nil, errors.New("bad priority value")
	}
	value := string(input[*current+1])
	*current += 3
	return &Token{tokenType: LEFT_PAREN, value: value}, nil
}

func handleDate(current *int, input string) (*Token, error) {
	pos := strings.Index(input, DASH.String())
	log.Printf("Pos starting value: %d\n", pos)
	// Parse the year backwards
	yearStart := pos - 4
	dateValue := ""
	if yearStart < 0 {
		return nil, errors.New("couldn't parse date")
	}

	year, month, day := "", "", ""
	year = input[yearStart:pos]
	pos++
	month = input[pos : pos+2]
	pos += 3
	day = input[pos : pos+2]
	pos += 3

	dateValue = year + "-" + month + "-" + day
	isValid, err := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, dateValue)
	// subtract 4 for year
	*current += len(dateValue) - 5
	log.Printf("Slice from new current value: %s, completion date: %s\n", input[*current:], dateValue)

	if isValid && err != nil {
		return nil, errors.New("bad format for date")
	}

	return &Token{tokenType: DASH, value: dateValue}, nil
}

func scan(input string) []Token {
	current := 0 // current char
	tokens := []Token{{tokenType: STRING, value: input}}

	for !isAtEnd(current, input) {
		char := string(input[current])
		switch char {
		case DONE_CHAR.String():
			tokens = append(tokens, Token{tokenType: DONE_CHAR})
		case LEFT_PAREN.String():
			token, err := handlePriority(&current, input)
			if err != nil {
				panic(err)
			}
			tokens = append(tokens, *token)
		case DASH.String():
			token, err := handleDate(&current, input)
			if err != nil {
				panic(err)
			}
			tokens = append(tokens, *token)
		case PLUS.String():
			tokens = append(tokens, projectLiteral(&current, input))
		case AT.String():
			tokens = append(tokens, contextLiteral(&current, input))
		case COLON.String():
			tokens = append(tokens, keyValueLiteral(&current, input))
		default:
		}

		current++
	}
	return tokens
}

func handleKeyValueTag(token Token) (*Tag, error) {
	colonPos := strings.Index(token.value, ":")
	if colonPos < 0 {
		log.Printf("String doesn't contain colon. String: %s", token.value)
		return nil, errors.New("colon character not found in string")
	}
	key := token.value[0:colonPos]
	value := token.value[colonPos+1:]
	if value == "" || value == ":" {
		return nil, errors.New("value cannot be empty")
	}

	kvTag := Tag{TagType: KeyValue, Key: &key, Value: value}
	log.Printf("kvTag: %v\n", kvTag)
	return &kvTag, nil
}

func stripRight(desc string, val string, input string) string {
	var b strings.Builder
	pos := strings.Index(desc, val)
	if pos >= 0 {
		before := desc[:pos-1]
		b.WriteString(before)
	} else {
		b.WriteString(desc)
	}
	return strings.TrimSpace(b.String())
}

func stripLeft(desc string, val string, input string) string {
	var b strings.Builder
	pos := strings.Index(desc, val)
	if pos >= 0 {
		afterLen := len(val) + 1
		if afterLen < len(input) {
			after := input[afterLen:]
			b.WriteString(after)
		}
	} else {
		b.WriteString(desc)
	}
	return strings.TrimSpace(b.String())
}

func Parse(input string) (*Todo, error) {
	log.Printf("Got: %s\n", input)
	input = strings.Trim(input, " ")

	todo := &Todo{
		Done:         false,
		CreationDate: time.Now().UTC(),
		Original:     input,
	}

	// Handle todo completion
	if string(input[0]) == DONE_CHAR.String() && len(input) > 1 && string(input[1]) == " " {
		todo.Done = true
		input = input[2:]
	}
	tokens := scan(input)

	for _, token := range tokens {
		switch token.tokenType {
		case LEFT_PAREN:
			tmp := token.value
			todo.Priority = &tmp
		case STRING:
			todo.Description.Text = token.value
		case PLUS:
			todo.Description.Tags = append(todo.Description.Tags, Tag{TagType: Project, Value: token.value})
		case AT:
			todo.Description.Tags = append(todo.Description.Tags, Tag{TagType: Context, Value: token.value})
		case COLON:
			kvTag, err := handleKeyValueTag(token)
			if err != nil {
				return nil, err
			}
			todo.Description.Tags = append(todo.Description.Tags, *kvTag)

			// We have to clean up the description text since
			// it's picking up on the key of the first key value tag
			colonPos := strings.Index(token.value, COLON.String())
			keyToRemove := token.value[0:colonPos]

			keyStartPos := strings.Index(todo.Description.Text, keyToRemove)
			if keyStartPos >= 0 {
				todo.Description.Text = todo.Description.Text[0:keyStartPos]
			}
		case DASH:
			date, err := time.Parse(YYYYMMDD, token.value)
			log.Printf("date: %s\n", date.Format(YYYYMMDD))
			if err != nil {
				log.Printf("bad date: %s\n", token.value)
				return nil, errors.New("could not parse completion date")
			}
			todo.Description.Text = stripLeft(todo.Description.Text, token.value, input)
			todo.CompletionDate = &date
		}
	}

	if todo.Priority != nil {
		todo.Description.Text = todo.Description.Text[4:]
	}

	for _, t := range todo.Description.Tags {
		todo.Description.Text = stripRight(todo.Description.Text, t.Value, input)
	}

	return todo, nil
}
