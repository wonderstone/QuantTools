// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Wonderstone (Digital Office Product Department #2)
// revisor:
package indicator

import (
	"fmt"
	"testing"
)

func TestAVGDEVEval(t *testing.T) {
	avgdev := NewAvgDev([]int{5}, []string{"Close"})
	avgdev.LoadData(map[string]float64{"Close": 1.0})
	fmt.Println(avgdev.Eval(), avgdev.DQ.Full())
	avgdev.LoadData(map[string]float64{"Close": 2.0})
	fmt.Println(avgdev.Eval(), avgdev.DQ.Full())
	avgdev.LoadData(map[string]float64{"Close": 3.0})
	fmt.Println(avgdev.Eval(), avgdev.DQ.Full())
	avgdev.LoadData(map[string]float64{"Close": 4.0})
	fmt.Println(avgdev.Eval(), avgdev.DQ.Full())
	avgdev.LoadData(map[string]float64{"Close": 5.0})
	fmt.Println(avgdev.Eval(), avgdev.DQ.Full())
	avgdev.LoadData(map[string]float64{"Close": 6.0})
	fmt.Println(avgdev.Eval(), avgdev.DQ.Full())

	if avgdev.Eval() != 1.2 {
		fmt.Println("avgdev.Eval() :  ", avgdev.Eval())
		t.Error("Expected --- , got ", avgdev.Eval())
	}
}
