// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Wonderstone (Digital Office Product Department #2)
// revisor:

package indicator

import (
	"fmt"
	"testing"
)

func TestEvalMA(t *testing.T) {
	ma := NewMA("MA3", []int{3}, []string{"Close"})

	ma.LoadData(map[string]float64{"Close": 1.0})
	fmt.Println(ma.EvalOld(), ma.Eval(), ma.DQ.Full())
	ma.LoadData(map[string]float64{"Close": 2.0})
	fmt.Println(ma.EvalOld(), ma.Eval(), ma.DQ.Full())
	ma.LoadData(map[string]float64{"Close": 3.0})
	fmt.Println(ma.EvalOld(), ma.Eval(), ma.DQ.Full())
	ma.LoadData(map[string]float64{"Close": 4.0})
	fmt.Println(ma.EvalOld(), ma.Eval(), ma.DQ.Full())
	ma.LoadData(map[string]float64{"Close": 5.0})
	fmt.Println(ma.EvalOld(), ma.Eval(), ma.DQ.Full())
	ma.LoadData(map[string]float64{"Close": 6.0})
	fmt.Println(ma.EvalOld(), ma.Eval(), ma.DQ.Full())
	ma.LoadData(map[string]float64{"Close": 7.0})
	fmt.Println(ma.EvalOld(), ma.Eval(), ma.DQ.Full())

	if ma.Eval() != 6.0 {
		fmt.Println("ma.Eval() :  ", ma.Eval())
		t.Error("Expected 3.0, got ", ma.Eval())
	}

}
func BenchmarkEvalOldMA(b *testing.B) {
	ma := NewMA("MA60", []int{60}, []string{"Close"})
	for i := 0; i < 60; i++ {
		ma.LoadData(map[string]float64{"Close": float64(i) + 1.0})
	}
	for i := 0; i < b.N; i++ {
		ma.EvalOld() //   216.0 ns/op
	}
}

func BenchmarkEval(b *testing.B) { //迭代式
	ma := NewMA("MA60", []int{60}, []string{"Close"})
	for i := 0; i < 60; i++ {
		ma.LoadData(map[string]float64{"Close": float64(i) + 1.0})

	}
	for i := 0; i < b.N; i++ {
		ma.Eval() // 0.6323 ns/op
	}
}
