// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"resenje.org/eas/internal/contracts"
)

type EASTimestamped struct {
	Data      UID
	Timestamp time.Time
	Raw       types.Log
}

func newEASTimestamped(r *contracts.EASTimestamped) *EASTimestamped {
	return &EASTimestamped{
		Data:      r.Data,
		Timestamp: time.Unix(int64(r.Timestamp), 0),
		Raw:       r.Raw,
	}
}

func (c *EASContract) Timestamp(ctx context.Context, data UID) (*types.Transaction, WaitTx[EASTimestamped], error) {
	txOpts, err := c.client.txOpts(ctx)
	if err != nil {
		return nil, nil, err
	}

	tx, err := c.contract.Timestamp(txOpts, data)
	if err != nil {
		return nil, nil, err
	}

	return tx, newWaitTx(tx, c.client, newParseProxy(c.contract.ParseTimestamped, newEASTimestamped)), nil
}

func (c *EASContract) MultiTimestamp(ctx context.Context, schemaUID UID, data []UID) (*types.Transaction, WaitTx[EASTimestamped], error) {
	txOpts, err := c.client.txOpts(ctx)
	if err != nil {
		return nil, nil, err
	}

	tx, err := c.contract.MultiTimestamp(txOpts, castUIDSlice(data))
	if err != nil {
		return nil, nil, err
	}
	return tx, newWaitTx(tx, c.client, newParseProxy(c.contract.ParseTimestamped, newEASTimestamped)), nil
}

func (c *EASContract) GetTimestamp(ctx context.Context, data UID) (uint64, error) {
	id, err := c.contract.GetTimestamp(&bind.CallOpts{Context: ctx}, data)
	if err != nil {
		return 0, err
	}
	return id, nil
}

type easTimestampedIterator struct {
	contracts.EASTimestampedIterator
}

func (i *easTimestampedIterator) Value() EASTimestamped {
	return *newEASTimestamped(i.Event)
}

func (c *EASContract) FilterTimestamped(ctx context.Context, start uint64, end *uint64, data [][32]byte, timestamp []uint64) (Iterator[EASTimestamped], error) {
	it, err := c.contract.FilterTimestamped(&bind.FilterOpts{Start: start, End: end, Context: ctx}, data, timestamp)
	if err != nil {
		return nil, err
	}
	return &easTimestampedIterator{*it}, nil
}

func (c *EASContract) WatchTimestamped(ctx context.Context, start *uint64, sink chan<- *EASTimestamped, data [][32]byte, timestamp []uint64) (event.Subscription, error) {
	return c.contract.WatchTimestamped(&bind.WatchOpts{Start: start, Context: ctx}, newChanProxy(ctx, sink, newEASTimestamped), data, timestamp)
}
