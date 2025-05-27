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
	Output string
	Error  error
}

type Job struct {
	Language Language
	Code     string
}
