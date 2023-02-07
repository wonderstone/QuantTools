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

func TestEvalCYE(t *testing.T) {
	c := NewCYE("CYE5", []int{5}, []string{"Close"})
	c.LoadData(map[string]float64{"High": 7.34, "Low": 7.21, "Close": 7.22, "Open": 7.28})
	fmt.Println(c.Eval())
	c.LoadData(map[string]float64{"High": 7.25, "Low": 7.17, "Close": 7.23, "Open": 7.25})
	fmt.Println(c.Eval())
	c.LoadData(map[string]float64{"High": 7.37, "Low": 7.24, "Close": 7.36, "Open": 7.28})
	fmt.Println(c.Eval())
	c.LoadData(map[string]float64{"High": 7.39, "Low": 7.30, "Close": 7.36, "Open": 7.32})
	fmt.Println(c.Eval())
	c.LoadData(map[string]float64{"High": 7.40, "Low": 7.30, "Close": 7.34, "Open": 7.35})
	fmt.Println(c.Eval())
	c.LoadData(map[string]float64{"High": 7.39, "Low": 7.31, "Close": 7.38, "Open": 7.32})
	fmt.Println(c.Eval())
	c.LoadData(map[string]float64{"High": 7.39, "Low": 7.30, "Close": 7.32, "Open": 7.37})
	fmt.Println(c.Eval())
	c.LoadData(map[string]float64{"High": 7.36, "Low": 7.23, "Close": 7.26, "Open": 7.31})
	fmt.Println(c.Eval())
	c.LoadData(map[string]float64{"High": 7.35, "Low": 7.25, "Close": 7.32, "Open": 7.27})
	fmt.Println(c.Eval())
	c.LoadData(map[string]float64{"High": 7.37, "Low": 7.26, "Close": 7.32, "Open": 7.37})
	fmt.Println(c.Eval())

	if roundDigits(c.Eval(), 3) != -0.055 {
		fmt.Println("c.Eval() :  ", c.Eval())
		t.Error("Expected -*0.055, got ", roundDigits(c.Eval(), 3))
	}
}
