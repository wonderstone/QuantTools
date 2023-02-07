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

func TestEvalZIG(t *testing.T) {
	z := NewZIG("ZIG", []int{3, 5}, []string{"Open", "High", "Low", "Close"})
	z.LoadData(map[string]float64{"Close": 100, "Volume": 1})
	fmt.Println(z.Eval())
	z.LoadData(map[string]float64{"Close": 96, "Volume": 1})
	fmt.Println(z.Eval())
	z.LoadData(map[string]float64{"Close": 104, "Volume": 1})
	fmt.Println(z.Eval())
	z.LoadData(map[string]float64{"Close": 105, "Volume": 1})
	fmt.Println(z.Eval())
	z.LoadData(map[string]float64{"Close": 103, "Volume": 1})
	fmt.Println(z.Eval())
	z.LoadData(map[string]float64{"Close": 102, "Volume": 1})
	fmt.Println(z.Eval())
	z.LoadData(map[string]float64{"Close": 101, "Volume": 1})
	fmt.Println(z.Eval())
	z.LoadData(map[string]float64{"Close": 109, "Volume": 1})
	fmt.Println(z.Eval())
	z.LoadData(map[string]float64{"Close": 103, "Volume": 1})
	fmt.Println(z.Eval())
	z.LoadData(map[string]float64{"Close": 110, "Volume": 1})
	fmt.Println(z.Eval())
	z.LoadData(map[string]float64{"Close": 105, "Volume": 1})
	fmt.Println(z.Eval())

	if z.Eval() != 117 {
		fmt.Println("z.Eval() :  ", z.Eval())
		t.Error("Expected 117, got ", z.Eval())
	}
}
