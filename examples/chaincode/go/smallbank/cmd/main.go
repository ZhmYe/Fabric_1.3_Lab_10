package cmd

import "fmt"

/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
"fmt"

"github.com/hyperledger/fabric/core/chaincode/shim"
"github.com/hyperledger/fabric/examples/chaincode/go/example01"
)

func main() {
	err := shim.Start(new(smallbank.SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
