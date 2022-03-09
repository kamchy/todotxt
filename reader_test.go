package main

import (
	"sort"
	"strings"
	"testing"

	mt "github.com/kamchy/gotodotxt/test"
)

type Strs = []string

type TaskReaderTestType struct {
	Tasks      string
	ValidCount int
	Prios      Strs
	Contents   Strs
	Contexts   Strs
	Projects   Strs
	Errors     Strs
	TagNames   Strs
}

func (tt TaskReaderTestType) GetBy(tp TaskPart) Strs {
	switch tp {
	case Cont:
		return tt.Contents
	case Prio:
		return tt.Prios
	case Contexts:
		return tt.Contexts
	case Projects:
		return tt.Projects
	case Tags:
		return tt.TagNames
	default:
		return Strs{}

	}
}

var TestReadTasksData = []TaskReaderTestType{
	{`foo bar
blabla
x done 1
x done 2 +home`, 4,
		Strs{},
		Strs{"foo bar", "blabla", "done 1", "done 2"},
		Strs{"home"},
		Strs{},
		Strs{},
		Strs{}},
	{`(A) foo bar +work
(B) blabla @life @kids
x done 1
x done 2 +home`, 4,
		Strs{"(A)", "(B)"},
		Strs{"foo bar", "blabla", "done 1", "done 2"},
		Strs{"home", "work"},
		Strs{"life", "kids"},
		Strs{},
		Strs{}},
	{`foo bar
(C) blabla
x 2022-02-12 lksdjf done 1
x done 2 +home foo:bar bing:bang`, 3,
		Strs{"(C)"},
		Strs{"foo bar", "blabla", "done 2"},
		Strs{"home"},
		Strs{},
		Strs{"could not parse task line"},
		Strs{"bing", "foo"}},
	{`
blabla
x done 1
x done 2 +home repeat:daily`, 3,
		Strs{},
		Strs{"blabla", "done 1", "done 2"}, Strs{"home"},
		Strs{},
		Strs{"Task is empty line"},
		Strs{"repeat"}},
}

func TestReadTasks(t *testing.T) {
	for _, td := range TestReadTasksData {
		tasks, errors := ReadTasks(strings.NewReader(td.Tasks))
		exp := td.ValidCount
		act := len(tasks)
		if exp != act {
			t.Errorf("Expected %d, got %d tasks", exp, act)
		}

		if el, expel := len(errors), len(td.Errors); el != expel {
			t.Errorf("Expeced %d errors, got %d", expel, el)
		} else {
			for i := 0; i < el; i++ {
				mt.AssertStartsWith(t, "Error should start with expected:", td.Errors[i], errors[i].Error())
			}
		}
	}
}
func TestParts(t *testing.T) {
	for _, part := range []TaskPart{Contexts, Projects, Tags, Prio, Cont} {
		testGetByPart(t, part)
	}
}
func testGetByPart(t *testing.T, part TaskPart) {
	for _, td := range TestReadTasksData {
		tasks, _ := ReadTasks(strings.NewReader(td.Tasks))
		exp := td.GetBy(part)
		sort.Strings(exp)
		act := Tasks{tasks}.Get(part)
		sort.Strings(act)
		if len(exp) != len(act) {
			t.Logf(`%v`, TaskPartNames[part])
			t.Errorf(
				`Expected num of %s: %v (%v), got %v (%v)`, td.Tasks, len(exp), exp, len(act), act)
		}
		for i, expAt := range exp {
			if expAt != act[i] {
				t.Logf(`%v`, TaskPartNames[part])
				t.Errorf("Expected %v: %v, got: %v", part, expAt, act[i])
			}
		}
	}
}
