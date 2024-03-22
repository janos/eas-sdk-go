# How to contribute

We'd love to accept your patches and contributions to this project. There are just a few small guidelines you need to follow.

1. Code should be `go fmt` formatted.
2. Exported types, constants, variables and functions should be documented.
3. Changes must be covered with tests.
4. All tests must pass constantly `go test .`.

## Versioning

Ethereum Attestation Service Go client follows semantic versioning. New functionality should be accompanied by increment to the minor version number.

## Releasing

Any code which is complete, tested, reviewed, and merged to master can be released.

1. Make a pull request with changes.
2. Once the pull request has been merged, visit [https://github.com/janos/eas-sdk-go/releases](https://github.com/janos/eas-sdk-go/release) and click `Draft a new release`.
3. Update the `Tag version` and `Release title` field with the new Ethereum Attestation Service Go client version. Be sure the version has a `v` prefixed in both places, e.g. `v1.25.0`.
4. Publish the release.
