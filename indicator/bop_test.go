// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
package indicator

import (
	"fmt"
	"testing"
)

func TestEvalBOP(t *testing.T) {
	bop := NewBOP("BOP", []string{"closing", "opening", "high", "low"})
	bop.LoadData(map[string]float64{"closing": 20, "opening": 10, "high": 40, "low": 4})
	fmt.Println(bop.Eval())
	bop.LoadData(map[string]float64{"closing": 15, "opening": 20, "high": 25, "low": 10})
	fmt.Println(bop.Eval())

	if roundDigits(bop.Eval(), 2) != float64(-0.33) {
		fmt.Println("bop.Eval() :  ", bop.Eval())
		t.Error("Expected -0.33, got ", bop.Eval())
	}

}
