package pipeline

import (
	"fmt"
	"reflect"
)

//Stage representation a step in the Pipeline, a Stage will lead to another Stage.
//If there is no Stage, this Stage will be the end of Pipeline
type Stage struct {
	Stages  *Stage
	task    func() error
	taskalt reflect.Value
}

//Next set the left Stage as the next stage of right Stage
//In the Pipeline, left Stage will run when the right is done
func (r *Stage) Next(l *Stage) {
	r.Stages = l
}

func (r *Stage) NextStage(task interface{}) {
	t := reflect.ValueOf(task).Type()
	olen := t.NumOut()
	if olen == 1 {
		var fn func() (struct{}, error)

		r.taskalt = reflect.MakeFunc(reflect.ValueOf(fn).Type(), func(args []reflect.Value) (results []reflect.Value) {
			result := reflect.ValueOf(task).Call(nil)
			result = append([]reflect.Value{reflect.ValueOf(struct{}{})}, result...)

			return result
		})

	}

	if olen == 2 {
		if t.Out(0).Kind() != reflect.Struct {
			panic("The first return type must be struct")
		}

		r.taskalt = reflect.ValueOf(task)
	}
}

func (s *Stage) Run() error {
	return s.task()
}

func NewStage(process func() error) *Stage {
	var stage *Stage = new(Stage)
	stage.task = process
	return stage
}

func init() {
	stage := NewStage(nil)
	fn := func() (struct {
		Field string
		Meo   int
	}, error) {

		return struct {
			Field string
			Meo   int
		}{"meomeocute", 1}, nil
	}

	var p *Pipeline = NewPipeline()
	stage.NextStage(fn)
	p.First = stage
	p.Run()

	fmt.Println(p.GetString("Field"))
	fmt.Println(p.GetFloat("Meo"))
}
