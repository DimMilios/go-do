package todo

import (
	"reflect"
	"testing"
)

func Test_Find_Todo_By_Description_Text(t *testing.T) {
	prio := "A"
	todo := Todo{Description: Description{Text: "measure space for"}, Original: "x (A) 2016-04-30 measure space for +chapelShelving @chapel due:2016-05-30", Done: true, Priority: &prio}
	todos := []Todo{
		todo,
		{Description: Description{Text: "update screenshots"}},
		{Description: Description{Text: "doctor appointment"}},
	}
	desc := "measure space for"
	if found := FindByDescrText(todos, desc); found == nil || !reflect.DeepEqual(*found, todo) {
		t.Errorf("Couldn't find todo by description text. Expected: %v, but got: %v\n", todo, found)
	}
}

// input := "x (A) 2016-04-30 measure space for +chapelShelving @chapel due:2016-05-30"
// Tags: []Tag{{TagType: Project, Value: "chapelShelving"}, {TagType: Context, Value: "chapel"}, {TagType: KeyValue, Key: &key, Value: "2016-05-30"}}
