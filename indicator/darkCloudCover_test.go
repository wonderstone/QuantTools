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

func TestEvalDarkCloudCover(t *testing.T) {
	d := NewDarkCloudCover("Dark Cloud Cover", []int{}, []string{"High", "Low", "Open", "Close"})
	d.LoadData(map[string]float64{"High": 7.40, "Low": 7.22, "Close": 7.32, "Open": 7.37})
	fmt.Println(d.cye.Eval(), d.single.Eval(), d.Eval())
	d.LoadData(map[string]float64{"High": 7.37, "Low": 7.24, "Close": 7.36, "Open": 7.38})
	fmt.Println(d.cye.Eval(), d.single.Eval(), d.Eval())
	d.LoadData(map[string]float64{"High": 7.68, "Low": 7.48, "Close": 7.41, "Open": 7.55})
	fmt.Println(d.cye.Eval(), d.single.Eval(), d.Eval())
	d.LoadData(map[string]float64{"High": 7.89, "Low": 7.40, "Close": 7.62, "Open": 7.40})
	fmt.Println(d.cye.Eval(), d.single.Eval(), d.Eval())
	d.LoadData(map[string]float64{"High": 7.70, "Low": 7.30, "Close": 7.44, "Open": 7.70})
	fmt.Println(d.cye.Eval(), d.single.Eval(), d.Eval())

	if !d.Eval() {
		fmt.Println("d.Eval() :  ", d.Eval())
		t.Error("Expected true, got ", d.Eval())
	}
}
