package workers

type Language int

const (
	Python Language = iota
	Java
)

var LanguageImageMap = map[Language]string{
	Python: "python_runner",
	Java:   "java_runner",
}

type JobResult struct {
	ID     string
	Output string
	Error  error
}

type Job struct {
	ID       string
	Language Language
	Code     string
}
