# Ethereum Attestation Service - Go SDK

[![Go](https://github.com/janos/eas/workflows/Go/badge.svg)](https://github.com/janos/eas/actions)
[![PkgGoDev](https://pkg.go.dev/badge/resenje.org/eas)](https://pkg.go.dev/resenje.org/eas)
[![NewReleases](https://newreleases.io/badge.svg)](https://newreleases.io/github/janos/eas)

This repository contains the Ethereum Attestation Service SDK for the Go programming language, used to interact with the Ethereum Attestation Service Protocol.

## Example

```go
package main

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"resenje.org/eas"
)

func main() {
	const (
		// network
		endpointSepolia        = "https://ethereum-sepolia-rpc.publicnode.com/"
		contractAddressSepolia = "0xC2679fBD37d54388Ce493F1DB75320D236e1815e"

		// account
		privateKeyHex = "933c798b990a6be3fb91ae2fd3b6593f61d6d478548091205ee948b1de9c9f19" // this is not a real user key
		)
	)

	privateKey, err := eas.HexParsePrivateKey(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(contractAddressSepolia)

	c, err := eas.NewClient(context.Background(), endpointSepolia, privateKey, contractAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

	attestationUID := eas.HexDecodeUID("0xac812932f5cee90a457d57a9fbd7b142b21ba99b809f982bbf86947f295281ff")

	a, err := c.EAS.GetAttestation(context.Background(), attestationUID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Attester", a.Attester, a.Time, a.Schema)

	var schemaUID eas.UID
	var name string

	if err := a.ScanValues(&schemaUID, &name); err != nil {
		log.Fatal(err)
	}

	log.Println("Attestation")
	log.Println("Schema UID:", schemaUID)
	log.Println("Name:", name)
}
```

## Versioning

Each version of the client is tagged and the version is updated accordingly.

This package uses Go modules.

To see the list of past versions, run `git tag`.

## Contributing

We love pull requests! Please see the [contribution guidelines](CONTRIBUTING.md).

## License

This library is distributed under the BSD-style license found in the [LICENSE](LICENSE) file.
