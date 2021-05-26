# sctool

`sctool` is a CLI tool that provides various utilities for working with SC files.

## Installing

Currently, `sctool` can only be installed from source. This requires Go.

```
go install github.com/sc-lang/sctool
```

## Usage

`sctool` contains the following tools:

- formatter
- validator

For details on each command and their usage run `sctool <command> --help`.

Here are some quick examples:

To format an SC file:

```
sctool fmt input.sc
```

To validate an SC file:

```
sctool validate input.sc
```

For details on commands and their usage
