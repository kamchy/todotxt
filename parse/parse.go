// parse is the package where parsing of todo lines takes place
package parse

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// DATE_SIMPLE defines format of the date used for parsing of task date input
const DATE_SIMPLE = "2006-01-02"

// Prio is the priority of the task
type Prio struct {
	P byte
}

func (p *Prio) String() string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("(%s)", string(p.P))
}

// IsByteValidPrio checks if given byte is valid as
// task priority
func IsByteValidPrio(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

// NewPrio creates Prio with given value or nil
// if the value is out of A-Z range
func NewPrio(c byte) *Prio {
	if IsByteValidPrio(c) {
		return &Prio{c}
	}
	return nil
}

// Task is the task representation
type Task struct {
	Data     string
	Status   DoneStatus
	Priority *Prio
	Contexts []string
	Projects []string
	Tags     map[string][]string
}

// TagNames returns list of all tag names in Task's Tags field
func (t Task) TagNames() []string {
	res := make([]string, 0)
	for k := range t.Tags {
		res = append(res, k)
	}
	return res
}

type DoneStatus struct {
	IsDone    bool
	EndDate   *time.Time
	StartDate *time.Time
}

type ParseError struct {
	Message string
}
func (p *ParseError) Error() string {
	return fmt.Sprintf("Error parsing task: %s", p.Message)
}
func ParseDone(fs []string) ([]string, DoneStatus, error) {
	if len(fs) == 0 {
		return fs, DoneStatus{}, &ParseError{"Task is empty line"}
	} 

	if fs[0] != "x" {
		return fs, DoneStatus{IsDone: false}, nil
	}

	if len(fs) <= 1 {
		return fs, DoneStatus{IsDone: true}, nil
	}

	endDate, err := time.ParseInLocation(DATE_SIMPLE, fs[1], time.Local)
	if err != nil {
		return fs[1:], DoneStatus{true, nil, nil}, nil
	}

	if len(fs) < 3 {
		return fs, DoneStatus{}, &ParseError{"Starting date missing"}
	}

	startDate, err := time.ParseInLocation(DATE_SIMPLE, fs[2], time.Local)
	if err != nil {
		return nil, DoneStatus{}, err

	}
	if startDate.After(endDate) {
		return nil, DoneStatus{}, 
		&ParseError{fmt.Sprintf("end Date %v is before Start Date %v", endDate, startDate)}

	}
	return fs[3:], DoneStatus{true, &endDate, &startDate}, nil
}

// Parse parses string s and returns Task and error
func Parse(s string) (Task, error) {
	s = strings.TrimSpace(s)
	fs := strings.Fields(s)
	fs, doneStatus, err := ParseDone(fs)
	if err != nil {
		return Task{}, err
	}

	fs, prio := ParsePriority(fs)
	contents := ParseContents(fs)

	return Task{
		Status:   doneStatus,
		Priority: prio,
		Data:     contents.Data,
		Projects: contents.Projects,
		Contexts: contents.Contexts,
		Tags:     contents.Tags,
	}, nil
}

type Contents struct {
	Data     string
	Contexts []string
	Projects []string
	Tags     map[string][]string
}

func ParseContents(fs []string) Contents {
	var result Contents = Contents{
		Data:     "",
		Contexts: make([]string, 0),
		Projects: make([]string, 0),
		Tags:     make(map[string][]string),
	}
	for _, s := range fs {
		if strings.HasPrefix(s, "+") {
			result.Contexts = append(result.Contexts, s[1:])
		} else if strings.HasPrefix(s, "@") {
			result.Projects = append(result.Projects, s[1:])
		} else if colonIdx := strings.Index(s, ":"); colonIdx > -1 {
			key, value := s[:colonIdx], s[colonIdx+1:]
			vals, found := result.Tags[key]
			if !found {
				result.Tags[key] = make([]string, 0)
			}
			result.Tags[key] = append(vals, value)
		} else {
			result.Data += " " + s
		}
	}
	result.Data = strings.TrimSpace(result.Data)
	return result
}

// ParsePriority parses optional priority. Returns
// rest of fields and pointer to Prio, possibly null
func ParsePriority(fs []string) ([]string, *Prio) {
	re := regexp.MustCompile(`^\([A-Z]\)$`)
	match := re.FindStringIndex(fs[0])
	// ther was no match or not at the begining of first field
	if match == nil {
		return fs, nil
	}
	b := fs[0][1]
	return fs[1:], NewPrio(byte(b))
}
