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
	"reflect"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	defaultGasFeeCap = big.NewInt(20_000_000_000) // 20 Gwei
	defaultGasTipCap = big.NewInt(2_000_000_000)  // 2 Gwei
)

type Backend interface {
	bind.ContractBackend
	bind.DeployBackend
	ethereum.ChainIDReader
}

type Client struct {
	backend            Backend
	from               common.Address
	pk                 *ecdsa.PrivateKey
	easContractAddress common.Address
	options            *Options

	chainID *big.Int

	// Contracts
	SchemaRegistry *SchemaRegistryContract
	EAS            *EASContract
}

type Options struct {
	SchemaRegistryContractAddress common.Address
	GasLimit                      uint64
	GasFeeCap                     *big.Int
	GasTipCap                     *big.Int
	Backend                       Backend
}

func NewClient(ctx context.Context, endpoint string, pk *ecdsa.PrivateKey, easContractAddress common.Address, o *Options) (*Client, error) {
	publicKeyECDSA, ok := pk.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not a valid ecdsa public key")
	}

	if o == nil {
		o = new(Options)
	}

	backend := o.Backend
	if backend == nil {
		b, err := ethclient.DialContext(ctx, endpoint)
		if err != nil {
			return nil, fmt.Errorf("connect to endpoint: %w", err)
		}
		backend = b
	}

	c := &Client{
		backend:            backend,
		from:               crypto.PubkeyToAddress(*publicKeyECDSA),
		pk:                 pk,
		easContractAddress: easContractAddress,
		options:            o,
	}

	chainID, err := c.backend.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}
	c.chainID = chainID

	schemaRegistryContract, err := newSchemaRegistryContract(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("construct schema registry contract: %w", err)
	}
	c.SchemaRegistry = schemaRegistryContract

	easContract, err := newEASContract(c)
	if err != nil {
		return nil, fmt.Errorf("construct eas contract: %w", err)
	}
	c.EAS = easContract

	return c, nil
}

func (c *Client) Backend() Backend {
	return c.backend
}

func Ptr[T any](v T) *T {
	return &v
}

type TxOptions []TxOption

type TxOption func(*bind.TransactOpts)

func GasLimit(v uint64) TxOption {
	return func(to *bind.TransactOpts) {
		to.GasLimit = v
	}
}

func GasFeeCap(v *big.Int) TxOption {
	return func(to *bind.TransactOpts) {
		to.GasFeeCap = v
	}
}

func GasTipCap(v *big.Int) TxOption {
	return func(to *bind.TransactOpts) {
		to.GasTipCap = v
	}
}

func (c *Client) newTxOpts(ctx context.Context, opts TxOptions) (*bind.TransactOpts, error) {
	nonce, err := c.backend.PendingNonceAt(ctx, c.from)
	if err != nil {
		return nil, fmt.Errorf("get padding nonce: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(c.pk, c.chainID)
	if err != nil {
		return nil, fmt.Errorf("construct transactor: %w", err)
	}
	auth.Context = ctx
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = c.options.GasLimit
	auth.GasFeeCap = defaultGasFeeCap
	if c.options.GasFeeCap != nil {
		auth.GasFeeCap = c.options.GasFeeCap
	}
	auth.GasTipCap = defaultGasTipCap
	if c.options.GasTipCap != nil {
		auth.GasTipCap = c.options.GasTipCap
	}

	for _, o := range opts {
		o(auth)
	}

	return auth, nil
}

type WaitTx[T any] func(ctx context.Context) (*T, error)

func newWaitTx[T any](tx *types.Transaction, client *Client, parse func(log types.Log) (*T, error)) WaitTx[T] {
	return func(ctx context.Context) (*T, error) {
		receipt, err := bind.WaitMined(ctx, client.backend, tx)
		if err != nil {
			return nil, err
		}

		return parse(*receipt.Logs[0])
	}
}

type Iterator[T any] interface {
	Value() T
	Close() error
	Error() error
	Next() bool
}

func newChanProxy[I, O any](ctx context.Context, sink chan<- O, constructor func(I) O) chan I {
	proxy := make(chan I)
	go func() {
		for {
			select {
			case v, ok := <-proxy:
				if !ok {
					return
				}
				select {
				case sink <- constructor(v):
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return proxy
}

func newParseProxy[I, O any](parse func(log types.Log) (I, error), constructor func(I) O) func(log types.Log) (O, error) {
	return func(log types.Log) (O, error) {
		r, err := parse(log)
		if err != nil {
			var o O
			return o, err
		}
		return constructor(r), nil
	}
}

func castUIDSlice(s []UID) [][32]byte {
	r := make([][32]byte, 0, len(s))
	for _, u := range s {
		r = append(r, u)
	}
	return r
}

func getTypeString(v any, suffix string) (string, error) {
	switch v.(type) {
	case common.Address:
		return "address" + suffix, nil
	case string:
		return "string" + suffix, nil
	case bool:
		return "bool" + suffix, nil
	case [32]byte, UID:
		return "bytes32" + suffix, nil
	case []byte:
		return "bytes" + suffix, nil
	case uint8:
		return "uint8" + suffix, nil
	case uint16:
		return "uint16" + suffix, nil
	case uint32:
		return "uint32" + suffix, nil
	case uint64:
		return "uint64" + suffix, nil
	case *big.Int:
		return "uint256" + suffix, nil
	default:
		t := reflect.TypeOf(v)
		switch t.Kind() {
		case reflect.Array:
			e := reflect.New(t.Elem()).Interface()
			return getTypeString(e, "[]"+suffix)
		case reflect.Slice:
			len := reflect.ValueOf(v).Len()
			e := reflect.New(t.Elem()).Interface()
			return getTypeString(e, fmt.Sprintf("[%v]", len)+suffix)
		default:
			return "", fmt.Errorf("unsupported type %T", v)
		}
	}
}
