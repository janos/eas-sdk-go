// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"

	"resenje.org/eas/internal/contracts"
)

type EASRevoked struct {
	Recipient common.Address
	Attester  common.Address
	Uid       UID
	Schema    UID
	Raw       types.Log
}

func newEASRevoked(r *contracts.EASRevoked) *EASRevoked {
	return &EASRevoked{
		Recipient: r.Recipient,
		Attester:  r.Attester,
		Uid:       r.Uid,
		Schema:    r.Schema,
		Raw:       r.Raw,
	}
}

type RevokeOptions struct {
	Value *big.Int
}

func newRevocationRequestData(attestationUID UID, o *RevokeOptions) contracts.RevocationRequestData {
	if o == nil {
		o = new(RevokeOptions)
	}
	if o.Value == nil {
		o.Value = big.NewInt(0)
	}
	return contracts.RevocationRequestData{
		Uid:   attestationUID,
		Value: o.Value,
	}
}

func (c *EASContract) Revoke(ctx context.Context, schemaUID, attestationUID UID, o *RevokeOptions) (*types.Transaction, WaitTx[EASRevoked], error) {
	txOpts, err := c.client.txOpts(ctx)
	if err != nil {
		return nil, nil, err
	}

	tx, err := c.contract.Revoke(txOpts, contracts.RevocationRequest{
		Schema: schemaUID,
		Data:   newRevocationRequestData(attestationUID, o),
	})
	if err != nil {
		return nil, nil, err
	}

	return tx, newWaitTx(tx, c.client, newParseProxy(c.contract.ParseRevoked, newEASRevoked)), nil
}

func (c *EASContract) MultiRevoke(ctx context.Context, schemaUID UID, o *AttestOptions, attestationUIDs []UID) (*types.Transaction, WaitTx[EASRevoked], error) {
	txOpts, err := c.client.txOpts(ctx)
	if err != nil {
		return nil, nil, err
	}

	var data []contracts.RevocationRequestData
	for _, u := range attestationUIDs {
		data = append(data, newRevocationRequestData(u, nil))
	}

	tx, err := c.contract.MultiRevoke(txOpts, []contracts.MultiRevocationRequest{
		{
			Schema: schemaUID,
			Data:   data,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return tx, newWaitTx(tx, c.client, newParseProxy(c.contract.ParseRevoked, newEASRevoked)), nil
}

type easRevokedIterator struct {
	contracts.EASRevokedIterator
}

func (i *easRevokedIterator) Value() EASRevoked {
	return *newEASRevoked(i.Event)
}

func (c *EASContract) FilterRevoked(ctx context.Context, start uint64, end *uint64, recipient []common.Address, attester []common.Address, schema []UID) (Iterator[EASRevoked], error) {
	it, err := c.contract.FilterRevoked(&bind.FilterOpts{Start: start, End: end, Context: ctx}, recipient, attester, castUIDSlice(schema))
	if err != nil {
		return nil, err
	}
	return &easRevokedIterator{*it}, nil
}

func (c *EASContract) WatchRevoked(ctx context.Context, start *uint64, sink chan<- *EASRevoked, recipient []common.Address, attester []common.Address, schema []UID) (event.Subscription, error) {
	return c.contract.WatchRevoked(&bind.WatchOpts{Start: start, Context: ctx}, newChanProxy(ctx, sink, newEASRevoked), recipient, attester, castUIDSlice(schema))
}
