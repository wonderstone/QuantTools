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

func TestEvalPSY(t *testing.T) {
	p := NewPSY("PSY", []int{3}, []string{"Close"})
	p.LoadData(map[string]float64{"Close": 7.22, "Volume": 1})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"Close": 7.23, "Volume": 1})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"Close": 7.36, "Volume": 1})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"Close": 7.36, "Volume": 1})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"Close": 7.34, "Volume": 1})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"Close": 7.38, "Volume": 1})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"Close": 7.32, "Volume": 1})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"Close": 7.26, "Volume": 1})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"Close": 7.32, "Volume": 1})
	fmt.Println(p.Eval())
	p.LoadData(map[string]float64{"Close": 7.32, "Volume": 1})
	fmt.Println(p.Eval())

	if roundDigits(p.Eval(), 3) != 33.333 {
		fmt.Println("p.Eval() :  ", p.Eval())
		t.Error("Expected 33.333, got ", p.Eval())
	}
}
