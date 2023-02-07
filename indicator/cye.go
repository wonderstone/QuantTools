// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:

package indicator

import "math"

type CYE struct {
	Name      string
	ParSlice  []int
	InfoSlice []string
	lv, val   float64
	ma        *MA
}

func NewCYE(Name string, ParSlice []int, InfoSlice []string) *CYE {
	return &CYE{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: InfoSlice,
		ma:        NewMA(Name, ParSlice, InfoSlice),
		val:       0,
		lv:        0,
	}
}

func (c *CYE) LoadData(data map[string]float64) {
	if c.val != 0 {
		c.lv = c.val
	}
	c.ma.LoadData(data)
	c.val = c.ma.Eval()
}

func (c *CYE) Eval() float64 {
	if c.lv == 0 {
		return math.NaN()
	}
	return 100 * (c.val - c.lv) / c.lv
}

func (c *CYE) GetName() string {
	return c.Name
}
