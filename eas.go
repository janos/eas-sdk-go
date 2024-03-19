// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"resenje.org/eas/internal/contracts"
)

type EASContract struct {
	client   *Client
	contract *contracts.EAS
	abi      *abi.ABI
}

func newEASContract(client *Client) (*EASContract, error) {
	contract, err := contracts.NewEAS(client.easContractAddress, client.backend)
	if err != nil {
		return nil, fmt.Errorf("construct abi bindings: %w", err)
	}

	abi, err := contracts.EASMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("get abi: %w", err)
	}

	return &EASContract{
		client:   client,
		contract: contract,
		abi:      abi,
	}, nil
}

func (c *EASContract) unpackError(err error) error {
	return unpackError(err, c.abi)
}

func (c *EASContract) Version(ctx context.Context) (string, error) {
	v, err := c.contract.Version(&bind.CallOpts{Context: ctx})
	if err != nil {
		return "", c.unpackError(err)
	}
	return v, nil
}
