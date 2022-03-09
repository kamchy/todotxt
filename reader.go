package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/gookit/color"
	p "github.com/kamchy/gotodotxt/parse"
)

func ReadTasks(r io.Reader) ([]p.Task, []error) {
	tasks := make([]p.Task, 0)
	errors := make([]error, 0)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		t, e := p.Parse(strings.TrimSpace(line))
		if e != nil {
			errors = append(errors, fmt.Errorf("could not parse task line %s: %v", line, e))
		} else {
			tasks = append(tasks, t)
		}
	}
	return tasks, errors
}
func check(desc string, e error) {
	if e != nil {
		color.Errorf("%s: %v", desc, e)
	}
}

// ReadFromFile reads tasks from file with given name
// and returns array of tasks and corresponding array of errors
func ReadFromFile(s string) ([]p.Task, []error) {
	f, err := os.Open(s)
	defer func(f *os.File) error {
		err := f.Close()
		return err
	}(f)
	check(fmt.Sprintf("Cannot open %s", s), err)
	return ReadTasks(bufio.NewReader(f))
}

type Tasks struct{ Ts []p.Task }

func dedup(s []string) []string {
	if len(s) == 0 {
		return s
	}
	sort.Strings(s)
	var res []string = make([]string, 0)
	curr := s[0]
	if len(curr) > 0 {
		res = append(res, curr)
	}
	for i := 1; i < len(s); i++ {
		if s[i] != curr && s[i] != "" {
			res = append(res, s[i])
			curr = s[i]
		}
	}
	return res
}

type partDataRetriever = func(p.Task) []string

func makePartMapper() map[TaskPart]partDataRetriever {
	var m map[TaskPart]partDataRetriever = make(map[TaskPart]partDataRetriever)
	m[Prio] = func(t p.Task) []string {
		if t.Priority != nil {
			return []string{t.Priority.String()}
		} else {
			return []string{}
		}
	}
	m[Cont] = func(t p.Task) []string { return []string{t.Data} }
	m[Contexts] = func(t p.Task) []string { return t.Contexts }
	m[Projects] = func(t p.Task) []string { return t.Projects }
	m[Tags] = func(t p.Task) []string { return t.TagNames() }
	return m
}

var partMapper = makePartMapper()

// Get returns given task part as list of strings
// in the order the strings appear in task
func (ts Tasks) Get(p TaskPart) []string {
	var cxs = make([]string, len(ts.Ts))
	for _, t := range ts.Ts {
		cxs = append(cxs, partMapper[p](t)...)
	}
	sort.Strings(cxs)
	cxs = dedup(cxs)
	return cxs
}
