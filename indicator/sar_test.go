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

func TestEvalSAR(t *testing.T) {
	s := NewSAR("SAR5", []int{5}, []string{"High", "Low", "Close"})
	s.LoadData(map[string]float64{"High": 7.34, "Low": 7.21, "Close": 7.22})
	fmt.Println(s.Eval(), s.upTrend)
	s.LoadData(map[string]float64{"High": 7.25, "Low": 7.17, "Close": 7.23})
	fmt.Println(s.Eval(), s.upTrend)
	s.LoadData(map[string]float64{"High": 7.37, "Low": 7.24, "Close": 7.36})
	fmt.Println(s.Eval(), s.upTrend)
	s.LoadData(map[string]float64{"High": 7.39, "Low": 7.30, "Close": 7.36})
	fmt.Println(s.Eval(), s.upTrend)
	s.LoadData(map[string]float64{"High": 7.40, "Low": 7.30, "Close": 7.34})
	fmt.Println(s.Eval(), s.upTrend)
	s.LoadData(map[string]float64{"High": 7.39, "Low": 7.31, "Close": 7.38})
	fmt.Println(s.Eval(), s.upTrend)
	s.LoadData(map[string]float64{"High": 7.39, "Low": 7.30, "Close": 7.32})
	fmt.Println(s.Eval(), s.upTrend)
	s.LoadData(map[string]float64{"High": 7.36, "Low": 7.23, "Close": 7.26})
	fmt.Println(s.Eval(), s.upTrend)
	s.LoadData(map[string]float64{"High": 7.35, "Low": 7.25, "Close": 7.32})
	fmt.Println(s.Eval(), s.upTrend)
	s.LoadData(map[string]float64{"High": 7.37, "Low": 7.26, "Close": 7.32})
	fmt.Println(s.Eval(), s.upTrend)

	if roundDigits(s.Eval(), 3) != 7.38 {
		fmt.Println("s.Eval() :  ", s.Eval())
		t.Error("Expected 7.38, got ", s.Eval())
	}
}
