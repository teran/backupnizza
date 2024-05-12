package config

type Environment struct {
	Name          string
	Value         string
	ValueFrom     string
	Prefetch      bool
	SourceOptions any
}
