package parse

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
)

type TaskExp struct {
	P *Prio
	S string
}

var TestPrioStrData = []TaskExp{
	{NewPrio('A'), "(A)"},
	{NewPrio('Z'), "(Z)"},
	{NewPrio('x'), "nil"},
}

func TestStr(t *testing.T) {
	for i, tsd := range TestPrioStrData {
		td, exp := tsd.P, tsd.S
		act := td.String()
		if exp != act {
			t.Errorf("%d, To string failed: exp %s, act %s", i, exp, act)
		}
	}
}

type ByteValid struct {
	B byte
	V bool
}

var bytesValidData = []ByteValid{
	{'A', true},
	{'B', true},
	{'C', true},
	{'Z', true},

	{'a', false},
	{'b', false},
	{'c', false},
	{'x', false},
	{'1', false},
}

func TestIsByteValid(t *testing.T) {

	for i, d := range bytesValidData {
		exp, act := d.V, IsByteValidPrio(d.B)
		if exp != act {
			t.Errorf("%d) expected byte %b to be %v, got %v", i, d.B, exp, act)
		}
	}
}

type TestParsePrio struct {
	Line     string
	Priority *Prio
}

var testParsePrioData = []TestParsePrio{
	{"(A) foo", NewPrio('A')},
	{"(B) baroo", NewPrio('B')},
	{"(Baca) baroo", nil},
	{"(z) baroo", nil},
	{"(D)ata", nil},
}


func comparePointersPrio(exp *Prio, act *Prio) bool {
	return (exp == nil && act != nil) || 
			(exp != nil && act == nil)|| 
			(exp != nil && act != nil && *exp  != *act)
}

func TestParse(t *testing.T) {
	for _, tpd := range testParsePrioData {
		task, err := Parse(tpd.Line)
		if err != nil {
			t.Fail()
		}
		exp, act := tpd.Priority, task.Priority

		if comparePointersPrio(exp, act) {
				t.Errorf("Prio is different: exp=%v, act=%v", exp, act);
		}
		
	}
}

type StrValidType struct {
	S string
	V bool
}
var data = []StrValidType {
	{"asd(A)", false},
	{"(B)asdasd", false},
	{"(C)", true},
	{"(a)", false},
	{"C", false},
}
func TestRe(t *testing.T) {
	t.Logf("testing %s", t.Name())
	re := regexp.MustCompile(`^\([A-Z]\)$`)
	for _, d := range data {
		match := re.FindStringIndex(d.S)
		if match == nil && d.V {
			t.Errorf("match is nil for %v: %v", d, match)
		}

		if match != nil && match[0] != 0  && d.V {
			t.Errorf("match start is not 0 for %v: %v", d, match)
		}

		if match != nil && match[0] == 0 && !d.V {
			t.Errorf("match at 0 for %v but should not be: %v", d, match)
		}
		

	}

}
type TestParseDoneStatusType struct {
	Line           string
	ExpectedStatus DoneStatus
	StatusCorrect bool
}

func MkTime(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Now().Location())
}

func Status(isDone bool, start time.Time, end time.Time) DoneStatus {
	return DoneStatus{isDone, &start, &end}
}

var testParseDoneStatusData = []TestParseDoneStatusType{
	{"x foo", DoneStatus{true, nil, nil}, true},
	{"x 2022-01-02 2022-01-01 foo", Status(true, MkTime(2022, 1, 2), MkTime(2022, 1, 1)), true},
	{"foo", DoneStatus{}, true},
	{"xfoo", DoneStatus{}, true},
	{"x 2022-01-02foo", DoneStatus{true, nil, nil}, true},
	{"x 2022-61-02 2022-61-03", DoneStatus{true, nil, nil}, true},

	{"x 2022-01-02 foo", DoneStatus{}, false},
	{"x 2022-01-02 2022-61-03", DoneStatus{}, false},
	{"x 2022-12-02 2023-12-03", DoneStatus{}, false},
}

func (ds *DoneStatus) equals(other *DoneStatus) bool {
	return (ds == nil && other == nil) ||
		(ds != nil && other != nil && 
			ds.IsDone == other.IsDone &&
			((ds.StartDate == nil && other.StartDate == nil) ||
			 (ds.StartDate != nil && other.StartDate != nil && ds.StartDate.Equal(*other.StartDate)) &&
			 ((ds.EndDate == nil && other.EndDate == nil) ||
			 (ds.EndDate != nil && other.EndDate != nil && ds.EndDate.Equal(*other.EndDate)))))
}
func TestTestParseDoneStatus(t *testing.T) {
	for i, td := range testParseDoneStatusData {
		// t.Logf("%d %s: Testing %v", i, t.Name(), td)
		exp := td.ExpectedStatus
		task, err := Parse(td.Line)
		if (err != nil && td.StatusCorrect) || ( err == nil && !td.StatusCorrect){
				t.Fail()
		}
		
		act := task.Status
		if !(&act).equals(&exp) {
			t.Errorf("%d) Expected status\n%v, got:\n%v", i, exp, act)
		}
	}
}

type ContentsDataType struct {
	Text string
	ExpText string
	ExpContexts string
	ExpProjects string
	ExpTags string
}
var contentsData = []ContentsDataType {
	{"foo +home", "foo", "home", "", ""},
	{"foo +home bar", "foo bar", "home", "", ""},
	{"foo bar +home", "foo bar", "home", "", ""},
	{"+home", "", "home", "", ""},
	{"foo +home bar +work", "foo bar", "home work", "", ""},
	{"foo +home @life", "foo", "home", "life", ""},
	{"foo +home @life +work @money", "foo", "home work", "life money", ""},
	{"foo +home @life feel:good", "foo", "home", "life", "feel good"},
	{"foo feel:good feel:sad time:high", "foo", "", "", "feel good sad time high"},
}

func assertEquals(t *testing.T, act string, exp string, desc string) {
	if act != exp {
		t.Errorf("%v: exp %v, act %v", desc, exp, act)
	}
}
func flatten(m map[string][]string) string {
	var res string
	var keys = make([]string, len(m))

	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {	
		sort.Strings(m[k])
		res += fmt.Sprintf(" %s %s", k, strings.TrimSpace(strings.Join(m[k], " ")))
	}
	return strings.TrimSpace(res)
}

func TestParseContents(t *testing.T) {
	for _, d := range contentsData {
		act := ParseContents(strings.Fields(d.Text))
		assertEquals(t, act.Data, d.ExpText, "task content")
		assertEquals(t, strings.Join(act.Contexts, " "), d.ExpContexts, "contextx")
		assertEquals(t, strings.Join(act.Projects, " "), d.ExpProjects, "projects")
		assertEquals(t, flatten(act.Tags), d.ExpTags, "projects")
	}
}