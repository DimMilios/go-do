package todo

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Test_Parse_Simple_Description(t *testing.T) {
	expected := "simple description"
	todo, _ := Parse(expected)
	got := todo.Description.Text

	if strings.Compare(got, expected) != 0 {
		t.Errorf("Description text is incorrect. Expected: \"%s\", but got: \"%s\"\n", expected, got)
	}
}

func Test_Parse_Simple_Description_In_Greek(t *testing.T) {
	expected := "απλή περιγραφή"
	todo, _ := Parse(expected)
	got := todo.Description.Text

	if strings.Compare(got, expected) != 0 {
		t.Errorf("Description text is incorrect. Expected: \"%s\", but got: \"%s\"\n", expected, got)
	}
}

func Test_Parse_Should_Ignore_Done(t *testing.T) {
	simple := "x simple description"
	expected := "simple description"
	todo, _ := Parse(simple)
	got := todo.Description.Text

	if strings.Compare(got, expected) != 0 {
		t.Errorf("Description text is incorrect. Expected: \"%s\", but got: \"%s\"\n", expected, got)
	}
}

func Test_Parse_Mark_As_Done(t *testing.T) {
	input := "x walk dog"
	todo, _ := Parse(input)

	if !todo.Done {
		t.Error("Completion character at index 0 failed to mark todo as done.")
	}
}

func Test_Parse_Mark_As_Done_Correct_Index(t *testing.T) {
	input := "walk x dog"
	todo, _ := Parse(input)

	if todo.Done {
		t.Error("Should not consider completion char at wrong index as valid.")
	}
}

func Test_Parse_Mark_As_Done_Is_Followed_By_Space(t *testing.T) {
	input := "xx walk dog"
	todo, _ := Parse(input)

	if todo.Done {
		t.Error("Completion char should be followed by a space.")
	}
}

func Test_Parse_Description_With_Project_Tag(t *testing.T) {
	input := "call customer +proj1"
	expected := "proj1"
	todo, _ := Parse(input)

	var projTag Tag
	for _, tag := range todo.Description.Tags {
		if tag.TagType == Project {
			projTag = tag
			break
		}
	}

	if projTag.Value != expected {
		fmt.Printf("Project tag: %v\n", projTag)
		t.Errorf("Couldn't parse project tag. Expected: \"%s\", but got: \"%s\"\n", expected, projTag.Value)
	}
}

func Test_Parse_Description_With_Two_Project_Tags(t *testing.T) {
	input := "call customer +proj1 +proj2"
	todo, _ := Parse(input)

	var projTags []Tag
	for _, tag := range todo.Description.Tags {
		if tag.TagType == Project {
			projTags = append(projTags, tag)
		}
	}

	if len(projTags) != 2 {
		fmt.Printf("Todo: %v\n", todo)
		t.Errorf("Couldn't parse many project tags. Tags: %v", projTags)
	}
}

func Test_Parse_Description_With_Context_Tag(t *testing.T) {
	input := "call customer @ctx1"
	expected := "ctx1"
	todo, _ := Parse(input)

	var ctxTag Tag
	for _, tag := range todo.Description.Tags {
		if tag.TagType == Context {
			ctxTag = tag
		}
	}

	if ctxTag.Value != expected {
		t.Errorf("Couldn't parse context tag. Expected: \"%s\", but got: \"%s\"\n", expected, ctxTag.Value)
	}
}

func Test_Parse_Description_With_Two_Context_Tags(t *testing.T) {
	input := "call customer @ctx1 @ctx2"
	todo, _ := Parse(input)

	var ctxTags []Tag
	for _, tag := range todo.Description.Tags {
		if tag.TagType == Context {
			ctxTags = append(ctxTags, tag)
		}
	}

	if len(ctxTags) != 2 {
		fmt.Printf("Todo: %v\n", todo)
		t.Errorf("Couldn't parse many context tags. Tags: %v", ctxTags)
	}
}

func Test_Parse_Description_With_Project_And_Context_Tags(t *testing.T) {
	input := "call customer +proj @ctx1"
	todo, _ := Parse(input)

	var tags []Tag
	for _, tag := range todo.Description.Tags {
		if tag.TagType == Project || tag.TagType == Context {
			tags = append(tags, tag)
		}
	}

	if len(tags) != 2 {
		fmt.Printf("Todo: %v\n", todo)
		t.Errorf("Couldn't parse mixed project and context tags. Tags: %v", tags)
	}
}

func Test_Parse_Mark_Done_Complex_Description(t *testing.T) {
	input := "x call customer +proj @ctx1"
	todo, _ := Parse(input)

	var tags []Tag
	for _, tag := range todo.Description.Tags {
		if tag.TagType == Project || tag.TagType == Context {
			tags = append(tags, tag)
		}
	}

	if len(tags) != 2 || !todo.Done {
		t.Errorf("Couldn't mark todo with tagged description as done. Todo: %v", todo)
	}
}

func Test_Parse_Description_With_Key_Value_Tag(t *testing.T) {
	input := "call customer due:now"
	todo, _ := Parse(input)
	expectedKey := "due"
	expectedVal := "now"

	var keyValueTag Tag
	for _, tag := range todo.Description.Tags {
		if tag.Key != nil && *tag.Key == expectedKey && tag.Value == expectedVal {
			keyValueTag = tag
		}
	}

	if keyValueTag.Key == nil || *keyValueTag.Key != expectedKey || keyValueTag.Value != expectedVal {
		t.Errorf("Couldn't parse key value tag. Tag: %v", keyValueTag)
	}
}

func Test_Parse_Key_Value_Tag_Empty_Value(t *testing.T) {
	input := "call customer due:"
	_, err := Parse(input)

	if err == nil {
		t.Error("Trying to pass a key value tag without suppling a value should return an error.")
	}
}

func Test_Parse_Description_With_Many_Kv_Tags(t *testing.T) {
	input := "call customer due:now who:me test:ing"
	todo, _ := Parse(input)

	var kvTags []Tag
	for _, tag := range todo.Description.Tags {
		if tag.TagType == KeyValue {
			kvTags = append(kvTags, tag)
		}
	}

	if len(kvTags) != 3 {
		fmt.Printf("Todo: %v\n", todo)
		t.Errorf("Couldn't parse many key value tags. Tags: %v", kvTags)
	}
}

func Test_Description_Doesnt_Contain_Key_Of_Kv_Tag(t *testing.T) {
	input := "call customer due:now"
	expected := "call customer"
	todo, _ := Parse(input)

	if todo == nil || todo.Description.Text != expected {
		t.Error("Key of key value tag shouldn't be included to the text description.")
	}
}

func Test_Parse_Description_With_All_Tags(t *testing.T) {
	input := "call customer +proj @ctx1 due:now"
	todo, _ := Parse(input)

	if len(todo.Description.Tags) != 3 {
		fmt.Printf("Todo: %v\n", todo)
		t.Errorf("Couldn't parse mixed types of tags. Tags: %v", todo.Description.Tags)
	}
}

func Test_Parse_Description_With_All_Tags_Reordered(t *testing.T) {
	input := "call customer due:now @ctx1 +proj1"
	todo, _ := Parse(input)

	if len(todo.Description.Tags) != 3 {
		fmt.Printf("Todo: %v\n", todo)
		t.Errorf("Couldn't parse mixed types of tags. Tags: %v", todo.Description.Tags)
	}
}

func Test_Parse_Description_With_Multiple_Of_Each_Tag(t *testing.T) {
	input := "call customer due:now @ctx1 who:john +proj1 +proj2 @ctx2"
	todo, _ := Parse(input)

	if len(todo.Description.Tags) != 6 {
		fmt.Printf("Todo: %v\n", todo)
		t.Errorf("Couldn't parse mixed types of tags. Tags: %v", todo.Description.Tags)
	}
}

func Test_Parse_Priority(t *testing.T) {
	input := "x (A) simple description"
	expected := "A"
	expectedDescription := "simple description"
	todo, _ := Parse(input)

	if todo.Priority == nil || strings.Compare(*todo.Priority, expected) != 0 || strings.Compare(todo.Description.Text, expectedDescription) != 0 {
		t.Errorf("Priority is incorrect. Expected: \"%s\", but got: \"%s\"\n", expected, *todo.Priority)
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

func Parse_Completion_Date(t *testing.T) {
	input := "2016-05-20 simple description"
	expected := "2016-05-20"
	todo, err := Parse(input)
	fmt.Printf("Todo: %v\n", todo)

	if todo.CompletionDate == nil || err != nil {
		t.Errorf("Priority is incorrect. Expected: \"%s\", but got: \"%s\"\n", expected, (*todo.CompletionDate).Format(YYYYMMDD))
	}
}

func Test_Parse_Creation_Date(t *testing.T) {
	input := "x (A) 2022-04-20 2022-04-21 update screenshots +proj"
	createDate := time.Now().Format(YYYYMMDD)
	complDate := "2022-04-21"
	todo, err := Parse(input)
	fmt.Printf("Todo: %v\n", todo)

	if strings.Compare(todo.CreationDate.Format(YYYYMMDD), createDate) != 0 {
		t.Errorf("Bad creation date. Expected: \"%s\", but got: \"%s\"\n", createDate, (todo.CreationDate).Format(YYYYMMDD))
	}

	if todo.CompletionDate == nil || err != nil {
		t.Errorf("Bad completion date. Expected: \"%s\", but got: \"%s\"\n", complDate, (*todo.CompletionDate).Format(YYYYMMDD))
	}
}

func Test_Parse_Bad_Date_Should_Panic(t *testing.T) {
	inputs := []string{"2015-5-20 simple description", "20-05-20 simple description", "2015--20 simple description", "2015-20 simple due:now @ctx1 +proj1", "x 2015-20 simple @ctx1 +proj1"}

	for _, input := range inputs {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered:", r)
			}
		}()
		Parse(input)
	}
}

func hasValue(s *string, value string) bool {
	if s == nil || strings.Compare(*s, value) != 0 {
		return false
	}
	return true
}

func getFirstTagOfType(tags []Tag, tagType TagType) *Tag {
	for _, t := range tags {
		if t.TagType == tagType {
			return &t
		}
	}
	return nil
}

func Test_Parse_Complete_Todo_Input(t *testing.T) {
	input := "x (A) 2016-04-30 measure space for +chapelShelving @chapel due:2016-05-30"
	todo, _ := Parse(input)

	if todo.CompletionDate == nil {
		t.Error("Bad completion date.")
	} else {
		s := (*todo.CompletionDate).Format(YYYYMMDD)
		expected := "2016-04-30"
		if !hasValue(&s, expected) {
			t.Errorf("Completion date incorrect. Expected: \"%s\", but got: \"%s\"\n", expected, s)
		}
	}

	expectedDesc := "measure space for"
	if strings.Compare(todo.Description.Text, expectedDesc) != 0 {
		fmt.Println(todo)
		t.Errorf("Description incorrect. Expected: \"%s\", but got: \"%s\"\n", expectedDesc, todo.Description.Text)
	}

	if !todo.Done || !hasValue(todo.Priority, "A") {
		t.Error("Todo is not done or has incorrect priority value.")
	}

	projTag := getFirstTagOfType(todo.Description.Tags, Project)
	if projTag == nil {
		t.Error("Failed to parse project tag.")
	}
	ctxTag := getFirstTagOfType(todo.Description.Tags, Context)
	if ctxTag == nil {
		t.Error("Failed to parse context tag.")
	}
	kvTag := getFirstTagOfType(todo.Description.Tags, KeyValue)
	if kvTag == nil {
		t.Error("Failed to parse key value tag.")
	}
}
