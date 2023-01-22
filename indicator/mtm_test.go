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

func TestEvalMTM(t *testing.T) {
	m := NewMTM("MTM6", []int{6}, []string{"Close"})
	m.LoadData(map[string]float64{"Close": 7.22})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"Close": 7.23})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"Close": 7.36})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"Close": 7.36})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"Close": 7.34})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"Close": 7.38})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"Close": 7.32})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"Close": 7.26})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"Close": 7.32})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"Close": 7.32})
	fmt.Println(m.Eval())

	if roundDigits(m.Eval(), 2) != -0.04 {
		fmt.Println("m.Eval() :  ", m.Eval())
		t.Error("Expected -0.04, got ", m.Eval())
	}
}
