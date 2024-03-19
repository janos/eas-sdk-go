// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"resenje.org/eas"
)

func TestSchemaRegistryContract_Version(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	version, err := client.SchemaRegistry.Version(ctx)
	assertNilError(t, err)

	assertEqual(t, "version", version, "1.0.0")
}

func TestSchemaRegistryContract_Register(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	schema := "byte32 uid, string secret"

	_, wait, err := client.SchemaRegistry.Register(ctx, schema, common.Address{1}, true)
	assertNilError(t, err)

	client.backend.Commit()

	r, err := wait(ctx)
	assertNilError(t, err)

	assertEqual(t, "registerer", r.Registerer, client.account)

	s, err := client.SchemaRegistry.GetSchema(ctx, r.UID)
	assertNilError(t, err)

	assertEqual(t, "uid", s.UID, r.UID)
	assertEqual(t, "schema", s.Schema, schema)
	assertEqual(t, "resolver", s.Resolver, common.Address{1})
	assertEqual(t, "revocable", s.Revocable, true)
}

func TestSchemaRegistryContract_FilterRegistered(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	blockNumber, err := client.Backend().(ethereum.BlockNumberReader).BlockNumber(ctx)
	assertNilError(t, err)

	schemas := []string{
		"bool ignore",               // ignore
		"byte32 uid, string secret", // ignore
		"uint256[] number",
		"uint64[] id",
		"uint64[] id, string text", // ignore
	}

	var uids []eas.UID

	for _, schema := range schemas {
		_, wait, err := client.SchemaRegistry.Register(ctx, schema, common.Address{}, true)
		assertNilError(t, err)

		client.backend.Commit()

		r, err := wait(ctx)
		assertNilError(t, err)

		uids = append(uids, r.UID)
	}

	t.Run("all", func(t *testing.T) {
		it, err := client.SchemaRegistry.FilterRegistered(ctx, 0, nil, nil)
		assertNilError(t, err)
		defer it.Close()

		count := 0

		for it.Next() {
			r := it.Value()

			s, err := client.SchemaRegistry.GetSchema(ctx, r.UID)
			assertNilError(t, err)

			assertEqual(t, "uid", s.UID, uids[count])
			assertEqual(t, "schema", s.Schema, schemas[count])
			assertEqual(t, "resolver", s.Resolver, common.Address{})
			assertEqual(t, "revocable", s.Revocable, true)

			count++
		}
		assertNilError(t, it.Error())

		assertEqual(t, "count", count, len(schemas))
	})

	t.Run("filter blocks", func(t *testing.T) {
		// start from the third block after adding schemas
		it, err := client.SchemaRegistry.FilterRegistered(ctx, blockNumber+3, eas.Ptr(blockNumber+4), nil)
		assertNilError(t, err)
		defer it.Close()

		count := 0

		for it.Next() {
			r := it.Value()

			s, err := client.SchemaRegistry.GetSchema(ctx, r.UID)
			assertNilError(t, err)

			assertEqual(t, "uid", s.UID, uids[count+2])
			assertEqual(t, "schema", s.Schema, schemas[count+2]) // ignore first two schemas
			assertEqual(t, "resolver", s.Resolver, common.Address{})
			assertEqual(t, "revocable", s.Revocable, true)

			count++
		}
		assertNilError(t, it.Error())

		assertEqual(t, "count", count, 2)
	})

	t.Run("filter uids", func(t *testing.T) {
		it, err := client.SchemaRegistry.FilterRegistered(ctx, 0, nil, []eas.UID{uids[1], uids[3]})
		assertNilError(t, err)
		defer it.Close()

		count := 0

		wantUID := 1
		for it.Next() {
			r := it.Value()

			s, err := client.SchemaRegistry.GetSchema(ctx, r.UID)
			assertNilError(t, err)

			assertEqual(t, "uid", s.UID, uids[wantUID])
			assertEqual(t, "schema", s.Schema, schemas[wantUID])
			assertEqual(t, "resolver", s.Resolver, common.Address{})
			assertEqual(t, "revocable", s.Revocable, true)

			count++
			wantUID = 3
		}
		assertNilError(t, it.Error())

		assertEqual(t, "count", count, 2)
	})
}

func TestSchemaRegistryContract_WatchRegistered(t *testing.T) {
	client := newClient(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sink := make(chan *eas.SchemaRegistryRegistered)
	sub, err := client.SchemaRegistry.WatchRegistered(ctx, nil, sink, nil)
	assertNilError(t, err)
	defer sub.Unsubscribe()

	count := 0

	schemas := []string{
		"bool ignore",
		"byte32 uid, string secret",
		"uint256[] number",
		"uint64[] id",
		"uint64[] id, string text",
	}

	go func() {
		defer sub.Unsubscribe()

		for _, schema := range schemas {
			_, wait, err := client.SchemaRegistry.Register(ctx, schema, common.Address{}, true)
			if err != nil {
				t.Error(err)
			}

			client.backend.Commit()

			if _, err := wait(ctx); err != nil {
				t.Error(err)
			}

			select {
			case <-time.After(100 * time.Millisecond):
			case <-ctx.Done():
			}
		}
	}()

loop:
	for {
		select {
		case r := <-sink:
			s, err := client.SchemaRegistry.GetSchema(ctx, r.UID)
			assertNilError(t, err)

			assertEqual(t, "uid", s.UID, r.UID)
			assertEqual(t, "schema", s.Schema, schemas[count])
			assertEqual(t, "resolver", s.Resolver, common.Address{})
			assertEqual(t, "revocable", s.Revocable, true)

			count++
		case err, ok := <-sub.Err():
			if !ok {
				break loop
			}
			if err != nil {
				t.Fatal(err)
			}
		case <-ctx.Done():
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	assertEqual(t, "count", count, 5)
}

func newSchema(t *testing.T, client *Client, schema string) eas.UID {
	t.Helper()
	ctx := context.Background()

	_, wait, err := client.SchemaRegistry.Register(ctx, schema, common.Address{}, true)
	assertNilError(t, err)

	client.backend.Commit()

	r, err := wait(ctx)
	assertNilError(t, err)

	return r.UID
}
