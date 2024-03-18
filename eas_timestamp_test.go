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
	"resenje.org/eas"
)

func TestEASContract_Timestamp(t *testing.T) {
	client, backend := newClient(t)
	ctx := context.Background()

	uid := eas.HexDecodeUID("0x00cf55814f6a9135f3458c52cb5cf8bc511e8f22b3e8688f0f52229bc12a4d15")

	_, wait, err := client.EAS.Timestamp(ctx, nil, uid)
	assertNilError(t, err)

	backend.Commit()

	r, err := wait(ctx)
	assertNilError(t, err)

	if time.Since(r.Timestamp.Time()) > time.Minute {
		t.Errorf("too old timestamp time %v", r.Timestamp)
	}
	assertEqual(t, "data", uid, r.Data)
}

func TestEASContract_GetTimestamp(t *testing.T) {
	client, backend := newClient(t)
	ctx := context.Background()

	uid := eas.HexDecodeUID("0x00cf55814f6a9135f3458c52cb5cf8bc511e8f22b3e8688f0f52229bc12a4d15")

	_, wait, err := client.EAS.Timestamp(ctx, nil, uid)
	assertNilError(t, err)

	backend.Commit()

	want, err := wait(ctx)
	assertNilError(t, err)

	timestamp, err := client.EAS.GetTimestamp(ctx, uid)
	assertNilError(t, err)

	assertEqual(t, "timestamp", timestamp, want.Timestamp)
}

func TestEASContract_MultiTimestamp(t *testing.T) {
	client, backend := newClient(t)
	ctx := context.Background()

	uids := []eas.UID{
		eas.HexDecodeUID("0xc06e2558b09b9846931ff3968fbd241e865167d1d9c46a77f8115daebe065ea9"),
		eas.HexDecodeUID("0xeea93ae9950464785cd53cfd61cfa9c5b0b04ed0d8c48efadf1b3d35813a49d6"),
		eas.HexDecodeUID("0x200ca4d156a9e8fcf8bcf55814f235f3458c52cb5cf5112b3f02fe86885529bc"),
	}

	_, wait, err := client.EAS.MultiTimestamp(ctx, nil, uids)
	assertNilError(t, err)

	backend.Commit()

	r, err := wait(ctx)
	assertNilError(t, err)

	for _, uid := range uids {
		timestamp, err := client.EAS.GetTimestamp(ctx, uid)
		assertNilError(t, err)

		assertEqual(t, "timestamp", timestamp, r.Timestamp)
	}
}

func TestEASContract_FilterRegistered(t *testing.T) {
	client, backend := newClient(t)
	ctx := context.Background()

	blockNumber, err := client.Backend().(ethereum.BlockNumberReader).BlockNumber(ctx)
	assertNilError(t, err)

	uids := []eas.UID{
		eas.HexDecodeUID("0x00cf55814f6a9135f3458c52cb8bc511e8f22b3e864f29325f52229bc12a4d15"),
		eas.HexDecodeUID("0x1605a2096b0bb3dd2ed53e268d0fd797271040a574dc1c995c3560357f9d7432"),
		eas.HexDecodeUID("0xcb8d8ccf0d47f2269ef59fcc16d6d4cdd64021b28c2edf41efce79dcb9c29a17"),
		eas.HexDecodeUID("0x3380b8d16d19cecce426d45e24cef6d1dc159a0f8d35808889f112cb2e003794"),
		eas.HexDecodeUID("0xa821c015ec279aa37500b66a4a65f3230cbeb5b0c707e4a9265771a67346bbe2"),
		eas.HexDecodeUID("0x204e9bc8300e1206dc9b41d096555ab70ec4423fbbb6a9b96a30691f30f673ba"),
		eas.HexDecodeUID("0x28ae72581c95e4664a63eaa5036a8f5e8515ed9f7c8a149058c9fcb3e8216da9"),
	}

	timestamps := make([]eas.Timestamp, 0, len(uids))

	for _, u := range uids {
		_, wait, err := client.EAS.Timestamp(ctx, nil, u)
		assertNilError(t, err)

		backend.Commit()

		r, err := wait(ctx)
		assertNilError(t, err)

		timestamps = append(timestamps, r.Timestamp)
	}

	t.Run("all", func(t *testing.T) {
		it, err := client.EAS.FilterTimestamped(ctx, 0, nil, nil, nil)
		assertNilError(t, err)
		defer it.Close()

		count := 0

		for it.Next() {
			r := it.Value()

			timestamp, err := client.EAS.GetTimestamp(ctx, uids[count])
			assertNilError(t, err)

			assertEqual(t, "timestamp", r.Timestamp, timestamp)
			assertEqual(t, "timestamp", r.Timestamp, timestamps[count])
			assertEqual(t, "data", r.Data, uids[count])

			count++
		}
		assertNilError(t, it.Error())

		assertEqual(t, "count", count, len(uids))
	})

	t.Run("filter blocks", func(t *testing.T) {
		// start from the third block after adding schemas
		it, err := client.EAS.FilterTimestamped(ctx, blockNumber+3, eas.Ptr(blockNumber+5), nil, nil)
		assertNilError(t, err)
		defer it.Close()

		count := 0

		for it.Next() {
			r := it.Value()

			index := count + 2

			timestamp, err := client.EAS.GetTimestamp(ctx, uids[index])
			assertNilError(t, err)

			assertEqual(t, "timestamp", r.Timestamp, timestamp)
			assertEqual(t, "timestamp", r.Timestamp, timestamps[index])
			assertEqual(t, "data", r.Data, uids[index])

			count++
		}
		assertNilError(t, it.Error())

		assertEqual(t, "count", count, 3)
	})

	t.Run("filter data", func(t *testing.T) {
		it, err := client.EAS.FilterTimestamped(ctx, 0, nil, []eas.UID{uids[1], uids[3]}, nil)
		assertNilError(t, err)
		defer it.Close()

		count := 0

		index := 1
		for it.Next() {
			r := it.Value()

			timestamp, err := client.EAS.GetTimestamp(ctx, uids[index])
			assertNilError(t, err)

			assertEqual(t, "timestamp", r.Timestamp, timestamp)
			assertEqual(t, "timestamp", r.Timestamp, timestamps[index])
			assertEqual(t, "data", r.Data, uids[index])

			count++

			index = 3
		}
		assertNilError(t, it.Error())

		assertEqual(t, "schema count", count, 2)
	})

	t.Run("filter timestamps", func(t *testing.T) {
		it, err := client.EAS.FilterTimestamped(ctx, 0, nil, nil, []eas.Timestamp{timestamps[3], timestamps[5]})
		assertNilError(t, err)
		defer it.Close()

		count := 0

		index := 3
		for it.Next() {
			r := it.Value()

			timestamp, err := client.EAS.GetTimestamp(ctx, uids[index])
			assertNilError(t, err)

			assertEqual(t, "timestamp", r.Timestamp, timestamp)
			assertEqual(t, "timestamp", r.Timestamp, timestamps[index])
			assertEqual(t, "data", r.Data, uids[index])

			count++

			index = 5
		}
		assertNilError(t, it.Error())

		assertEqual(t, "schema count", count, 2)
	})
}

func TestSchemaEAS_WatchTimestamped(t *testing.T) {
	client, backend := newClient(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sink := make(chan *eas.EASTimestamped)
	sub, err := client.EAS.WatchTimestamped(ctx, nil, sink, nil, nil)
	assertNilError(t, err)

	count := 0

	uids := []eas.UID{
		eas.HexDecodeUID("0x00cf55814f6a9135f3458c52cb8bc511e8f22b3e864f29325f52229bc12a4d15"),
		eas.HexDecodeUID("0x1605a2096b0bb3dd2ed53e268d0fd797271040a574dc1c995c3560357f9d7432"),
		eas.HexDecodeUID("0xcb8d8ccf0d47f2269ef59fcc16d6d4cdd64021b28c2edf41efce79dcb9c29a17"),
		eas.HexDecodeUID("0x3380b8d16d19cecce426d45e24cef6d1dc159a0f8d35808889f112cb2e003794"),
		eas.HexDecodeUID("0xa821c015ec279aa37500b66a4a65f3230cbeb5b0c707e4a9265771a67346bbe2"),
		eas.HexDecodeUID("0x204e9bc8300e1206dc9b41d096555ab70ec4423fbbb6a9b96a30691f30f673ba"),
		eas.HexDecodeUID("0x28ae72581c95e4664a63eaa5036a8f5e8515ed9f7c8a149058c9fcb3e8216da9"),
	}

	go func() {
		defer sub.Unsubscribe()

		for _, uid := range uids {
			_, wait, err := client.EAS.Timestamp(ctx, nil, uid)
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
		case r := <-sink:
			assertEqual(t, "uid", r.Data, uids[count])

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

	assertEqual(t, "count", count, 7)
}
