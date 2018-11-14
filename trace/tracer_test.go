package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		msg := "Hello, trace package."
		tracer.Trace(msg)
		if buf.String() != msg+"\n" {
			t.Errorf("文字列が一致しません : '%s'", buf.String())
		}
	}
}
func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("data")
}
