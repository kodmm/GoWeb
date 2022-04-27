package trace

import (
	"fmt"
	"io"
)

//Tracerはコード内での出来事を記録できるオブジェクトを表すインタフェース
// type Tracer interface {
// 	Trace(...interface{})
// }

func New(w io.Writer) *Tracer {
	return &Tracer{out: w}
}

type Tracer struct {
	out io.Writer
}

func (t *Tracer) Trace(a ...interface{}) {
	if t == nil || t.out == nil {
		return
	}
	fmt.Fprintln(t.out, a...)
}

// type nilTracer struct{}

// func (t *nilTracer) Trace(a ...interface{}) {}

// OffはTraceメソッドの呼び出しを無視するTracerを返します。
// func Off() Tracer {
// 	return &nilTracer{}
// }
