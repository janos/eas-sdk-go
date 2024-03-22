# Ethereum Attestation Service - Go SDK

[![Go](https://github.com/janos/eas-sdk-go/workflows/Go/badge.svg)](https://github.com/janos/eas-sdk-go/actions)
[![PkgGoDev](https://pkg.go.dev/badge/resenje.org/eas)](https://pkg.go.dev/resenje.org/eas)
[![NewReleases](https://newreleases.io/badge.svg)](https://newreleases.io/github/janos/eas-sdk-go)

This repository contains the Ethereum Attestation Service SDK for the Go programming language, used to interact with the Ethereum Attestation Service Protocol.

[Ethereum Attestation Service](https://attest.sh/) (EAS) is an open-source infrastructure public good for making attestations onchain or offchain.

Go SDK interacts with [EAS Smart Contracts](https://github.com/ethereum-attestation-service/eas-contracts) deployed on different EVM-compatible blockchains using smart contract bindings. The list off deployed contracts could be found in [EAS Contracts README file](https://github.com/ethereum-attestation-service/eas-contracts?tab=readme-ov-file#deployments). EAS contract address for a desired network in that list should be passed as an argument to the client constructor `eas.NewCLient`.

## Installing the Go EAS SDK

Run `go get resenje.org/eas` from command line in your Go module directory.

## Usage

Please refer to the generated package documentation on <https://pkg.go.dev/resenje.org/eas>, examples bellow and, of course, tests and code in this repository as the last resource of open source projects.

## Schemas

Attestations are structured by defining and registering Schemas. Schemas follow the Solidity ABI for acceptable types. Below is a list of current Solidity types and corresponding Go types.

| Go type                                    | Solidity type                               |
|--------------------------------------------|---------------------------------------------|
| `common.Address`                           | `address`                                   |
| `string`                                   | `string`                                    |
| `bool`                                     | `bool`                                      |
| `[32]byte`                                 | `bytes32`                                   |
| `eas.UID`                                  | `bytes32`                                   |
| `[]byte`                                   | `bytes`                                     |
| `uint8`                                    | `uint8`                                     |
| `uint16`                                   | `uint16`                                    |
| `uint32`                                   | `uint32`                                    |
| `uint64`                                   | `uint64`                                    |
| `uint256`                                  | `uint256`                                   |
| `struct{<name> <type>; <name> <type>;...}` | `tuple = (<type> <name>, <type> <name>...)` |
| `slice = []<type>`                         | `<type>[]`                                  |
| `array = [<size>]<type>`                   | `<type>[<size>]`                            |

All supported types can be nested inside tuples (Go structs), fixed-sized arrays (Go arrays) and variable-length arrays (So slices).

### Field names

Solidity tuples are represented with Go struct type where names of the struct fields are used for tuple field names. It is possible to set a custom name with Go struct field tag `abi`. En example of a tuple related type:

```go
type MyTuple struct {
	ID        eas.UID `abi:"id"`
	Msg       string  `abi:"message"`
	Timestamp uint64  `abi:"timeStamp"`
	RawData   []byte  `abi:"raw_data"`
	Sender    common.Address
}
```

which corresponds to this schema definition:

```solidity
bytes32 id, string message, uint64 timeStamp, bytes raw_data, address Sender
```

## Examples

### Get an existing attestation

```go
package main

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"resenje.org/eas"
)

var (
	endpointSepolia        = "https://ethereum-sepolia-rpc.publicnode.com/"
	contractAddressSepolia = common.HexToAddress("0xC2679fBD37d54388Ce493F1DB75320D236e1815e")
)

func main() {
	ctx := context.Background()

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	c, err := eas.NewClient(ctx, endpointSepolia, privateKey, contractAddressSepolia, nil)
	if err != nil {
		log.Fatal(err)
	}

	attestationUID := eas.HexDecodeUID("0xac812932f5cee90a457d57a9fbd7b142b21ba99b809f982bbf86947f295281ff")

	a, err := c.EAS.GetAttestation(ctx, attestationUID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Attestation UID", a.UID)
	log.Println("Attestation Time", a.Time)

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

### Structured schema

Create a structured schema, make attestation and get attestation.

```go
// Attest a road trip by defining a schema.
package main

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"resenje.org/eas"
)

var (
	endpointSepolia        = "https://ethereum-sepolia-rpc.publicnode.com/"
	contractAddressSepolia = common.HexToAddress("0xC2679fBD37d54388Ce493F1DB75320D236e1815e")
)

type RoadTrip struct {
	ID           uint64      `abi:"id"`
	VIN          string      `abi:"vin"` // Vehicle Identification Number
	VehicleOwner string      `abi:"vehicleOwner"`
	Passengers   []Passenger `abi:"passengers"`
}

type Passenger struct {
	Name     string `abi:"name"`
	CanDrive bool   `abi:"canDrive"`
}

func (p Passenger) CanDriveString() string {
	if p.CanDrive {
		return "can drive"
	}
	return "cannot drive"
}

// Attest a road trip by defining a schema.
func main() {
	ctx := context.Background()

	// Use a fake key here. Use your own funded key to be able to send transactions.
	privateKey, err := eas.HexParsePrivateKey("a896e1f28a6453e8db4794f11ea185befd04c4e4f06790e37e8d1cc90a611948")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Wallet address:", crypto.PubkeyToAddress(privateKey.PublicKey))

	// Construct a client that will interact with EAS contracts.
	c, err := eas.NewClient(ctx, endpointSepolia, privateKey, contractAddressSepolia, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create the Schema on chain.
	tx, waitRegistration, err := c.SchemaRegistry.Register(ctx, eas.MustNewSchema(RoadTrip{}), common.Address{}, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Waiting schema registration transaction:", tx.Hash())
	schemaRegistration, err := waitRegistration(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Just check the schema definition.
	schema, err := c.SchemaRegistry.GetSchema(ctx, schemaRegistration.UID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Schema UID:", schema.UID)
	log.Println("Schema:", schema.Schema)

	// Attest a road trip on chain.
	tx, waitAttestation, err := c.EAS.Attest(ctx,
		schema.UID,
		&eas.AttestOptions{Revocable: true},
		RoadTrip{
			ID:           4242,
			VIN:          "1FA6P8CF5L5100421",
			VehicleOwner: "Richard Hammond",
			Passengers: []Passenger{
				{
					Name:     "James May",
					CanDrive: true,
				},
				{
					Name:     "Jeremy Clarkson",
					CanDrive: true,
				},
				{
					Name:     "The Stig",
					CanDrive: false,
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Waiting attest transaction:", tx.Hash())
	attestConfirmation, err := waitAttestation(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Get the attestation to verify it.
	a, err := c.EAS.GetAttestation(ctx, attestConfirmation.UID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Attestation UID", a.UID)
	log.Println("Attestation Time", a.Time)

	var roadTrip RoadTrip
	if err := a.ScanValues(&roadTrip); err != nil {
		log.Fatal(err)
	}

	log.Println("Road trip:", roadTrip.ID)
	log.Println("Vehicle Identification Number:", roadTrip.VIN)
	log.Println("Vehicle owner:", roadTrip.VehicleOwner)
	for i, p := range roadTrip.Passengers {
		log.Printf("Passenger %v: %s (%s)", i, p.Name, p.CanDriveString())
	}
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
