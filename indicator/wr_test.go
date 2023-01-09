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

func TestEvalWR(t *testing.T) {
	w := NewWR("WR4", []int{4}, []string{"Close", "Low", "High"})
	w.LoadData(map[string]float64{"Close": 1.0, "High": 10.0, "Low": 0.0})
	fmt.Println(w.Eval())
	w.LoadData(map[string]float64{"Close": 2.0, "High": 4.0, "Low": 1.0})
	fmt.Println(w.Eval())
	w.LoadData(map[string]float64{"Close": 3.0, "High": 5.0, "Low": 2.0})
	fmt.Println(w.Eval())
	w.LoadData(map[string]float64{"Close": 4.0, "High": 6.0, "Low": 3.0})
	fmt.Println(w.Eval())
	w.LoadData(map[string]float64{"Close": 4.0, "High": 8.0, "Low": 4.0})
	fmt.Println(w.Eval())

	if roundDigits(w.Eval(), 3) != 57.143 {
		fmt.Println("w.Eval() :  ", w.Eval())
		t.Error("Expected 57.142, got ", w.Eval())
	}
}
