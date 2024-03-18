// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas_test

import (
	"context"
	"testing"
)

func TestEASContract_Version(t *testing.T) {
	client, _ := newClient(t)
	ctx := context.Background()

	version, err := client.EAS.Version(ctx)
	assertNilError(t, err)

	assertEqual(t, "version", version, "1.0.0")
}
