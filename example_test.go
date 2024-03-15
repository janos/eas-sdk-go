// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas_test

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"resenje.org/eas"
)

const (
	// network
	endpointSepolia        = "https://ethereum-sepolia-rpc.publicnode.com/"
	contractAddressSepolia = "0xC2679fBD37d54388Ce493F1DB75320D236e1815e"

	// account
	privateKeyHex = "933c798b990a6be3fb91ae2fd3b6593f61d6d478548091205ee948b1de9c9f19"
)

func Example() {
	privateKey, err := eas.HexParsePrivateKey(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(contractAddressSepolia)

	c, err := eas.NewClient(context.Background(), endpointSepolia, privateKey, contractAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

	attestationUID := eas.HexDecodeUID("0xac812932f5cee90a457d57a9fbd7b142b21ba99b809f982bbf86947f295281ff")

	a, err := c.EAS.GetAttestation(context.Background(), attestationUID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Attester", a.Attester, a.Time, a.Schema)

	var schemaUID eas.UID
	var name string

	if err := a.ScanValues(&schemaUID, &name); err != nil {
		log.Fatal(err)
	}

	log.Println("Attestation")
	log.Println("Schema UID:", schemaUID)
	log.Println("Name:", name)

	// Output:
}
