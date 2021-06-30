A small tool to make packaging for Solus even easier.

[![Go Report Card](https://goreportcard.com/badge/github.com/EbonJaeger/soltools)](https://goreportcard.com/report/github.com/EbonJaeger/soltools) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Features

- Remove all packages from the local solbuild repo.
- Copy eopkg files to the default local repo and re-index the repo.
- Initialize a new package repo
- Clone an existing package repo

## Usage

`./soltools subcommand`

Available subcommands:

- `clean`
- `copy`
- `init NAME SOURCE_ARCHIVE_URL`
- `clone NAME`

## License

Copyright &copy; 2021 Evan Maddock maddock.evan@vivaldi.net

`soltools` is available under the terms of the Apache 2.0 license.
