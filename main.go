package main

import "github.com/deanishe/awgo"

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

func run() error {
	wf.NewItem("First result!")
	wf.SendFeedback()
	return nil
}

func main() {
	wf.Run(func() {
		if err := run(); err != nil {
			wf.FatalError(err)
			return
		}
	})
}
