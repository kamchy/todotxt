package main

import (
	"os"

	c "github.com/gookit/color"
	p "github.com/kamchy/gotodotxt/parse"
)

type TaskPart int

const (
	Prio TaskPart = iota
	Cont
	Data
	Contexts
	Projects
	Tags
)

var TaskPartNames = map[TaskPart]string{
	Prio:     "Part: Priority",
	Cont:     "Part: Contents",
	Data:     "Part: Data",
	Contexts: "Part: Contexts",
	Projects: "Part: Projects",
	Tags:     "Part: Tags",
}

type Styler interface {
	GetStyle(n TaskPart) *c.RGBStyle
}
type StylerData struct {
	Map map[TaskPart]*c.RGBStyle
}

func (d StylerData) GetStyle(n TaskPart) *c.RGBStyle {
	return d.Map[n]
}

var sLight = StylerData{
	Map: map[TaskPart]*c.RGBStyle{
		Prio:     c.NewRGBStyle(c.HslInt(20, 60, 60)),
		Cont:     c.NewRGBStyle(c.HslInt(100, 60, 60)),
		Data:     c.NewRGBStyle(c.HslInt(50, 60, 60)),
		Contexts: c.NewRGBStyle(c.HslInt(250, 60, 60)),
		Projects: c.NewRGBStyle(c.HslInt(300, 60, 60)),
		Tags:     c.NewRGBStyle(c.HslInt(350, 60, 60)),
	},
}

func Render(t p.Task, s Styler) {
	if t.Status.IsDone {
		c.Grayln(t.Line)
		return
	}
	if t.Priority != nil {
		s.GetStyle(Prio).Printf("%v ", t.Priority)
	}
	s.GetStyle(Cont).Print(t.Data)

	for _, l := range t.Contexts {
		s.GetStyle(Contexts).Printf(" +%v", l)
	}
	for _, p := range t.Projects {
		s.GetStyle(Projects).Printf(" @%v", p)
	}
	for k, tl := range t.Tags {
		s.GetStyle(Tags).Printf(" %v: %v", k, tl)
	}
	c.Println()
	c.Reset()
}

func main() {
	args := os.Args
	if len(args) < 1 {
		c.Blue.Println("Expected name of file")
		os.Exit(0)
	}

	tasks, errors := ReadFromFile(args[1])
	for _, task := range tasks {
		Render(task, sLight)
	}
	for i, err := range errors {
		c.Errorf("Line %d: %v\n", i, err)
	}
}
