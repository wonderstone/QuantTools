// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Zhangweixuan (Digital Office Product Department #2)
// revisor:

package indicator

type ROC struct {
	Name      string
	ParSlice  []int
	InfoSlice []string
	today     float64
	ref       *Ref
}

func NewROC(Name string, ParSlice []int, InfoSlice []string) *ROC {
	return &ROC{
		Name:      Name,
		ParSlice:  ParSlice,
		InfoSlice: InfoSlice,
		today:     0,
		ref:       NewRef(Name, []int{ParSlice[0] + 1}, []string{"Close"}),
	}
}

// LoadData loads 1 tick info datas into the indicator
func (r *ROC) LoadData(data map[string]float64) {
	r.today = data[r.InfoSlice[0]]
	r.ref.LoadData(data)
}

func (r *ROC) Eval() float64 {
	return (r.today - r.ref.Eval()) / r.ref.Eval()
}

func (r *ROC) GetName() string {
	return r.Name
}
