package indicator

import (
	"fmt"
	"testing"
)

func TestEval(t *testing.T) {
	ma := NewMA(3)
	ma.DQ.Enqueue(1.0)
	ma.DQ.Enqueue(2.0)

	if ma.Eval() != 1.5 {
		fmt.Println("ma.Eval() :  ", ma.Eval())
		t.Error("Expected 1.5, got ", ma.Eval())
	}
}

func TestEvalCB(t *testing.T) {
	ma := NewMA(3)
	ma.DQ.Enqueue(1.0)
	ma.DQ.Enqueue(2.0)
	ma.DQ.Enqueue(3.0)
	ma.DQ.Enqueue(nil)
	ma.DQ.Enqueue(5.0)
	ma.DQ.Enqueue(6.0)
	if ma.Eval() != 5.0 {
		fmt.Println("ma.Eval() :  ", ma.Eval())
		t.Error("Expected 5.0, got ", ma.Eval())
	}
}
