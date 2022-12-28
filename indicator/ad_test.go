// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Maminghui (Digital Office Product Department #2)
// revisor:
package indicator

import (
	"fmt"
	"testing"
)

func TestEvalAD(t *testing.T) {
	ad := NewAD([]int{3}, []string{"high", "low", "close", "volume"})
	ad.LoadData(map[string]float64{"high": 10, "low": 6, "close": 9, "volume": 100})
	fmt.Println(ad.Eval())
	ad.LoadData(map[string]float64{"high": 9, "low": 7, "close": 11, "volume": 110})
	fmt.Println(ad.Eval())
	ad.LoadData(map[string]float64{"high": 12, "low": 9, "close": 7, "volume": 80})
	actual := ad.Eval()
	fmt.Println(actual)
	actual_2 := roundDigits(actual, 2)

	if actual_2 != float64(193.33) {
		fmt.Println("ad.Eval() :  ", actual_2)
		t.Error("Expected 193.33, got ", actual_2)
	}

}
