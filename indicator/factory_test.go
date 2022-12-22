// All rights reserved. This is part of West Securities ltd. proprietary source code.
// No part of this file may be reproduced or transmitted in any form or by any means,
// electronic or mechanical, including photocopying, recording, or by any information
// storage and retrieval system, without the prior written permission of West Securities ltd.

// author:  Wonderstone (Digital Office Product Department #2)
// revisor:

package indicator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the indicator factory
func TestFactory(t *testing.T) {
	// Test the indicator factory
	indis := []IndiInfo{
		{"MA10", "MA", []int{3}, []string{"Close"}},
		{"Var10", "Var", []int{3}, []string{"Close"}},
	}

	for _, indi := range indis {
		indicator := IndiFactory(indi)
		assert.NotNil(t, indicator, "Indicator should not be nil")
	}

}