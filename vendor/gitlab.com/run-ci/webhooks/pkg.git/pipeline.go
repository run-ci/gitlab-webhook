package pkg

type Pipeline struct {
	Name   string `json:"name"`
	Remote string `json:"remote"`
	Branch string `json:"branch" yaml:"branch"`
	Tag    string `json:"tag" yaml:"tag"`
	Steps  []Step `json:"steps" yaml:"steps"`
}

type Step struct {
	Name  string `json:"name" yaml:"name"`
	Tasks []Task `json:"task" yaml:"tasks"`
}

type Task struct {
	Name      string                 `json:"name" yaml:"name"`
	Arguments map[string]interface{} `yaml:"arguments"`
}

type PipelineSender interface {
	SendPipeline(Pipeline) error
}
