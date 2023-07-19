# Change log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.0] - 2023-07-20

### Added

- Added the possibility to authenticate against Hydra Admin API using custom `Authorization` header.

- Added back/front-channel logout arguments. Thanks to @agouil.

### Changed

- The provider is now built with Go 1.20.

## [0.4.0] - 2022-05-13

### Added

- macOS arm binaries are now available.

### Changed

- The provider is now built with Go 1.16.

## [0.3.0] - 2020-10-25

### Added

- Added the possibility to authenticate against Hydra Admin API using OAuth2.0 client credentials.

## [0.2.0] - 2020-10-01

This release is mainly about publishing the provider to the Terraform Registry and improving documentation.

### Added

- The provider is now available in the [Terraform Registry](https://registry.terraform.io/providers/hypnoglow/oryhydra/latest).

## [0.1.1] - 2020-06-10

### Fixed

- Fixed handling non-existing OAuth2 Clients in the provider.

## [0.1.0] - 2020-06-08

Initial release.
