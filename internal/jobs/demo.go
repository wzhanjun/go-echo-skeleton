package jobs

import "fmt"

type demoJobBuilder struct {
}

func NewDemoJobBuilder() *demoJobBuilder {
	return &demoJobBuilder{}
}

func (s *demoJobBuilder) Build() func() {
	return func() {
		fmt.Println("this is a demo job")
	}
}
