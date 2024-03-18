// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"

	"resenje.org/eas/internal/contracts"
)

type EASRevokedOffchain struct {
	Revoker   common.Address
	Data      UID
	Timestamp time.Time
	Raw       types.Log
}

func newEASRevokedOffchain(r *contracts.EASRevokedOffchain) *EASRevokedOffchain {
	return &EASRevokedOffchain{
		Revoker:   r.Revoker,
		Data:      r.Data,
		Timestamp: time.Unix(int64(r.Timestamp), 0),
		Raw:       r.Raw,
	}
}

func (c *EASContract) RevokeOffchain(ctx context.Context, opts TxOptions, uid UID) (*types.Transaction, WaitTx[EASRevokedOffchain], error) {
	txOpts, err := c.client.newTxOpts(ctx, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("construct transaction options: %w", err)
	}

	tx, err := c.contract.RevokeOffchain(txOpts, uid)
	if err != nil {
		return nil, nil, fmt.Errorf("call revoke offchain contract method: %w", err)
	}

	return tx, newWaitTx(tx, c.client, newParseProxy(c.contract.ParseRevokedOffchain, newEASRevokedOffchain)), nil
}

func (c *EASContract) MultiRevokeOffchain(ctx context.Context, opts TxOptions, schemaUID UID, uids []UID) (*types.Transaction, WaitTx[EASRevokedOffchain], error) {
	txOpts, err := c.client.newTxOpts(ctx, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("construct transaction options: %w", err)
	}

	tx, err := c.contract.MultiRevokeOffchain(txOpts, castUIDSlice(uids))
	if err != nil {
		return nil, nil, fmt.Errorf("call multiple revoke offchain contract method: %w", err)
	}
	return tx, newWaitTx(tx, c.client, newParseProxy(c.contract.ParseRevokedOffchain, newEASRevokedOffchain)), nil
}

func (c *EASContract) GetRevokeOffchain(ctx context.Context, revoker common.Address, uid UID) (uint64, error) {
	id, err := c.contract.GetRevokeOffchain(&bind.CallOpts{Context: ctx}, revoker, uid)
	if err != nil {
		return 0, err
	}
	return id, nil
}

type easRevokedOffchainIterator struct {
	contracts.EASRevokedOffchainIterator
}

func (i *easRevokedOffchainIterator) Value() EASRevokedOffchain {
	return *newEASRevokedOffchain(i.Event)
}

func (c *EASContract) FilterRevokedOffchain(ctx context.Context, start uint64, end *uint64, revoker []common.Address, data []UID, timestamp []uint64) (Iterator[EASRevokedOffchain], error) {
	it, err := c.contract.FilterRevokedOffchain(&bind.FilterOpts{Start: start, End: end, Context: ctx}, revoker, castUIDSlice(data), timestamp)
	if err != nil {
		return nil, err
	}
	return &easRevokedOffchainIterator{*it}, nil
}

func (c *EASContract) WatchRevokedOffchain(ctx context.Context, start *uint64, sink chan<- *EASRevokedOffchain, revoker []common.Address, data []UID, timestamp []uint64) (event.Subscription, error) {
	return c.contract.WatchRevokedOffchain(&bind.WatchOpts{Start: start, Context: ctx}, newChanProxy(ctx, sink, newEASRevokedOffchain), revoker, castUIDSlice(data), timestamp)
}
