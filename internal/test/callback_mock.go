// Package generatortest is an external functional test
package generatortest

func newCallbackMock() *callbackMock {
	return &callbackMock{}
}

type callbackMock struct {
	params      []interface{}
	initCalled  int
	deferCalled int
	called      int
	initErr     error
	callbackErr error
	deferErr    error
}

func (c *callbackMock) callbackFunc(i int, params ...interface{}) (interface{}, error) {
	c.params = params
	c.called++
	if i >= 5 {
		return nil, nil
	}

	return true, c.callbackErr
}

func (c *callbackMock) withInitError(err error) *callbackMock {
	c.initErr = err
	return c
}

func (c *callbackMock) withCallbackError(err error) *callbackMock {
	c.callbackErr = err
	return c
}

func (c *callbackMock) initFunc(_ ...interface{}) ([]interface{}, error) {
	c.initCalled++
	return []interface{}{1, 2, 3}, c.initErr
}

func (c *callbackMock) deferFunc(_ ...interface{}) error {
	c.deferCalled++
	return c.deferErr
}
