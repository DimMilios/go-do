package todo

import (
	"reflect"
	"strings"
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

func todoLiterals() []string {
	return []string{"(A) 2022-04-20 2022-04-22 update screenshots +project 10",
		"x (B) 2022-04-20 2022-04-22 walk dog +project",
		"x 2011-03-02 2011-03-01 Review Tim's pull request +TodoTxtTouch @github",
		"x 2011-03-03 Call Mom due:now",
		"(A) Call Mom +Family +PeaceLoveAndHappiness @iphone @phone",
		"2011-03-02 Document +TodoTxt task format",
		"(A) Thank Mom for the meatballs @phone",
		"(B) Schedule Goodwill pickup +GarageSale @phone",
		"2018-04-12 2018-04-28 Post signs around the neighborhood +GarageSale ends:tomorrow",
		"@GroceryStore Eskimo pies",
		"x Schedule +dentist @phone",
		"(A) 2022-01-01 2022-01-02 doctor appointment @personal",
	}
}

func Test_Skip_First_Occurrence(t *testing.T) {
	testcases := []struct{ query, line string }{
		{"Mom", "x 2011-03-03 Call Mom due:now"},
		{"Thank Mom for the", "(A) Thank Mom for the meatballs @phone"},
		{"Goodwill", "(B) Schedule Goodwill pickup +GarageSale @phone"},
		{"around the neighborhood", "2018-04-12 2018-04-28 Post signs around the neighborhood +GarageSale ends:tomorrow"},
		{"2018-04-12 2018-04-28 Post signs around the neighborhood +GarageSale ends:tomorrow", "2018-04-12 2018-04-28 Post signs around the neighborhood +GarageSale ends:tomorrow"},
		{"pies", "@GroceryStore Eskimo pies"},
		{"@GroceryStore Eskimo pies", "@GroceryStore Eskimo pies"},
		{"x (B) 2022-04-20 2022-04-22 walk dog +project", "x (B) 2022-04-20 2022-04-22 walk dog +project"},
	}

	t.Run("One by one", func(t *testing.T) {
		lines := strings.Join(todoLiterals(), "\n")

		for _, tc := range testcases {
			query, line := tc.query, tc.line
			newLines, _ := SkipFirst(strings.NewReader(lines), query)
			if pos := FindLineByText(newLines, line); pos >= 0 {
				t.Errorf("Entry %s should be removed, but it was found in the list", line)
			}
		}
	})

	t.Run("Bulk remove, bulk check", func(t *testing.T) {
		lines := strings.Join(todoLiterals(), "\n")
		newLines, _ := SkipFirst(strings.NewReader(lines), testcases[0].query)

		for i := 1; i < len(testcases); i++ {
			newLines, _ = SkipFirst(strings.NewReader(strings.Join(newLines, "\n")), testcases[i].query)
		}

		for _, tc := range testcases {
			if pos := FindLineByText(newLines, tc.line); pos >= 0 {
				t.Errorf("Entry %s should be removed, but it was found in the list", tc.line)
			}
		}
	})
}

func Test_Filter_By_Completion_Status(t *testing.T) {
	t.Run("Filter done", func(t *testing.T) {
		done := []string{
			"x (B) 2022-04-20 2022-04-22 walk dog +project",
			"x 2011-03-02 2011-03-01 Review Tim's pull request +TodoTxtTouch @github",
			"x 2011-03-03 Call Mom due:now",
			"x Schedule +dentist @phone",
		}

		got := GetCompleteLines(todoLiterals())
		if len(got) != len(done) {
			t.Error("Slices are not of same length")
		}
		for i, l := range got {
			if l != done[i] {
				t.Errorf("Slices not equal. Line %q is not equal to line %q\n", l, done[i])
			}
		}
	})

	t.Run("Filter incomplete", func(t *testing.T) {
		todo := []string{"(A) 2022-04-20 2022-04-22 update screenshots +project 10",
			"(A) Call Mom +Family +PeaceLoveAndHappiness @iphone @phone",
			"2011-03-02 Document +TodoTxt task format",
			"(A) Thank Mom for the meatballs @phone",
			"(B) Schedule Goodwill pickup +GarageSale @phone",
			"2018-04-12 2018-04-28 Post signs around the neighborhood +GarageSale ends:tomorrow",
			"@GroceryStore Eskimo pies",
			"(A) 2022-01-01 2022-01-02 doctor appointment @personal",
		}

		got := GetIncompleteLines(todoLiterals())
		if len(got) != len(todo) {
			t.Error("Slices are not of same length")
		}
		for i, l := range got {
			if l != todo[i] {
				t.Errorf("Slices not equal. Line %q is not equal to line %q\n", l, todo[i])
			}
		}
	})

}
