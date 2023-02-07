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

func TestEvalPiercingLine(t *testing.T) {
	h := NewPiercingLine("Piercing Line", []int{}, []string{"High", "Low", "Open", "Close"})
	h.LoadData(map[string]float64{"High": 7.34, "Low": 7.21, "Close": 7.42, "Open": 7.48})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	h.LoadData(map[string]float64{"High": 7.25, "Low": 7.17, "Close": 7.41, "Open": 7.45})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	h.LoadData(map[string]float64{"High": 7.37, "Low": 7.24, "Close": 7.36, "Open": 7.38})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	h.LoadData(map[string]float64{"High": 7.42, "Low": 7.30, "Close": 7.3, "Open": 7.42})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	h.LoadData(map[string]float64{"High": 7.44, "Low": 7.20, "Close": 7.36, "Open": 7.20})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())

	if !h.Eval() {
		fmt.Println("h.Eval() :  ", h.Eval())
		t.Error("Expected true, got ", h.Eval())
	}
}
