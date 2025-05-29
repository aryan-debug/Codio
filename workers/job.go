package workers

type Language int

const (
	Python Language = iota
	Java
	CPP
)

var LanguageMap = map[string]Language{
	"python": Python,
	"java":   Java,
	"c++":    CPP,
}

var LanguageImageMap = map[Language]string{
	Python: "python_runner",
	Java:   "java_runner",
}

type JobResult struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

type Job struct {
	Language Language
	Code     string
}
