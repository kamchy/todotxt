package main

import (
	"fmt"
	c "github.com/gookit/color"
	p "github.com/kamchy/gotodotxt/parse"
)
type TaskStyleName int
const (
	Prio TaskStyleName = iota;
	Cont
	Data
	Contexts
	Projects
	Tags
)
type Styler interface {
	GetStyle(n TaskStyleName) *c.RGBStyle
}
type StylerData struct {
	Map map[TaskStyleName] *c.RGBStyle

} 

func (d StylerData) GetStyle(n TaskStyleName) *c.RGBStyle {
	return d.Map[n]
}

var sLight = StylerData{
	Map: map[TaskStyleName]*c.RGBStyle{
		Prio: c.NewRGBStyle(c.HslInt(20, 60, 60)),
		Cont: c.NewRGBStyle(c.HslInt(100, 60, 60)),
		Data: c.NewRGBStyle(c.HslInt(50, 60, 60)),
		Contexts: c.NewRGBStyle(c.HslInt(250, 60, 60)),
		Projects: c.NewRGBStyle(c.HslInt(300, 60, 60)),
		Tags: c.NewRGBStyle(c.HslInt(350, 60, 60)),
	},
}
func Render(t p.Task, s Styler){
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
	for k, tl := range t.Tags{
		s.GetStyle(Tags).Printf(" %v: %v", k, tl)
	}
}

func main() {
	t, e := p.Parse("(A) nauczyć się go @praca +kuchnia feel:good feel:down")
	if e != nil {
		fmt.Printf("Error: Cannot parse %v", e)

	}
	Render(t, sLight)
}
