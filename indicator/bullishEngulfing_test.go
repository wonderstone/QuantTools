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

func TestEvalBullishEngulfing(t *testing.T) {
	h := NewBullishEngulfing("Bullish Engulfing", []int{}, []string{"High", "Low", "Open", "Close"})
	h.LoadData(map[string]float64{"High": 7.89, "Low": 7.48, "Close": 7.56, "Open": 7.62})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	h.LoadData(map[string]float64{"High": 7.68, "Low": 7.48, "Close": 7.41, "Open": 7.55})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	h.LoadData(map[string]float64{"High": 7.37, "Low": 7.24, "Close": 7.36, "Open": 7.38})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	h.LoadData(map[string]float64{"High": 7.40, "Low": 7.22, "Close": 7.32, "Open": 7.37})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	h.LoadData(map[string]float64{"High": 7.42, "Low": 7.05, "Close": 7.38, "Open": 7.25})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())

	if !h.Eval() {
		fmt.Println("h.Eval() :  ", h.Eval())
		t.Error("Expected true, got ", h.Eval())
	}
}
