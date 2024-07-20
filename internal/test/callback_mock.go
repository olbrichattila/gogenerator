// Package generatortest is an external functional test
package generatortest

func newCallbackMock() *callbackMock {
	return &callbackMock{}
}

type callbackMock struct {
	params     []interface{}
	initCalled int
	called     int
}

func (c *callbackMock) callbackFunc(i int, params ...interface{}) interface{} {
	c.params = params
	c.called++
	if i >= 5 {
		return nil
	}

	return true
}

func (c *callbackMock) initFunc(_ ...interface{}) []interface{} {
	c.initCalled++
	return []interface{}{1, 2, 3}
}
