// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:

package indicator

import (
	"fmt"
	"testing"
)

func TestEvalRSI(t *testing.T) {
	r := NewRSI("RSI6", []int{6}, []string{"Close"})
	r.LoadData(map[string]float64{"Close": 1.0})
	fmt.Println(r.Eval())
	r.LoadData(map[string]float64{"Close": 2.0})
	fmt.Println(r.Eval())
	r.LoadData(map[string]float64{"Close": 3.0})
	fmt.Println(r.Eval())
	r.LoadData(map[string]float64{"Close": 4.0})
	fmt.Println(r.Eval())
	r.LoadData(map[string]float64{"Close": 3.0})
	fmt.Println(r.Eval())
	r.LoadData(map[string]float64{"Close": 2.0})
	fmt.Println(r.Eval())
	r.LoadData(map[string]float64{"Close": 1.0})
	fmt.Println(r.Eval())
	r.LoadData(map[string]float64{"Close": 0.0})
	fmt.Println(r.Eval())

	if roundDigits(r.Eval(), 3) != 33.333 {
		fmt.Println("r.Eval() :  ", r.Eval())
		t.Error("Expected 33.333, got ", r.Eval())
	}
}
