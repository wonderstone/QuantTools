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

func TestEvalBearishAbandonedBaby(t *testing.T) {
	h := NewBearishAbandonedBaby("Bearish Abandoned Baby", []int{}, []string{"High", "Low", "Open", "Close"})
	h.LoadData(map[string]float64{"High": 7.34, "Low": 7.21, "Close": 7.25, "Open": 7.22})
	fmt.Println(h.Eval())
	h.LoadData(map[string]float64{"High": 7.25, "Low": 7.17, "Close": 7.28, "Open": 7.23})
	fmt.Println(h.Eval())
	h.LoadData(map[string]float64{"High": 7.37, "Low": 7.24, "Close": 7.30, "Open": 7.26})
	fmt.Println(h.Eval())
	h.LoadData(map[string]float64{"High": 7.39, "Low": 7.30, "Close": 7.36, "Open": 7.30})
	fmt.Println(h.Eval())
	h.LoadData(map[string]float64{"High": 7.40, "Low": 7.34, "Close": 7.40, "Open": 7.38})
	fmt.Println(h.Eval())
	h.LoadData(map[string]float64{"High": 7.39, "Low": 7.31, "Close": 7.32, "Open": 7.37})
	fmt.Println(h.Eval())

	if !h.Eval() {
		fmt.Println("h.Eval() :  ", h.Eval())
		t.Error("Expected true, got ", h.Eval())
	}
}
