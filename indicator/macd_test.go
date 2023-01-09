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

func TestEvalMACD(t *testing.T) {
	m := NewMACD("MACD", []int{}, []string{"Close"})
	m.LoadData(map[string]float64{"Close": 1.0})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"Close": 2.0})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"Close": 3.0})
	fmt.Println(m.Eval())
	m.LoadData(map[string]float64{"Close": 4.0})
	fmt.Println(m.Eval())

	if roundDigits(m.Eval(), 3) != 0.563 {
		fmt.Println("m.Eval() :  ", m.Eval())
		t.Error("Expected 0.563, got ", m.Eval())
	}
}
