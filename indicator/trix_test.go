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

func TestEvalTRIX(t *testing.T) {
	trix := NewTRIX("TRIX3", []int{3}, []string{"Close"})
	trix.LoadData(map[string]float64{"Close": 8996.96})
	fmt.Println(trix.Eval())
	trix.LoadData(map[string]float64{"Close": 9003.19})
	fmt.Println(trix.Eval())
	trix.LoadData(map[string]float64{"Close": 9010.41})
	fmt.Println(trix.Eval())
	trix.LoadData(map[string]float64{"Close": 9008.07})
	fmt.Println(trix.Eval())
	trix.LoadData(map[string]float64{"Close": 9018.03})
	fmt.Println(trix.Eval())
	trix.LoadData(map[string]float64{"Close": 9009.80})
	fmt.Println(trix.Eval())
	trix.LoadData(map[string]float64{"Close": 9011.55})
	fmt.Println(trix.Eval())
	trix.LoadData(map[string]float64{"Close": 9020.26})
	fmt.Println(trix.Eval())
	trix.LoadData(map[string]float64{"Close": 9013.29})
	fmt.Println(trix.Eval())
	trix.LoadData(map[string]float64{"Close": 9018.78})
	fmt.Println(trix.Eval())

	if roundDigits(trix.Eval(), 4) != 0.0137 {
		fmt.Println("trix.Eval() :  ", roundDigits(trix.Eval(), 4))
		t.Error("Expected 0.137, got ", trix.Eval())
	}
}
