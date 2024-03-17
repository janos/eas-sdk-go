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

func TestSchemaRegistryContract_Register(t *testing.T) {
	client, backend := newClient(t)
	ctx := context.Background()

	schema := "byte32 uid, string secret"

	_, wait, err := client.SchemaRegistry.Register(ctx, schema, common.Address{1}, true)
	assertNilError(t, err)

	backend.Commit()

	r, err := wait(ctx)
	assertNilError(t, err)

	s, err := client.SchemaRegistry.GetSchema(ctx, r.UID)
	assertNilError(t, err)

	assertEqual(t, "uid", s.UID, r.UID)
	assertEqual(t, "schema", s.Schema, schema)
	assertEqual(t, "resolver", s.Resolver, common.Address{1})
	assertEqual(t, "revocable", s.Revocable, true)
}

func TestSchemaRegistryContract_FilterRegistered(t *testing.T) {
	client, backend := newClient(t)
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

		backend.Commit()

		r, err := wait(ctx)
		assertNilError(t, err)

		uids = append(uids, r.UID)
	}

	t.Run("all", func(t *testing.T) {
		it, err := client.SchemaRegistry.FilterRegistered(ctx, 0, nil, nil)
		assertNilError(t, err)
		defer it.Close()

		schemaCount := 0

		for it.Next() {
			r := it.Value()

			s, err := client.SchemaRegistry.GetSchema(ctx, r.UID)
			assertNilError(t, err)

			assertEqual(t, "uid", s.UID, uids[schemaCount])
			assertEqual(t, "schema", s.Schema, schemas[schemaCount])
			assertEqual(t, "resolver", s.Resolver, common.Address{})
			assertEqual(t, "revocable", s.Revocable, true)

			schemaCount++
		}
		assertNilError(t, it.Error())

		assertEqual(t, "schema count", schemaCount, len(schemas))
	})

	t.Run("filter blocks", func(t *testing.T) {
		// start from the third block after adding schemas
		it, err := client.SchemaRegistry.FilterRegistered(ctx, blockNumber+3, eas.Ptr(blockNumber+4), nil)
		assertNilError(t, err)
		defer it.Close()

		schemaCount := 0

		for it.Next() {
			r := it.Value()

			s, err := client.SchemaRegistry.GetSchema(ctx, r.UID)
			assertNilError(t, err)

			assertEqual(t, "uid", s.UID, uids[schemaCount+2])
			assertEqual(t, "schema", s.Schema, schemas[schemaCount+2]) // ignore first two schemas
			assertEqual(t, "resolver", s.Resolver, common.Address{})
			assertEqual(t, "revocable", s.Revocable, true)

			schemaCount++
		}
		assertNilError(t, it.Error())

		assertEqual(t, "schema count", schemaCount, 2)
	})
}

func TestSchemaRegistryContract_WatchRegistered(t *testing.T) {
	client, backend := newClient(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sink := make(chan *eas.SchemaRegistryRegistered)
	sub, err := client.SchemaRegistry.WatchRegistered(ctx, nil, sink, nil)
	assertNilError(t, err)
	defer sub.Unsubscribe()

	schemaCount := 0

	schemas := []string{
		"bool ignore",
		"byte32 uid, string secret",
		"uint256[] number",
		"uint64[] id",
		"uint64[] id, string text",
	}

	go func() {
		defer close(sink)

		for _, schema := range schemas {
			_, wait, err := client.SchemaRegistry.Register(ctx, schema, common.Address{}, true)
			if err != nil {
				t.Error(err)
			}

			backend.Commit()

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
		case r, ok := <-sink:
			if !ok {
				break loop
			}
			s, err := client.SchemaRegistry.GetSchema(ctx, r.UID)
			assertNilError(t, err)

			assertEqual(t, "uid", s.UID, r.UID)
			assertEqual(t, "schema", s.Schema, schemas[schemaCount])
			assertEqual(t, "resolver", s.Resolver, common.Address{})
			assertEqual(t, "revocable", s.Revocable, true)

			schemaCount++
		case err := <-sub.Err():
			if err != nil {
				t.Fatal(err)
			}
		case <-ctx.Done():
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	assertEqual(t, "schema count", schemaCount, 5)
}
