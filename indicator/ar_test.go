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

func TestEvalAR(t *testing.T) {
	br := NewAR("BR5", []int{5}, []string{"High", "Low", "Open"})
	br.LoadData(map[string]float64{"High": 7.34, "Low": 7.21, "Open": 7.22})
	fmt.Println(br.Eval())
	br.LoadData(map[string]float64{"High": 7.25, "Low": 7.17, "Open": 7.23})
	fmt.Println(br.Eval())
	br.LoadData(map[string]float64{"High": 7.37, "Low": 7.24, "Open": 7.36})
	fmt.Println(br.Eval())
	br.LoadData(map[string]float64{"High": 7.39, "Low": 7.30, "Open": 7.36})
	fmt.Println(br.Eval())
	br.LoadData(map[string]float64{"High": 7.40, "Low": 7.30, "Open": 7.34})
	fmt.Println(br.Eval())
	br.LoadData(map[string]float64{"High": 7.39, "Low": 7.31, "Open": 7.38})
	fmt.Println(br.Eval())
	br.LoadData(map[string]float64{"High": 7.39, "Low": 7.30, "Open": 7.32})
	fmt.Println(br.Eval())
	br.LoadData(map[string]float64{"High": 7.36, "Low": 7.23, "Open": 7.26})
	fmt.Println(br.Eval())
	br.LoadData(map[string]float64{"High": 7.35, "Low": 7.25, "Open": 7.32})
	fmt.Println(br.Eval())
	br.LoadData(map[string]float64{"High": 7.37, "Low": 7.26, "Open": 7.32})
	fmt.Println(br.Eval())

	if roundDigits(br.Eval(), 2) != 104.00 {
		fmt.Println("br.Eval() :  ", br.Eval())
		t.Error("Expected 104.00, got ", br.Eval())
	}
}
