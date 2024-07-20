// Package gogenerator provides an iterator to return elements via callbacks, suitable for processing large files, databases, etc.
package gogenerator

// Generator interface defines the iterator
type Generator interface {
	SetInitFunc(InitFunc)
	Next() <-chan interface{}
	GetLastError() error
}

// InitFunc is the type of the function which will be called when iteration starts (once)
type InitFunc func(...interface{}) ([]interface{}, error)

// CallbackFunc is the callback function format which will be called for each iteration
type CallbackFunc func(int, ...interface{}) (interface{}, error)

// New creates a new generator
func New(f CallbackFunc, params ...interface{}) Generator {
	return &IterateGenerator{
		callbackFunc: f,
		params:       params,
	}
}

// IterateGenerator is the iterator struct
type IterateGenerator struct {
	initFunc       InitFunc
	initFuncResult []interface{}
	params         []interface{}
	callbackFunc   CallbackFunc
	err            error
}

// SetInitFunc can be called with a callback function, which will be called once, and the result will be passed to the callback as params
func (g *IterateGenerator) SetInitFunc(fn InitFunc) {
	g.initFunc = fn
}

// Next is to be used in the for loop for res := range generator.Next()
func (g *IterateGenerator) Next() <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		defer close(ch)

		if g.initFunc != nil {
			var err error
			g.initFuncResult, err = g.initFunc(g.params...)
			if err != nil {
				g.err = err
				return
			}
		}

		i := 0
		for {
			res, err := g.callbackFunc(i, g.initFuncResult...)
			if err != nil || res == nil {
				g.err = err
				break
			}

			ch <- res
			i++
		}
	}()

	return ch
}

// Error returns any error encountered during initialization or iteration
func (g *IterateGenerator) GetLastError() error {
	return g.err
}
