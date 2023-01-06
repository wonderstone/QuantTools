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

func TestEvalRef(t *testing.T) {
	ref := NewRef("Ref3", []int{3}, []string{"Close"})

	ref.LoadData(map[string]float64{"Close": 1.0})
	fmt.Println(ref.Eval())
	ref.LoadData(map[string]float64{"Close": 2.0})
	fmt.Println(ref.Eval())
	ref.LoadData(map[string]float64{"Close": 3.0})
	fmt.Println(ref.Eval())
	ref.LoadData(map[string]float64{"Close": 4.0})
	fmt.Println(ref.Eval())

	if ref.Eval() != 2.0 {
		fmt.Println("ref.Eval() :  ", ref.Eval())
		t.Error("Expected 2.0, got ", ref.Eval())
	}

}
