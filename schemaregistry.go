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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"

	"resenje.org/eas/internal/contracts"
)

type SchemaRecord struct {
	UID       UID
	Resolver  common.Address
	Revocable bool
	Schema    string
}

func newSchemaRecord(r *contracts.SchemaRecord) *SchemaRecord {
	return &SchemaRecord{
		UID:       UID(r.Uid),
		Resolver:  r.Resolver,
		Revocable: r.Revocable,
		Schema:    r.Schema,
	}
}

type SchemaRegistryRegistered struct {
	UID        UID
	Registerer common.Address
	Raw        types.Log
}

func newSchemaRegistryRegistered(r *contracts.SchemaRegistryRegistered) *SchemaRegistryRegistered {
	return &SchemaRegistryRegistered{
		UID:        UID(r.Uid),
		Registerer: r.Registerer,
		Raw:        r.Raw,
	}
}

type SchemaRegistryContract struct {
	client   *Client
	contract *contracts.SchemaRegistry
	abi      *abi.ABI
}

func newSchemaRegistryContract(ctx context.Context, client *Client) (*SchemaRegistryContract, error) {
	contractAddress := client.options.SchemaRegistryContractAddress

	var zeroAddress common.Address
	if contractAddress == zeroAddress {
		easContract, err := contracts.NewEAS(client.easContractAddress, client.backend)
		if err != nil {
			return nil, fmt.Errorf("construct abi bindings: %w", err)
		}

		a, err := easContract.GetSchemaRegistry(&bind.CallOpts{Context: ctx})
		if err != nil {
			return nil, fmt.Errorf("get schema registry address: %w", err)
		}
		contractAddress = a
	}

	contract, err := contracts.NewSchemaRegistry(contractAddress, client.backend)
	if err != nil {
		return nil, fmt.Errorf("construct abi bindings: %w", err)
	}

	abi, err := contracts.EASMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("get abi: %w", err)
	}

	return &SchemaRegistryContract{
		client:   client,
		contract: contract,
		abi:      abi,
	}, nil
}

func (c *SchemaRegistryContract) parseError(err error) error {
	return unpackError(err, c.abi)
}

func (c *SchemaRegistryContract) Version(ctx context.Context) (string, error) {
	v, err := c.contract.Version(&bind.CallOpts{Context: ctx})
	if err != nil {
		return "", c.parseError(err)
	}
	return v, nil
}

func (c *SchemaRegistryContract) Register(ctx context.Context, schema string, resolver common.Address, revocable bool) (*types.Transaction, WaitTx[SchemaRegistryRegistered], error) {
	txOpts, err := c.client.newTxOpts(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("construct transaction options: %w", err)
	}

	tx, err := c.contract.Register(txOpts, schema, resolver, revocable)
	if err != nil {
		return nil, nil, fmt.Errorf("call register contract method: %w", c.parseError(err))
	}

	return tx, newWaitTx(tx, c.client, newParseProxy(c.contract.ParseRegistered, newSchemaRegistryRegistered)), nil
}

func (c *SchemaRegistryContract) GetSchema(ctx context.Context, uid UID) (*SchemaRecord, error) {
	r, err := c.contract.GetSchema(&bind.CallOpts{Context: ctx}, uid)
	if err != nil {
		return nil, c.parseError(err)
	}
	return newSchemaRecord(&r), nil
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
		return nil, c.parseError(err)
	}
	return &schemaRegistryRegisteredIterator{*it}, nil
}

func (c *SchemaRegistryContract) WatchRegistered(ctx context.Context, start *uint64, sink chan<- *SchemaRegistryRegistered, uids []UID) (event.Subscription, error) {
	s, err := c.contract.WatchRegistered(&bind.WatchOpts{Start: start, Context: ctx}, newChanProxy(ctx, sink, newSchemaRegistryRegistered), castUIDSlice(uids))
	if err != nil {
		return nil, c.parseError(err)
	}
	return s, nil
}
