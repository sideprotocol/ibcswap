package testsuite

import "fmt"

type Logger interface {
	CleanLog(title string, data ...interface{})
}
type logger struct{}

func NewLogger() Logger {
	return &logger{}
}

func (l *logger) CleanLog(title string, data ...interface{}) {
	fmt.Printf("=======[%s]========\n", title)
	fmt.Println(data...)
	fmt.Printf("=======END=========\n\n")
}
