// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
// revisor:
package indicator

import (
	"fmt"
	"testing"
)

func TestEvalEMA(t *testing.T) {
	e := NewEMA([]int{12}, []string{"Close"})
	e.LoadData(map[string]float64{"Close": 1.0})
	fmt.Println(e.ptoday, e.Eval())
	e.LoadData(map[string]float64{"Close": 2.0})
	fmt.Println(e.ptoday, e.Eval())
	e.LoadData(map[string]float64{"Close": 3.0})
	fmt.Println(e.ptoday, e.Eval())
	e.LoadData(map[string]float64{"Close": 4.0})
	actual := e.Eval()
	fmt.Println(e.ptoday, actual)
	actual_3 := roundDigits(actual, 3)

	if actual_3 != float64(1.832) {
		fmt.Println("e.Eval() :  ", actual_3)
		t.Error("Expected 1.832, got ", actual_3)
	}

}
