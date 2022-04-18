package todo

import (
	"errors"
	"fmt"
	"log"
	_ "log"
	"regexp"
	"strings"
	"time"
)

const (
	CHAR_A      = "A"
	CHAR_B      = "B"
	CHAR_C      = "C"
	CHAR_D      = "D"
	CHAR_E      = "E"
	CHAR_F      = "F"
	CHAR_G      = "G"
	CHAR_H      = "H"
	CHAR_I      = "I"
	CHAR_J      = "J"
	CHAR_K      = "K"
	CHAR_L      = "L"
	CHAR_M      = "M"
	CHAR_N      = "N"
	CHAR_O      = "O"
	CHAR_P      = "P"
	CHAR_Q      = "Q"
	CHAR_R      = "R"
	CHAR_S      = "S"
	CHAR_T      = "T"
	CHAR_U      = "U"
	CHAR_V      = "V"
	CHAR_W      = "W"
	CHAR_X      = "X"
	CHAR_Y      = "Y"
	CHAR_Z      = "Z"
	DONE_CHAR   = "x"
	LEFT_PAREN  = "("
	RIGHT_PAREN = ")"
	PLUS        = "+"
	AT          = "@"
	COLON       = ":"
	DASH        = "-"
	EOL         = "EOL"
	STRING      = "STRING"

	YYYYMMDD = "2006-01-02"
)

type Token struct {
	tokenType string
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
	tagType TagType
	value   string

	// A key exists if a tag is a key-value.
	key *string
}

func (t Tag) String() string {
	return fmt.Sprintf("{ tagType: %v, value: %s, key: %s }", t.tagType, t.value, *t.key)
}

type Description struct {
	// Text content for the description.
	text string
	// List of description's tags.
	tags []Tag
}

type Todo struct {
	// Mandatory: Description + tags section of the todo.
	description Description

	// Optional: Todo is complete. Can get 'x' as value
	done *rune
	// Optional: The todo's priority is defined as a capital letter (A-Z)
	// enclosed in parentheses, e.g., (A)
	priority *string

	// Auto-generated: Date the todo was created at (YYYY-MM-DD)
	creationDate time.Time
	// Optional: Date the todo was completed (YYYY-MM-DD).
	// Its existence is dependent on creationDate.
	completionDate *time.Time
}

func (t Todo) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "{ description: %v, ", t.description)
	fmt.Fprintf(&b, "creationDate: %s, ", t.creationDate.Format(YYYYMMDD))

	if t.done != nil {
		fmt.Fprintf(&b, "done: %c, ", *t.done)
	}
	if t.priority != nil {
		fmt.Fprintf(&b, "priority: %s, ", *t.priority)
	}
	if t.completionDate != nil {
		fmt.Fprintf(&b, "completionDate: %s", t.completionDate.Format(YYYYMMDD))
	}

	b.WriteString(" }")
	return b.String()
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
	fmt.Printf("Project Literal==== Content: %s, current: %v\n", content, *current)
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
	fmt.Printf("Context Literal==== Content: %s, current: %v\n", content, *current)
	return Token{tokenType: AT, value: content}
}

func literal(current *int, input string) Token {
	var content string
	var startIndex = *current
	for ; !isAtEnd(*current, input); *current++ {
		if isTagChar(*current, input) {
			break
		}
	}

	content += input[startIndex:*current]
	fmt.Printf("String Literal==== Content: %s, current: %v\n", content, *current)
	return Token{tokenType: STRING, value: content}
}

func isAtEnd(current int, input string) bool {
	return current >= len(input)
}

func isWhiteSpace(current int, input string) bool {
	matched, err := regexp.MatchString(`(\s+)`, string(input[current]))
	if err != nil {
		fmt.Printf("regex for string %s with index %v\n", input, current)
	}
	return matched
}

func isTagChar(current int, input string) bool {
	asStr := string(input[current])
	return asStr == AT || asStr == PLUS || asStr == COLON
}

func isCapitalLetter(current int, input string) bool {
	matched, err := regexp.MatchString(`[A-Z]{1}`, string(input[current]))
	if err != nil {
		fmt.Printf("regex for string %s with index %v\n", input, current)
	}
	return matched
}

func isNumeric(current int, input string) bool {
	matched, err := regexp.MatchString(`[0-9]`, string(input[current]))
	if err != nil {
		fmt.Printf("regex for string %s with index %v\n", input, current)
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
	fmt.Printf("KV Token: { key: %s, value: %s }\n", input[keyBegin:colonPos], value)

	return Token{tokenType: COLON, value: input[keyBegin:colonPos] + value}
}

func handlePriority(current *int, input string) (*Token, error) {
	if string(input[*current+2]) != RIGHT_PAREN || !isCapitalLetter(*current+1, input) {
		return nil, errors.New("bad priority value")
	}
	value := string(input[*current+1])
	*current += 3
	return &Token{tokenType: LEFT_PAREN, value: value}, nil
}

func handleCompletionDate(current *int, input string) (*Token, error) {
	pos := strings.Index(input, DASH)
	// Parse the year backwards
	yearStart := pos - 4
	year, month, day := "", "", ""
	if yearStart >= 0 {
		year = input[yearStart:pos]
		pos++
		month = input[pos : pos+2]
		pos += 3
		day = input[pos : pos+2]
		pos += 3

		_, yerr := regexp.MatchString(`[0-9]{4}`, year)
		_, merr := regexp.MatchString(`[0-9]{2}`, month)
		_, derr := regexp.MatchString(`[0-9]{2}`, day)

		if yerr != nil || merr != nil || derr != nil {
			return nil, errors.New("bad format for date")
		}
	}

	*current += pos

	fmt.Printf("year: %s, month: %s, day: %s\n", year, month, day)

	return &Token{tokenType: DASH, value: year + "-" + month + "-" + day}, nil
}

func scan(input string) []Token {
	current := 0 // current char
	var tokens []Token

	alphaRegex, _ := regexp.Compile(`\p{L}`)

	for !isAtEnd(current, input) {
		char := string(input[current])
		fmt.Printf("Char: %s\n", char)
		switch char {
		case DONE_CHAR:
			tokens = append(tokens, Token{tokenType: DONE_CHAR})
		case LEFT_PAREN:
			token, err := handlePriority(&current, input)
			if err != nil {
				panic(err)
			}
			tokens = append(tokens, *token)
		case DASH:
			token, err := handleCompletionDate(&current, input)
			fmt.Printf("token %v\n", token)
			if err != nil {
				panic(err)
			}
			tokens = append(tokens, *token)
		case PLUS:
			tokens = append(tokens, projectLiteral(&current, input))
		case AT:
			tokens = append(tokens, contextLiteral(&current, input))
		case COLON:
			tokens = append(tokens, keyValueLiteral(&current, input))
		case " ", "\r", "\t", "\n":
			break
		default:
			isAlpha, _ := regexp.MatchString(alphaRegex.String(), string(input[current]))
			if isAlpha {
				tokens = append(tokens, literal(&current, input))
				// Move one step back to include the next tag character
				current--
			} else if isNumeric(current, input) {
				current++
				continue
			} else {
				log.Fatalf("Unexpected character: %s", char)
			}
		}

		current++
	}
	return tokens
}

func handleKeyValueTag(token Token) (*Tag, error) {
	colonPos := strings.Index(token.value, ":")
	if colonPos < 0 {
		fmt.Printf("String doesn't contain colon. String: %s", token.value)
		return nil, errors.New("colon character not found in string")
	}
	key := token.value[0:colonPos]
	value := token.value[colonPos+1:]
	if value == "" || value == ":" {
		return nil, errors.New("value cannot be empty")
	}

	kvTag := Tag{tagType: KeyValue, key: &key, value: value}
	fmt.Printf("kvTag: %v\n", kvTag)
	return &kvTag, nil
}

func Parse(input string) (*Todo, error) {
	tokens := scan(input)
	todo := &Todo{
		creationDate: time.Now().UTC(),
	}

	for _, token := range tokens {
		switch token.tokenType {
		case LEFT_PAREN:
			tmp := token.value
			todo.priority = &tmp
		case STRING:
			todo.description.text = token.value
		case PLUS:
			todo.description.tags = append(todo.description.tags, Tag{tagType: Project, value: token.value})
		case AT:
			todo.description.tags = append(todo.description.tags, Tag{tagType: Context, value: token.value})
		case COLON:
			kvTag, err := handleKeyValueTag(token)
			if err != nil {
				return nil, err
			}
			todo.description.tags = append(todo.description.tags, *kvTag)

			// We have to clean up the description text since
			// it's picking up on the key of the first key value tag
			colonPos := strings.Index(token.value, COLON)
			keyToRemove := token.value[0:colonPos]

			keyStartPos := strings.Index(todo.description.text, keyToRemove)
			if keyStartPos >= 0 {
				todo.description.text = todo.description.text[0:keyStartPos]
			}
		case DASH:
			date, err := time.Parse(YYYYMMDD, token.value)
			fmt.Printf("date: %s\n", date.Format(YYYYMMDD))
			if err != nil {
				fmt.Printf("bad date: %s\n", token.value)
				return nil, errors.New("could not parse completion date")
			}
			todo.completionDate = &date
		}
	}

	// Handle todo completion
	if string(input[0]) == DONE_CHAR && len(input) > 1 && string(input[1]) == " " {
		x := rune(DONE_CHAR[0])
		todo.done = &x
	}

	return todo, nil
}
