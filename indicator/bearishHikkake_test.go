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

func TestEvalBearishHikkake(t *testing.T) {
	d := NewBearishHikkake("Bearish Hikkake", []int{}, []string{"High", "Low", "Open", "Close"})
	d.LoadData(map[string]float64{"High": 7.99, "Low": 7.22, "Close": 7.90, "Open": 7.37})
	fmt.Println(d.cye.Eval(), d.single.Eval(), d.Eval())
	d.LoadData(map[string]float64{"High": 7.99, "Low": 7.22, "Close": 7.76, "Open": 7.79})
	fmt.Println(d.cye.Eval(), d.single.Eval(), d.Eval())
	d.LoadData(map[string]float64{"High": 7.99, "Low": 7.22, "Close": 7.65, "Open": 7.75})
	fmt.Println(d.cye.Eval(), d.single.Eval(), d.Eval())
	d.LoadData(map[string]float64{"High": 7.99, "Low": 7.22, "Close": 7.40, "Open": 7.45})
	fmt.Println(d.cye.Eval(), d.single.Eval(), d.Eval())
	d.LoadData(map[string]float64{"High": 7.99, "Low": 7.22, "Close": 7.70, "Open": 7.45})
	fmt.Println(d.cye.Eval(), d.single.Eval(), d.Eval())

	if !d.Eval() {
		fmt.Println("d.Eval() :  ", d.Eval())
		t.Error("Expected true, got ", d.Eval())
	}
}
