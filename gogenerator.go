// Package gogenerator is an iterator to return elements via callbacks, effectively can be used for example processing large files, databases
package gogenerator

// Generator interface defines the iterator
type Generator interface {
	SetInitFunc(InitFunc)
	Next() chan interface{}
}

// InitFunc is the type of the function which will be called when iteration starts (once)
type InitFunc = func(...interface{}) []interface{}

// CallbackFunc is the callback function format what will be called for each iteration
type CallbackFunc = func(int, ...interface{}) interface{}

// New creates new generator
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
}

// SetInitFunc can be called with a callback func, which will be called once, and the result of the function will be passed to the callback as params
func (g *IterateGenerator) SetInitFunc(fn InitFunc) {
	g.initFunc = fn
}

// Next is to be used in the for loop for res := range generator.Next()
func (g *IterateGenerator) Next() chan interface{} {
	ch := make(chan interface{}, 1)

	go func() {
		defer close(ch)
		if g.initFunc != nil {
			g.initFuncResult = g.initFunc(g.params...)
		}

		i := 0
		for {
			res := g.callbackFunc(i, g.initFuncResult...)
			if res == nil {
				break
			}

			ch <- res
			i++
		}
	}()

	return ch
}
