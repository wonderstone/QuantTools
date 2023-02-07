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

func TestEvalBearishEngulfing(t *testing.T) {
	h := NewBearishEngulfing("Bearish Engulfing", []int{}, []string{"High", "Low", "Open", "Close"})
	h.LoadData(map[string]float64{"High": 7.37, "Low": 7.24, "Close": 7.26, "Open": 7.30})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	h.LoadData(map[string]float64{"High": 7.40, "Low": 7.22, "Close": 7.37, "Open": 7.32})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	h.LoadData(map[string]float64{"High": 7.42, "Low": 7.05, "Close": 7.38, "Open": 7.25})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	h.LoadData(map[string]float64{"High": 7.68, "Low": 7.48, "Close": 7.55, "Open": 7.41})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	h.LoadData(map[string]float64{"High": 7.89, "Low": 7.60, "Close": 7.69, "Open": 7.56})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	h.LoadData(map[string]float64{"High": 7.78, "Low": 7.44, "Close": 7.49, "Open": 7.72})
	fmt.Println(h.cye.Eval(), h.single.Eval(), h.Eval())
	if !h.Eval() {
		fmt.Println("h.Eval() :  ", h.Eval())
		t.Error("Expected true, got ", h.Eval())
	}
}
