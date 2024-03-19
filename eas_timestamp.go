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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"resenje.org/eas/internal/contracts"
)

type Timestamp uint64

func (t Timestamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}

func castTimestampSlice(s []Timestamp) []uint64 {
	r := make([]uint64, 0, len(s))
	for _, t := range s {
		r = append(r, uint64(t))
	}
	return r
}

type EASTimestamped struct {
	Data      UID
	Timestamp Timestamp
	Raw       types.Log
}

func newEASTimestamped(r *contracts.EASTimestamped) *EASTimestamped {
	return &EASTimestamped{
		Data:      r.Data,
		Timestamp: Timestamp(r.Timestamp),
		Raw:       r.Raw,
	}
}

func (c *EASContract) Timestamp(ctx context.Context, data UID) (*types.Transaction, WaitTx[EASTimestamped], error) {
	txOpts, err := c.client.newTxOpts(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("construct transaction options: %w", err)
	}

	tx, err := c.contract.Timestamp(txOpts, data)
	if err != nil {
		return nil, nil, fmt.Errorf("call timestamp contract method: %w", c.unpackError(err))
	}

	return tx, newWaitTx(tx, c.client, newParseProxy(c.contract.ParseTimestamped, newEASTimestamped)), nil
}

func (c *EASContract) MultiTimestamp(ctx context.Context, data []UID) (*types.Transaction, WaitTxMulti[EASTimestamped], error) {
	txOpts, err := c.client.newTxOpts(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("construct transaction options: %w", err)
	}

	tx, err := c.contract.MultiTimestamp(txOpts, castUIDSlice(data))
	if err != nil {
		return nil, nil, fmt.Errorf("call multi timestamp contract method: %w", c.unpackError(err))
	}
	return tx, newWaitTxMulti(tx, c.client, newParseProxy(c.contract.ParseTimestamped, newEASTimestamped)), nil
}

func (c *EASContract) GetTimestamp(ctx context.Context, data UID) (Timestamp, error) {
	timestamp, err := c.contract.GetTimestamp(&bind.CallOpts{Context: ctx}, data)
	if err != nil {
		return 0, c.unpackError(err)
	}
	return Timestamp(timestamp), nil
}

type easTimestampedIterator struct {
	contracts.EASTimestampedIterator
}

func (i *easTimestampedIterator) Value() EASTimestamped {
	return *newEASTimestamped(i.Event)
}

func (c *EASContract) FilterTimestamped(ctx context.Context, start uint64, end *uint64, data []UID, timestamps []Timestamp) (Iterator[EASTimestamped], error) {
	it, err := c.contract.FilterTimestamped(&bind.FilterOpts{Start: start, End: end, Context: ctx}, castUIDSlice(data), castTimestampSlice(timestamps))
	if err != nil {
		return nil, c.unpackError(err)
	}
	return &easTimestampedIterator{*it}, nil
}

func (c *EASContract) WatchTimestamped(ctx context.Context, start *uint64, sink chan<- *EASTimestamped, data []UID, timestamps []Timestamp) (event.Subscription, error) {
	s, err := c.contract.WatchTimestamped(&bind.WatchOpts{Start: start, Context: ctx}, newChanProxy(ctx, sink, newEASTimestamped), castUIDSlice(data), castTimestampSlice(timestamps))
	if err != nil {
		return nil, c.unpackError(err)
	}
	return s, nil
}
