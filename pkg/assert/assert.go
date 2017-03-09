package assert

import (
	"reflect"
	"runtime"
	"testing"
)

func assert(t *testing.T, result bool, f func(), cd int) {
	if !result {
		_, file, line, _ := runtime.Caller(cd + 1)
		t.Errorf("%s:%d", file, line)
		f()
		t.FailNow()
	}
}

func Equal(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Did not get expected value at %s:%d expected %v, actual %v", file, line, expected, actual)
	}
}

func NotEqual(t *testing.T, expected interface{}, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Got unexpected value at %s:%d expected %v, actual %v", file, line, expected, actual)
	}
}
