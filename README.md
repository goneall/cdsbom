# ClearlyDefined SBOM tools

This repository contains tools for using SBOMs with [ClearlyDefined](https://docs.clearlydefined.io/)

## cdsbom

### Install:

```sh
go install github.com/jeffmendoza/cdsbom@latest
```

Make sure `$GOBIN` is in your path.

- `$GOBIN` defaults to `$GOPATH/bin`
- `$GOPATH` defaults to `$HOME/go` on Unix and `%USERPROFILE%\go` on Windows

### Use:

Example:
```sh
cdsbom -out enhanced-sbom.json input-sbom.json
```

This will read `input-sbom.json` and query ClearlyDefined for Licnese
information. The License fields in the SBOM will be replaced to use the license
data returned from ClearlyDefined. A new sbom will be written to
`enhanced-sbom.json` with the updated fields in the same format as the input
sbom.

Supported formats are the [same as
Protobom](https://github.com/protobom/protobom/blob/main/README.md#supported-versions-and-formats).

## Future / TODO

Another tool to generate a NOTICE file from an SBOM using ClealyDefined.

## Thanks

This project is possible due to
[Protobom](https://github.com/protobom/protobom) for SBOM parsing, and [GUAC
sw-id-core](https://github.com/guacsec/sw-id-core) to convert PURL to
ClearlyDefined Coordinates.
