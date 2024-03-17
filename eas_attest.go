// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"

	"resenje.org/eas/internal/contracts"
)

type Attestation struct {
	UID            UID
	Schema         UID
	Time           time.Time
	ExpirationTime time.Time
	RevocationTime time.Time
	RefUID         UID
	Recipient      common.Address
	Attester       common.Address
	Revocable      bool
	Data           []byte
}

func newAttestation(a *contracts.Attestation) *Attestation {
	return &Attestation{
		UID:            a.Uid,
		Schema:         a.Schema,
		Time:           time.Unix(int64(a.Time), 0),
		ExpirationTime: time.Unix(int64(a.ExpirationTime), 0),
		RevocationTime: time.Unix(int64(a.RevocationTime), 0),
		RefUID:         a.RefUID,
		Recipient:      a.Recipient,
		Attester:       a.Attester,
		Revocable:      a.Revocable,
		Data:           a.Data,
	}
}

func (a Attestation) Fields(schema string) ([]SchemaItem, error) {
	return decodeAttestationValues(a.Data, schema)
}

func (a Attestation) ScanValues(fields ...any) error {
	return scanAttestationValues(a.Data, fields...)
}

type EASAttested struct {
	Recipient common.Address
	Attester  common.Address
	UID       UID
	Schema    UID
	Raw       types.Log
}

func newEASAttested(r *contracts.EASAttested) *EASAttested {
	return &EASAttested{
		Recipient: r.Recipient,
		Attester:  r.Attester,
		UID:       r.Uid,
		Schema:    r.Schema,
		Raw:       r.Raw,
	}
}

type AttestOptions struct {
	Recipient      common.Address
	ExpirationTime time.Time
	Revocable      bool
	RefUID         UID
	Value          *big.Int
}

func newAttestationRequestData(data []byte, o *AttestOptions) contracts.AttestationRequestData {
	if o == nil {
		o = new(AttestOptions)
	}
	if o.Value == nil {
		o.Value = big.NewInt(0)
	}
	var expirationTime uint64
	if !o.ExpirationTime.IsZero() {
		expirationTime = uint64(o.ExpirationTime.Unix())
	}
	return contracts.AttestationRequestData{
		Recipient:      o.Recipient,
		ExpirationTime: expirationTime,
		Revocable:      o.Revocable,
		RefUID:         o.RefUID,
		Data:           data,
		Value:          o.Value,
	}
}

func (c *EASContract) Attest(ctx context.Context, schemaUID UID, o *AttestOptions, values ...any) (*types.Transaction, WaitTx[EASAttested], error) {
	txOpts, err := c.client.newTxOpts(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("construct transaction options: %w", err)
	}

	data, err := encodeAttestationValues(values)
	if err != nil {
		return nil, nil, fmt.Errorf("encode attestation values: %w", err)
	}

	tx, err := c.contract.Attest(txOpts, contracts.AttestationRequest{
		Schema: schemaUID,
		Data:   newAttestationRequestData(data, o),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("call attest contract method: %w", err)
	}

	return tx, newWaitTx(tx, c.client, newParseProxy(c.contract.ParseAttested, newEASAttested)), nil
}

func (c *EASContract) MultiAttest(ctx context.Context, schemaUID UID, o *AttestOptions, attestations ...[]any) (*types.Transaction, WaitTx[EASAttested], error) {
	txOpts, err := c.client.newTxOpts(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("construct transaction options: %w", err)
	}

	var data []contracts.AttestationRequestData
	for _, a := range attestations {
		d, err := encodeAttestationValues(a)
		if err != nil {
			return nil, nil, err
		}
		data = append(data, newAttestationRequestData(d, o))
	}

	tx, err := c.contract.MultiAttest(txOpts, []contracts.MultiAttestationRequest{
		{
			Schema: schemaUID,
			Data:   data,
		},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("call multi attest contract method: %w", err)
	}

	return tx, newWaitTx(tx, c.client, newParseProxy(c.contract.ParseAttested, newEASAttested)), nil
}

func (c *EASContract) GetAttestation(ctx context.Context, uid UID) (*Attestation, error) {
	a, err := c.contract.GetAttestation(&bind.CallOpts{Context: ctx}, uid)
	if err != nil {
		return nil, err
	}
	return newAttestation(&a), nil
}

func (c *EASContract) IsAttestationValid(ctx context.Context, uid UID) (bool, error) {
	return c.contract.IsAttestationValid(&bind.CallOpts{Context: ctx}, uid)
}

type easAttestedIterator struct {
	contracts.EASAttestedIterator
}

func (i *easAttestedIterator) Value() EASAttested {
	return *newEASAttested(i.Event)
}

func (c *EASContract) FilterAttested(ctx context.Context, start uint64, end *uint64, recipient []common.Address, attester []common.Address, schema []UID) (Iterator[EASAttested], error) {
	it, err := c.contract.FilterAttested(&bind.FilterOpts{Start: start, End: end, Context: ctx}, recipient, attester, castUIDSlice(schema))
	if err != nil {
		return nil, err
	}
	return &easAttestedIterator{*it}, nil
}

func (c *EASContract) WatchAttested(ctx context.Context, start *uint64, sink chan<- *EASAttested, recipient []common.Address, attester []common.Address, schema []UID) (event.Subscription, error) {
	return c.contract.WatchAttested(&bind.WatchOpts{Start: start, Context: ctx}, newChanProxy(ctx, sink, newEASAttested), recipient, attester, castUIDSlice(schema))
}
