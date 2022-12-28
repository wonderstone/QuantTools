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

func TestEvalSum(t *testing.T) {
	s := NewSum([]int{4}, []string{"Close"})
	s.LoadData(map[string]float64{"Close": 1.0})
	fmt.Println(s.Eval())
	s.LoadData(map[string]float64{"Close": 2.0})
	fmt.Println(s.Eval())
	s.LoadData(map[string]float64{"Close": 3.0})
	fmt.Println(s.Eval())
	s.LoadData(map[string]float64{"Close": 4.0})
	fmt.Println(s.Eval())
	s.LoadData(map[string]float64{"Close": 5.0})
	fmt.Println(s.Eval())

	if s.Eval() != float64(14) {
		fmt.Println("s.Eval() :  ", s.Eval())
		t.Error("Expected 14, got ", s.Eval())
	}

}
