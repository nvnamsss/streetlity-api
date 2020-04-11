package pipeline

//Stage representation a step in the Pipeline, a Stage will lead to another Stage.
//If there is no Stage, this Stage will be end of Pipeline
type Stage struct {
	Stages *Stage
	Func   func() error
}

func (r *Stage) Next(l *Stage) {
	r.Stages = l
}

func (s *Stage) Run() error {
	return s.Func()
}

func NewStage(process func() error) *Stage {
	var stage *Stage = new(Stage)
	stage.Func = process
	return stage
}
