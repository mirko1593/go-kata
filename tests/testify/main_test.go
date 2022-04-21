package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSomething(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(123, 123, "should be equla")
}

// MyMockedObject is a mocked object that implements an interface
// that describes an object that the code I am testing relies on.
type MyMockedObject struct {
	mock.Mock
}

// DoSomething is a method on MyMockedObject that implements some interface
// and just records the activity, and returns what the Mock object tells it to.
//
// In the real object, this method would do something useful, but since this
// is a mocked object - we're just going to stub it out.
//
// NOTE: This method is not being tested here, code that uses this object is.
func (m *MyMockedObject) DoSomething(number int) (bool, error) {
	args := m.Called(number)
	return args.Bool(0), args.Error(1)
}

// TestSomethingElse ...
func TestSomethingElse(t *testing.T) {
	testObj := new(MyMockedObject)
	defer testObj.AssertExpectations(t)

	testObj.On("DoSomething", 123).Return(true, nil)

	// testObj.ShouldCall("DoSomething", 123).Return(true, nil)

	// testing.Mock(MyMockedObject).ShouldCall("DoSomething", 123).Return(true, nil)

	testObj.DoSomething(123)
}

// Back ...
type Back struct{}

// SetName ...
func (b *Back) SetName(name string) {
}

// Save ...
func (b *Back) Save() {}

func TestBack(t *testing.T) {

	assert := assert.New(t)
	b := &Back{}
	b.SetName("newname")
	b.Save()

	assert.Equal("", "")
}
