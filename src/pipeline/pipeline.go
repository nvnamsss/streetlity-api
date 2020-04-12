package pipeline

//Pipeline representation a process of work needs to be run by order and will be stopped if a step is failed.
type Pipeline struct {
	First *Stage
}

//Run start the pipeline, success when the pipeline is passed when all stages are passed.
func (p Pipeline) Run() error {
	var current *Stage = p.First
	for current != nil {
		err := current.Run()

		if err != nil {
			return err
		}

		current = current.Stages
	}

	return nil
}

//Constructor for creating new Pipeline
func NewPipeline() *Pipeline {
	var pipe *Pipeline = new(Pipeline)

	return pipe
}
