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

func TestEvalEMV(t *testing.T) {
	e := NewEMV("EMV9", []int{9}, []string{"High", "Low", "Volume"})
	e.LoadData(map[string]float64{"Close": 1.0, "High": 10.0, "Low": 0.0, "Volume": 1})
	fmt.Println(e.Eval())
	e.LoadData(map[string]float64{"Close": 2.0, "High": 4.0, "Low": 1.0, "Volume": 1})
	fmt.Println(e.Eval())
	e.LoadData(map[string]float64{"Close": 3.0, "High": 5.0, "Low": 2.0, "Volume": 1})
	fmt.Println(e.Eval())
	e.LoadData(map[string]float64{"Close": 4.0, "High": 6.0, "Low": 3.0, "Volume": 1})
	fmt.Println(e.Eval())
	e.LoadData(map[string]float64{"Close": 4.0, "High": 8.0, "Low": 4.0, "Volume": 1})
	fmt.Println(e.Eval())

	if e.Eval() != 54.5 {
		fmt.Println("e.Eval() :  ", e.Eval())
		t.Error("Expected 54.5, got ", e.Eval())
	}
}
