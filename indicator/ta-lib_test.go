package indicator

import (
	"fmt"
	"testing"
)

func TestEval(t *testing.T) {
	ma := NewMA(3)
	ma.queue.Enqueue(1.0)
	ma.queue.Enqueue(2.0)

	if ma.Eval() != 3.0 {
		fmt.Println("ma.Eval() :  ", ma.Eval())
		t.Error("Expected 3.0, got ", ma.Eval())
	}
}

func TestEvalCB(t *testing.T) {
	ma := NewMA(3)
	ma.queue.Enqueue(1.0)
	ma.queue.Enqueue(2.0)
	ma.queue.Enqueue(3.0)
	ma.queue.Enqueue(nil)
	ma.queue.Enqueue(5.0)
	ma.queue.Enqueue(6.0)
	if ma.Eval() != 5.0 {
		fmt.Println("ma.Eval() :  ", ma.Eval())
		t.Error("Expected 5.0, got ", ma.Eval())
	}
}
