package workers

type LanguageInfo struct {
	ImageName string
	FileName  string
}

var LanguageMap = map[string]LanguageInfo{
	"python": {
		ImageName: "python_runner",
		FileName:  "test.py",
	},
	"cpp": {
		ImageName: "cpp_runner",
		FileName:  "test.cpp",
	},
	"go": {
		ImageName: "go_runner",
		FileName:  "test.go",
	},
}

type JobResult struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

type Job struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}
