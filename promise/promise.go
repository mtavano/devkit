package promise

import (
	"sync"
)

// Func represents a function that will be executed into a go routine.
type Func func() (interface{}, error)

// ErrorHandler is the custom error handler for manage error on our Promise funcs.
type ErrorHandler func(err error) error

// Group is the structure that holds the Promises and handles the syncronization between them
type Group struct {
	tasks []Func
	wg    *sync.WaitGroup
}

func NewGroup(tasks []Func) *Group {
	wg := new(sync.WaitGroup)
	wg.Add(len(tasks))

	return &Group{
		tasks: tasks,
		wg:    wg,
	}
}

// ExecAll will execute all promise Func no matter if any fails or not. This will return a slice
// of interfaces that must be casted in order to properly handle the desired returned data by
// the defined promise Func
func (g *Group) ExecAll() ([]interface{}, []error, bool) {
	results := make([]interface{}, len(g.tasks), len(g.tasks))
	errs := make([]error, len(g.tasks), len(g.tasks))
	hasError := false

	for idx, fn := range g.tasks {
		go func(idx int, fn Func) {
			defer g.wg.Done()

			res, err := fn()
			if err != nil {
				hasError = true
				errs[idx] = err
				return
			}

			results[idx] = res
		}(idx, fn)
	}

	g.wg.Wait()
	return results, errs, hasError
}

// All will run the promise Func into the slice but if any fails, will skip the subsequent funcs// and will return the error returned by the promise Fun
func (g *Group) All() (_ []interface{}, err error) {
	results := make([]interface{}, 0, len(g.tasks))

	for idx, fn := range g.tasks {
		go func(idx int, fn Func) {
			defer g.wg.Done()

			if err != nil {
				return
			}

			res, e := fn()
			if e != nil {
				err = e
				return
			}

			results = append(results, res)
		}(idx, fn)
	}
	g.wg.Wait()

	if err != nil {
		return nil, err
	}

	return results, nil
}
