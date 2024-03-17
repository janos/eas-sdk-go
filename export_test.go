// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"resenje.org/eas/internal/contracts"
)

type DeployBackend interface {
	bind.ContractBackend
	bind.DeployBackend
	ChainID(context.Context) (*big.Int, error)
}

func DeployEAS(ctx context.Context, backend DeployBackend, pk *ecdsa.PrivateKey, registry common.Address) (common.Address, *types.Transaction, *contracts.EAS, WaitDeployment, error) {
	txOpts, err := newTxOpts(ctx, backend, pk)
	if err != nil {
		return common.Address{}, nil, nil, nil, err
	}

	address, tx, contract, err := contracts.DeployEAS(txOpts, backend, registry)
	if err != nil {
		return common.Address{}, nil, nil, nil, err
	}

	return address, tx, contract, func(ctx context.Context) (common.Address, error) {
		return bind.WaitDeployed(ctx, backend, tx)
	}, nil
}

func DeploySchemaRegistry(ctx context.Context, backend DeployBackend, pk *ecdsa.PrivateKey) (common.Address, *types.Transaction, *contracts.SchemaRegistry, WaitDeployment, error) {
	txOpts, err := newTxOpts(ctx, backend, pk)
	if err != nil {
		return common.Address{}, nil, nil, nil, err
	}

	address, tx, contract, err := contracts.DeploySchemaRegistry(txOpts, backend)
	if err != nil {
		return common.Address{}, nil, nil, nil, err
	}

	return address, tx, contract, func(ctx context.Context) (common.Address, error) {
		return bind.WaitDeployed(ctx, backend, tx)
	}, nil
}

type WaitDeployment func(ctx context.Context) (common.Address, error)

func newTxOpts(ctx context.Context, backend DeployBackend, pk *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	chainID, err := backend.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		return nil, fmt.Errorf("construct transactor: %w", err)
	}

	publicKeyECDSA, ok := pk.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not a valid ecdsa public key")
	}

	from := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := backend.PendingNonceAt(ctx, from)
	if err != nil {
		return nil, fmt.Errorf("get padding nonce: %w", err)
	}

	gasPrice, err := backend.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("suggest gas price: %w", err)
	}

	opts.Nonce = big.NewInt(int64(nonce))
	opts.Value = big.NewInt(0)
	opts.GasLimit = 30000000
	opts.GasPrice = gasPrice
	opts.Context = ctx

	return opts, nil
}
