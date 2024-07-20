package generatortest

import (
	"fmt"
	"testing"

	generator "github.com/olbrichattila/gogenerator"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) TestNoInitFuncAndIterationWorks() {
	callbackMock := newCallbackMock()

	g := generator.New(callbackMock.callbackFunc)

	i := 0
	for res := range g.Next() {
		t.True(res.(bool))
		i++
	}

	t.Nil(g.GetLastError())
	t.Nil(callbackMock.params)
	t.Equal(5, i)
	t.Equal(0, callbackMock.initCalled)
	t.Equal(6, callbackMock.called)
}

func (t *TestSuite) TestInitFuncAndIterationWorks() {
	callbackMock := newCallbackMock()
	g := generator.New(callbackMock.callbackFunc)

	g.SetInitFunc(callbackMock.initFunc)

	i := 0
	for res := range g.Next() {
		t.True(res.(bool))
		i++
	}

	t.Nil(g.GetLastError())
	t.NotNil(callbackMock.params)
	t.Len(callbackMock.params, 3)

	t.Equal(5, i)
	t.Equal(1, callbackMock.initCalled)
	t.Equal(6, callbackMock.called)
}

func (t *TestSuite) TestInitFuncReturnsError() {
	callbackMock := newCallbackMock().withInitError(fmt.Errorf("init error"))
	g := generator.New(callbackMock.callbackFunc)

	g.SetInitFunc(callbackMock.initFunc)

	i := 0
	for res := range g.Next() {
		t.True(res.(bool))
		i++
	}

	t.Error(g.GetLastError())
	t.Nil(callbackMock.params)
	t.Len(callbackMock.params, 0)

	t.Equal(0, i)
	t.Equal(1, callbackMock.initCalled)
	t.Equal(0, callbackMock.called)
}

func (t *TestSuite) TestCallbackFuncReturnsError() {
	callbackMock := newCallbackMock().withCallbackError(fmt.Errorf("callback error"))
	g := generator.New(callbackMock.callbackFunc)

	g.SetInitFunc(callbackMock.initFunc)

	i := 0
	for res := range g.Next() {
		t.True(res.(bool))
		i++
	}

	t.Error(g.GetLastError())
	t.NotNil(callbackMock.params)
	t.Len(callbackMock.params, 3)

	t.Equal(0, i)
	t.Equal(1, callbackMock.initCalled)
	t.Equal(1, callbackMock.called)
}
