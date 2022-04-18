package todo

import (
	"fmt"
	"strings"
	"testing"
)

func Test_Parse_Simple_Description(t *testing.T) {
	expected := "simple description"
	todo, _ := Parse(expected)
	got := todo.description.text

	if strings.Compare(got, expected) != 0 {
		t.Fatalf("Description text is incorrect. Expected: \"%s\", but got: \"%s\"\n", expected, got)
	}
}

func Test_Parse_Simple_Description_In_Greek(t *testing.T) {
	expected := "απλή περιγραφή"
	todo, _ := Parse(expected)
	got := todo.description.text

	if strings.Compare(got, expected) != 0 {
		t.Fatalf("Description text is incorrect. Expected: \"%s\", but got: \"%s\"\n", expected, got)
	}
}

func Test_Parse_Should_Ignore_Done(t *testing.T) {
	simple := "x simple description"
	expected := "simple description"
	todo, _ := Parse(simple)
	got := todo.description.text

	if strings.Compare(got, expected) != 0 {
		t.Fatalf("Description text is incorrect. Expected: \"%s\", but got: \"%s\"\n", expected, got)
	}
}

func Test_Parse_Mark_As_Done(t *testing.T) {
	input := "x walk dog"
	todo, _ := Parse(input)

	got := string(*todo.done)

	if todo.done == nil || got != DONE_CHAR {
		t.Fatalf("Completion character at index 0 failed to mark todo as done. Got: %s\n", got)
	}
}

func Test_Parse_Mark_As_Done_Correct_Index(t *testing.T) {
	input := "walk x dog"
	todo, _ := Parse(input)

	if todo.done != nil {
		t.Fatal("Should not consider completion char at wrong index as valid.")
	}
}

func Test_Parse_Mark_As_Done_Is_Followed_By_Space(t *testing.T) {
	input := "xx walk dog"
	todo, _ := Parse(input)

	if todo.done != nil {
		t.Fatal("Completion char should be followed by a space.")
	}
}

func Test_Parse_Description_With_Project_Tag(t *testing.T) {
	input := "call customer +proj1"
	expected := "proj1"
	todo, _ := Parse(input)

	var projTag Tag
	for _, tag := range todo.description.tags {
		if tag.tagType == Project {
			projTag = tag
			break
		}
	}

	if projTag.value != expected {
		fmt.Printf("Project tag: %v\n", projTag)
		t.Fatalf("Couldn't parse project tag. Expected: \"%s\", but got: \"%s\"\n", expected, projTag.value)
	}
}

func Test_Parse_Description_With_Two_Project_Tags(t *testing.T) {
	input := "call customer +proj1 +proj2"
	todo, _ := Parse(input)

	var projTags []Tag
	for _, tag := range todo.description.tags {
		if tag.tagType == Project {
			projTags = append(projTags, tag)
		}
	}

	if len(projTags) != 2 {
		fmt.Printf("Todo: %v\n", todo)
		t.Fatalf("Couldn't parse many project tags. Tags: %v", projTags)
	}
}

func Test_Parse_Description_With_Context_Tag(t *testing.T) {
	input := "call customer @ctx1"
	expected := "ctx1"
	todo, _ := Parse(input)

	var ctxTag Tag
	for _, tag := range todo.description.tags {
		if tag.tagType == Context {
			ctxTag = tag
		}
	}

	if ctxTag.value != expected {
		t.Fatalf("Couldn't parse context tag. Expected: \"%s\", but got: \"%s\"\n", expected, ctxTag.value)
	}
}

func Test_Parse_Description_With_Two_Context_Tags(t *testing.T) {
	input := "call customer @ctx1 @ctx2"
	todo, _ := Parse(input)

	var ctxTags []Tag
	for _, tag := range todo.description.tags {
		if tag.tagType == Context {
			ctxTags = append(ctxTags, tag)
		}
	}

	if len(ctxTags) != 2 {
		fmt.Printf("Todo: %v\n", todo)
		t.Fatalf("Couldn't parse many context tags. Tags: %v", ctxTags)
	}
}

func Test_Parse_Description_With_Project_And_Context_Tags(t *testing.T) {
	input := "call customer +proj @ctx1"
	todo, _ := Parse(input)

	var tags []Tag
	for _, tag := range todo.description.tags {
		if tag.tagType == Project || tag.tagType == Context {
			tags = append(tags, tag)
		}
	}

	if len(tags) != 2 {
		fmt.Printf("Todo: %v\n", todo)
		t.Fatalf("Couldn't parse mixed project and context tags. Tags: %v", tags)
	}
}

func Test_Parse_Mark_Done_Complex_Description(t *testing.T) {
	input := "x call customer +proj @ctx1"
	todo, _ := Parse(input)

	var tags []Tag
	for _, tag := range todo.description.tags {
		if tag.tagType == Project || tag.tagType == Context {
			tags = append(tags, tag)
		}
	}

	if len(tags) != 2 || todo.done == nil {
		t.Fatalf("Couldn't mark todo with tagged description as done. Todo: %v", todo)
	}
}

func Test_Parse_Description_With_Key_Value_Tag(t *testing.T) {
	input := "call customer due:now"
	todo, _ := Parse(input)
	expectedKey := "due"
	expectedVal := "now"

	var keyValueTag Tag
	for _, tag := range todo.description.tags {
		if tag.key != nil && *tag.key == expectedKey && tag.value == expectedVal {
			keyValueTag = tag
		}
	}

	if keyValueTag.key == nil || *keyValueTag.key != expectedKey || keyValueTag.value != expectedVal {
		t.Fatalf("Couldn't parse key value tag. Tag: %v", keyValueTag)
	}
}

func Test_Parse_Key_Value_Tag_Empty_Value(t *testing.T) {
	input := "call customer due:"
	_, err := Parse(input)

	if err == nil {
		t.Fatal("Trying to pass a key value tag without suppling a value should return an error.")
	}
}

func Test_Parse_Description_With_Many_Kv_Tags(t *testing.T) {
	input := "call customer due:now who:me test:ing"
	todo, _ := Parse(input)

	var kvTags []Tag
	for _, tag := range todo.description.tags {
		if tag.tagType == KeyValue {
			kvTags = append(kvTags, tag)
		}
	}

	if len(kvTags) != 3 {
		fmt.Printf("Todo: %v\n", todo)
		t.Fatalf("Couldn't parse many key value tags. Tags: %v", kvTags)
	}
}

func Test_Description_Doesnt_Contain_Key_Of_Kv_Tag(t *testing.T) {
	input := "call customer due:now"
	expected := "call customer "
	todo, _ := Parse(input)

	if todo == nil || todo.description.text != expected {
		t.Fatal("Key of key value tag shouldn't be included to the text description.")
	}
}

func Test_Parse_Description_With_All_Tags(t *testing.T) {
	input := "call customer +proj @ctx1 due:now"
	todo, _ := Parse(input)

	if len(todo.description.tags) != 3 {
		fmt.Printf("Todo: %v\n", todo)
		t.Fatalf("Couldn't parse mixed types of tags. Tags: %v", todo.description.tags)
	}
}

func Test_Parse_Description_With_All_Tags_Reordered(t *testing.T) {
	input := "call customer due:now @ctx1 +proj1"
	todo, _ := Parse(input)

	if len(todo.description.tags) != 3 {
		fmt.Printf("Todo: %v\n", todo)
		t.Fatalf("Couldn't parse mixed types of tags. Tags: %v", todo.description.tags)
	}
}

func Test_Parse_Description_With_Multiple_Of_Each_Tag(t *testing.T) {
	input := "call customer due:now @ctx1 who:john +proj1 +proj2 @ctx2"
	todo, _ := Parse(input)

	if len(todo.description.tags) != 6 {
		fmt.Printf("Todo: %v\n", todo)
		t.Fatalf("Couldn't parse mixed types of tags. Tags: %v", todo.description.tags)
	}
}

func Test_Parse_Priority(t *testing.T) {
	input := "(A) simple description"
	expected := "A"
	todo, _ := Parse(input)

	if todo.priority == nil || strings.Compare(*todo.priority, expected) != 0 {
		t.Fatalf("Priority is incorrect. Expected: \"%s\", but got: \"%s\"\n", expected, *todo.priority)
	}
}

func Test_Parse_Parentheses_In_Description(t *testing.T) {
	input := "x simple (A) description"
	todo, _ := Parse(input)

	if todo.priority != nil {
		t.Fatal("Todo should not have priority.")
	}
}

func Test_Parse_Bad_Priority_Should_Panic(t *testing.T) {
	input := "x (AB) simple description"
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	Parse(input)
}

func Test_Parse_Completion_Date(t *testing.T) {
	input := "2016-05-20 simple description"
	expected := "2016-05-20"
	todo, err := Parse(input)
	fmt.Printf("Todo: %v\n", todo)

	if todo.completionDate == nil || err != nil {
		t.Fatalf("Priority is incorrect. Expected: \"%s\", but got: \"%s\"\n", expected, (*todo.completionDate).Format(YYYYMMDD))
	}
}

func Test_Parse_Bad_Date_Should_Panic(t *testing.T) {
	input := "2016-5-20 simple description"
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	Parse(input)
}
