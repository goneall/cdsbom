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

This will read `input-sbom.json` and query ClearlyDefined for License
information. The License fields in the SBOM will be replaced to use the license
data returned from ClearlyDefined. A new sbom will be written to
`enhanced-sbom.json` with the updated fields in the same format as the input
sbom.

Supported formats are the [same as
Protobom](https://github.com/protobom/protobom/blob/main/README.md#supported-versions-and-formats).

## sbomnotice

### Install:

```sh
go install github.com/jeffmendoza/cdsbom/sbomnotice@latest
```

Make sure `$GOBIN` is in your path.

- `$GOBIN` defaults to `$GOPATH/bin`
- `$GOPATH` defaults to `$HOME/go` on Unix and `%USERPROFILE%\go` on Windows

### Use:

Example:
```sh
sbomnotice -out NOTICE input-sbom.json
```

This will read `input-sbom.json` and parse all the dependencies found. Then
query ClearlyDefined for a NOTICE file with all dependencies from the SBOM. The
file contents will be written to the provided output file name, or `NOTICE` if
not specified.

Supported formats are the [same as
Protobom](https://github.com/protobom/protobom/blob/main/README.md#supported-versions-and-formats).

## Thanks

This project is possible due to
[Protobom](https://github.com/protobom/protobom) for SBOM parsing, and [GUAC
sw-id-core](https://github.com/guacsec/sw-id-core) to convert PURL to
ClearlyDefined Coordinates.
