// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas

import (
	"fmt"

	"resenje.org/eas/internal/contracts"
)

type EASContract struct {
	client   *Client
	contract *contracts.EAS
}

func newEASContract(client *Client) (*EASContract, error) {
	contract, err := contracts.NewEAS(client.easContractAddress, client.backend)
	if err != nil {
		return nil, fmt.Errorf("construct eas abi bindings: %w", err)
	}

	return &EASContract{
		client:   client,
		contract: contract,
	}, nil
}
