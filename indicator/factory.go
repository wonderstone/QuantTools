// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Wonderstone (Digital Office Product Department #2)
// revisor:

package indicator

type IndiInfo struct {
	Name      string
	IndiType  string
	ParSlice  []int
	InfoSlice []string
}

type IIndicator interface {
	LoadData(data map[string]float64)
	Eval() float64
}

// factory pattern
func IndiFactory(ii IndiInfo) IIndicator {
	switch ii.IndiType {
	case "MA":
		return NewMA(ii.ParSlice, ii.InfoSlice)
	case "Var":
		return NewVar(ii.ParSlice, ii.InfoSlice)
	case "EMA":
		return NewEMA(ii.ParSlice, ii.InfoSlice)
	case "BETA":
		return NewBeta(ii.ParSlice, ii.InfoSlice)
	case "Cov":
		return NewConv(ii.ParSlice, ii.InfoSlice)
	case "AD":
		return NewAvgDev(ii.ParSlice, ii.InfoSlice)
	default:
		return nil
	}
}
