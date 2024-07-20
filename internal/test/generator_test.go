package generatortest

import (
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

	t.NotNil(callbackMock.params)
	t.Len(callbackMock.params, 3)

	t.Equal(5, i)
	t.Equal(1, callbackMock.initCalled)
	t.Equal(6, callbackMock.called)
}
