A small tool to make packaging for Solus even easier.

## Features

- Remove all packages from the local solbuild repo.
- Copy eopkg files to the default local repo and re-index the repo.
- Initialize a new package repo

## Usage

`./soltools subcommand`

Available subcommands:

- `clean`
- `copy`
- `init NAME SOURCE_ARCHIVE_URL`

## License

Copyright &copy; 2021 Evan Maddock maddock.evan@vivaldi.net

`soltools` is available under the terms of the Apache 2.0 license.
