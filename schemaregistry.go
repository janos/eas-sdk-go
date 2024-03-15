// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"

	"resenje.org/eas/internal/contracts"
)

type SchemaRecord struct {
	Uid       UID
	Resolver  common.Address
	Revocable bool
	Schema    string
}

func newSchemaRecord(r *contracts.SchemaRecord) *SchemaRecord {
	return &SchemaRecord{
		Uid:       UID(r.Uid),
		Resolver:  r.Resolver,
		Revocable: r.Revocable,
		Schema:    r.Schema,
	}
}

type SchemaRegistryRegistered struct {
	Uid        UID
	Registerer common.Address
	Raw        types.Log
}

func newSchemaRegistryRegistered(r *contracts.SchemaRegistryRegistered) *SchemaRegistryRegistered {
	return &SchemaRegistryRegistered{
		Uid:        UID(r.Uid),
		Registerer: r.Registerer,
		Raw:        r.Raw,
	}
}

type SchemaRegistryContract struct {
	client   *Client
	contract *contracts.SchemaRegistry
}

func newSchemaRegistryContract(ctx context.Context, client *Client) (*SchemaRegistryContract, error) {
	contractAddress := client.options.SchemaRegistryContractAddress

	var zeroAddress common.Address
	if contractAddress == zeroAddress {
		easContract, err := contracts.NewEAS(client.easContractAddress, client.backend)
		if err != nil {
			return nil, err
		}

		a, err := easContract.GetSchemaRegistry(&bind.CallOpts{Context: ctx})
		if err != nil {
			return nil, err
		}
		contractAddress = a
	}

	contract, err := contracts.NewSchemaRegistry(contractAddress, client.backend)
	if err != nil {
		return nil, err
	}

	return &SchemaRegistryContract{
		client:   client,
		contract: contract,
	}, nil
}

func (c *SchemaRegistryContract) Register(ctx context.Context, schema string, resolver common.Address, revocable bool) (*types.Transaction, WaitTx[SchemaRegistryRegistered], error) {
	txOpts, err := c.client.txOpts(ctx)
	if err != nil {
		return nil, nil, err
	}

	tx, err := c.contract.Register(txOpts, schema, resolver, revocable)
	if err != nil {
		return nil, nil, err
	}

	return tx, newWaitTx(tx, c.client, newParseProxy(c.contract.ParseRegistered, newSchemaRegistryRegistered)), nil
}

func (c *SchemaRegistryContract) GetSchema(ctx context.Context, uid UID) (*SchemaRecord, error) {
	r, err := c.contract.GetSchema(&bind.CallOpts{Context: ctx}, uid)
	if err != nil {
		return nil, err
	}
	return newSchemaRecord(&r), err
}

type schemaRegistryRegisteredIterator struct {
	contracts.SchemaRegistryRegisteredIterator
}

func (i *schemaRegistryRegisteredIterator) Value() *SchemaRegistryRegistered {
	return newSchemaRegistryRegistered(i.Event)
}

func (c *SchemaRegistryContract) FilterRegistered(ctx context.Context, start uint64, end *uint64, uids []UID) (Iterator[*SchemaRegistryRegistered], error) {
	it, err := c.contract.FilterRegistered(&bind.FilterOpts{Start: start, End: end, Context: ctx}, castUIDSlice(uids))
	if err != nil {
		return nil, err
	}
	return &schemaRegistryRegisteredIterator{*it}, nil
}

func (c *SchemaRegistryContract) WatchRegistered(ctx context.Context, start *uint64, sink chan<- *SchemaRegistryRegistered, uids []UID) (event.Subscription, error) {
	proxy := newChanProxy(ctx, sink, newSchemaRegistryRegistered)
	return c.contract.WatchRegistered(&bind.WatchOpts{Start: start, Context: ctx}, proxy, castUIDSlice(uids))
}
