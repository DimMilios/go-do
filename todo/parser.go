package todo

import (
	"errors"
	"fmt"
	"log"
	_ "log"
	"regexp"
	"strings"
)

const (
	CHAR_A      string = "A"
	CHAR_B      string = "B"
	CHAR_C      string = "C"
	CHAR_D      string = "D"
	CHAR_E      string = "E"
	CHAR_F      string = "F"
	CHAR_G      string = "G"
	CHAR_H      string = "H"
	CHAR_I      string = "I"
	CHAR_J      string = "J"
	CHAR_K      string = "K"
	CHAR_L      string = "L"
	CHAR_M      string = "M"
	CHAR_N      string = "N"
	CHAR_O      string = "O"
	CHAR_P      string = "P"
	CHAR_Q      string = "Q"
	CHAR_R      string = "R"
	CHAR_S      string = "S"
	CHAR_T      string = "T"
	CHAR_U      string = "U"
	CHAR_V      string = "V"
	CHAR_W      string = "W"
	CHAR_X      string = "X"
	CHAR_Y      string = "Y"
	CHAR_Z      string = "Z"
	DONE_CHAR   string = "x"
	LEFT_PAREN  string = "("
	RIGHT_PAREN string = ")"
	PLUS        string = "+"
	AT          string = "@"
	COLON       string = ":"
	EOL         string = "EOL"
	STRING      string = "STRING"
)

type Token struct {
	tokenType string
	lexeme    string
	value     string
}

func (t Token) String() string {
	return fmt.Sprintf("{ tokenType: %s, lexeme: %s, value: %s }", t.tokenType, t.lexeme, t.value)
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
	/// Project, context or key-value.
	tagType TagType
	value   string

	/// A key exists if a tag is a key-value.
	key *string
}

func (t Tag) String() string {
	return fmt.Sprintf("{ tagType: %v, value: %s, key: %s }", t.tagType, t.value, *t.key)
}

type Description struct {
	/// Text content for the description.
	text string
	/// List of description's tags.
	tags []Tag
}

type Todo struct {
	/// Mandatory: Description + tags section of the todo.
	description Description

	/// Optional: Todo is complete. Can get 'x' as value
	done *rune
	/// Optional: The todo's priority is defined as a capital letter (A-Z)
	/// enclosed in parentheses, e.g., (A)
	priority *string

	/// Optional: Date the todo was created at (YYYY-MM-DD)
	creationDate *string
	/// Optional: Date the todo was completed (YYYY-MM-DD).
	/// Its existence is dependent on creationDate.
	completionDate *string
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

func scan(input string) []Token {
	current := 0 // current char
	var tokens []Token

	alphaRegex, _ := regexp.Compile(`[a-z]|[A-Z]`)

	for !isAtEnd(current, input) {
		char := string(input[current])
		fmt.Printf("Char: %s\n", char)
		switch char {
		case DONE_CHAR:
			tokens = append(tokens, Token{tokenType: DONE_CHAR})
		case LEFT_PAREN:
			tokens = append(tokens, Token{tokenType: LEFT_PAREN, value: string(input[current+1])})
		case RIGHT_PAREN:
			tokens = append(tokens, Token{tokenType: RIGHT_PAREN})
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
	todo := &Todo{}

	for _, token := range tokens {
		switch token.tokenType {
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
			colonPos := strings.Index(token.value, ":")
			keyToRemove := token.value[0:colonPos]

			keyStartPos := strings.Index(todo.description.text, keyToRemove)
			if keyStartPos >= 0 {
				todo.description.text = todo.description.text[0:keyStartPos]
			}
		}
	}

	// Handle todo completion
	if string(input[0]) == DONE_CHAR && len(input) > 1 && string(input[1]) == " " {
		x := rune(DONE_CHAR[0])
		todo.done = &x
	}

	return todo, nil
}

// func peek(current int, input string) rune {
// 	if isAtEnd(current, input) {
// 		return EOL
// 	}
// 	return rune(input[current])
// }
