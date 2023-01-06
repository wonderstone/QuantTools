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

func TestEvalVar(t *testing.T) {
	vari := NewVar("Var3", []int{3}, []string{"Close"})

	vari.LoadData(map[string]float64{"Close": 1.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())
	vari.LoadData(map[string]float64{"Close": 2.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())
	vari.LoadData(map[string]float64{"Close": 3.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())
	vari.LoadData(map[string]float64{"Close": 4.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())
	vari.LoadData(map[string]float64{"Close": 5.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())
	vari.LoadData(map[string]float64{"Close": 6.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())
	vari.LoadData(map[string]float64{"Close": 7.0})
	fmt.Println(vari.Eval(), vari.DQ.Full())

	if vari.Eval() != float64(2.0/3.0) {
		fmt.Println("vari.Eval() :  ", vari.Eval())
		t.Error("Expected 3.0, got ", vari.Eval())
	}

}
