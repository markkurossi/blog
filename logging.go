//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
)

// Verbose prints a verbose output message if the verbose output is
// enabled.
func Verbose(format string, a ...interface{}) {
	if !flagVerbose {
		return
	}
	fmt.Printf(format, a...)
}
