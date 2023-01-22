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

func TestEvalOBV(t *testing.T) {
	o := NewOBV("OBV", []int{}, []string{"Close", "Volume"})
	o.LoadData(map[string]float64{"Close": 7.22, "Volume": 1})
	fmt.Println(o.Eval())
	o.LoadData(map[string]float64{"Close": 7.23, "Volume": 1})
	fmt.Println(o.Eval())
	o.LoadData(map[string]float64{"Close": 7.36, "Volume": 1})
	fmt.Println(o.Eval())
	o.LoadData(map[string]float64{"Close": 7.36, "Volume": 1})
	fmt.Println(o.Eval())
	o.LoadData(map[string]float64{"Close": 7.34, "Volume": 1})
	fmt.Println(o.Eval())
	o.LoadData(map[string]float64{"Close": 7.38, "Volume": 1})
	fmt.Println(o.Eval())
	o.LoadData(map[string]float64{"Close": 7.32, "Volume": 1})
	fmt.Println(o.Eval())
	o.LoadData(map[string]float64{"Close": 7.26, "Volume": 1})
	fmt.Println(o.Eval())
	o.LoadData(map[string]float64{"Close": 7.32, "Volume": 1})
	fmt.Println(o.Eval())
	o.LoadData(map[string]float64{"Close": 7.32, "Volume": 1})
	fmt.Println(o.Eval())

	if o.Eval() != 1 {
		fmt.Println("o.Eval() :  ", o.Eval())
		t.Error("Expected 1, got ", o.Eval())
	}
}
