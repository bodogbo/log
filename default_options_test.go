package log

import (
	"reflect"
	"testing"
)

func Test_GetDefaultOptions_SetDefaultOptions(t *testing.T) {
	// no call SetDefaultOptions
	opts := getDefaultOptions()
	if opts != nil {
		t.Errorf("have:%v, want:nil", opts)
		return
	}

	// call SetDefaultOptions with non-nil Options
	opt := WithFormatter(JsonFormatter)
	opts = []Option{
		opt,
	}
	SetDefaultOptions(opts)
	opts2 := getDefaultOptions()
	if !reflect.DeepEqual(opts, opts2) {
		t.Errorf("have:%v, want:%v", opts2, opts)
		return
	}

	// call SetDefaultOptions with nil Options
	SetDefaultOptions(nil)
	opts = getDefaultOptions()
	if opts != nil {
		t.Errorf("have:%v, want:nil", opts)
		return
	}

	// call SetDefaultOptions with non-nil Options
	opt = WithFormatter(JsonFormatter)
	opts = []Option{
		opt,
	}
	SetDefaultOptions(opts)
	opts2 = getDefaultOptions()
	if !reflect.DeepEqual(opts, opts2) {
		t.Errorf("have:%v, want:%v", opts2, opts)
		return
	}
}