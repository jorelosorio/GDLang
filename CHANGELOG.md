# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- The AST tree and the compiler can now handle `unsigned integers of 16 bits` as names for variables and functions. This change will decrease binary size and increase performance in the stack.

### Changed

- A refactor was made in the AST tree to improve the performance of the compiler.
- The stack map was updated to handle `any` type as a key instead of a `GDIdent` type.
- Tests are now performed twice to test for `uint16` and `string` based variables and function names.

## [0.0.1-alpha] - 2024-09-22

### Added

- Initial release of the project