package internal

import (
	"log"
	"strings"
	"sync"

	"github.com/deanishe/awgo"
)

var Default *Workflow
var once sync.Once

type Workflow struct {
	*aw.Workflow
}

func (r *Workflow) Run(handler func() error) {
	r.Workflow.Run(func() {
		if err := handler(); err != nil {
			r.FatalError(err)
			return
		}
	})
}

func (r *Workflow) Input() string {
	var query string
	if args := r.Workflow.Args(); len(args) > 0 {
		query = args[0]
	}
	log.Printf("input: %s\n", query)

	return strings.TrimSpace(query)
}

func (r *Workflow) InputWithPre() (string, string) {
	query := r.Input()
	queries := strings.Split(query, "=")
	if len(queries) >= 2 {
		return strings.TrimSpace(queries[0]), strings.TrimSpace(queries[1])
	}

	return query, query
}

func Init() {
	once.Do(func() {
		// wf
		Default = &Workflow{
			Workflow: aw.New(),
		}

		// log
		log.SetPrefix("[alfred_wf_yapi]")
	})
}
