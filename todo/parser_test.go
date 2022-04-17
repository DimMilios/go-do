package todo

import (
	"fmt"
	"strings"
	"testing"
)

func Test_Parse_Simple_Description(t *testing.T) {
	expected := "simple description"
	todo := Parse(expected)
	got := todo.description.text

	if strings.Compare(got, expected) != 0 {
		t.Fatalf("Description text is incorrect. Expected: \"%s\", but got: \"%s\"\n", expected, got)
	}
}

func Test_Parse_Should_Ignore_Done(t *testing.T) {
	simple := "x simple description"
	expected := "simple description"
	todo := Parse(simple)
	got := todo.description.text

	if strings.Compare(got, expected) != 0 {
		t.Fatalf("Description text is incorrect. Expected: \"%s\", but got: \"%s\"\n", expected, got)
	}
}

func Test_Parse_Mark_As_Done(t *testing.T) {
	input := "x walk dog"
	todo := Parse(input)

	got := string(*todo.done)

	if todo.done == nil || got != DONE_CHAR {
		t.Fatalf("Completion character at index 0 failed to mark todo as done. Got: %s\n", got)
	}
}

func Test_Parse_Mark_As_Done_Correct_Index(t *testing.T) {
	input := "walk x dog"
	todo := Parse(input)

	if todo.done != nil {
		t.Fatal("Should not consider completion char at wrong index as valid.")
	}
}

func Test_Parse_Mark_As_Done_Is_Followed_By_Space(t *testing.T) {
	input := "xx walk dog"
	todo := Parse(input)

	if todo.done != nil {
		t.Fatal("Completion char should be followed by a space.")
	}
}

func Test_Parse_Description_With_Project_Tag(t *testing.T) {
	src := "call customer +proj1"
	expected := "proj1"
	todo := Parse(src)

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
	src := "call customer +proj1 +proj2"
	todo := Parse(src)

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
	src := "call customer @ctx1"
	expected := "ctx1"
	todo := Parse(src)

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
	src := "call customer @ctx1 @ctx2"
	todo := Parse(src)

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
	src := "call customer +proj @ctx1"
	todo := Parse(src)

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
	src := "x call customer +proj @ctx1"
	todo := Parse(src)

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
