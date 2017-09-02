# Change Log


## Unreleased

### Added

- Support for virtual addresses and servers without network connectivity
- Functional options to constructors
- App server implementing all functionalities from `named` and `aio` servers

### Changed

- Queue accepts an address as an argument

### Removed

- `named` and `aio` servers
- Direct [go-kit](https://github.com/go-kit/kit) dependency
- Manager is no longer necessary (most of it's functionality is moved to `Queue` and `AppServer`)


## 0.11.0 - 2017-08-11

### Added

- `Close` call to server queue

### Changed

- Renamed `Stop` to `Shutdown` in server queue

### Removed

- Manager parameter from Queue constructor


## 0.10.0 - 2017-08-07

### Added

- Manager now accepts a network in `ListenAndStartServer`
- `AddrServer` to represent a server which holds network information

### Changed

- Moved mocks to internal package
- Renamed `Name` method to `GetName` (and `ServerName` fields to `Name`)
- Server queue only accepts `AddrServer`


### Removed

- Unused logger interface
- Unused constructors
- gRPC dependency from glide


## 0.9.0 - 2017-07-05

### Added

- Testify library and `mocks.Server` server mock
- `aio.Server` including all additional server features

### Changed

- Moved `NamedServer` to `named.Server`


## 0.8.3 - 2017-06-28

### Changed

- Use go kit log package for levels and nop implementation


## 0.8.2 - 2017-06-24

### Fixed

- Error return type


## 0.8.1 - 2017-06-24

### Fixed

- WaitGroup add


## 0.8.0 - 2017-06-24

### Added

- Multi error returned from `Queue.StopServer`

### Changed

- Use Go-kit's log interface
- Panics are not handled by the error handler anymore
- Manager stop function returns an error

### Removed

- Error handler dependency


## 0.7.0 - 2016-06-19

### Added

- Constructor for Queue

### Changed

- Renamed `ServerQueue` to `Queue`


## 0.6.0 - 2016-06-18

### Changed

- Moved back to `_test` test packages
- Rewritten logger interface


## 0.5.0 - 2017-06-12

## Added

- `ServerQueue` to fully manage servers

## Changed

- Moved `GrpcServer` to `grpc.Server`


## 0.4.0 - 2017-06-08

### Fixed

- [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) dependency


## 0.3.0 - 2017-05-16

### Removed

- Shutdown manager (use [sagikazarmark/utilz](https://github.com/sagikazarmark/utilz) instead)


## 0.2.0 - 2017-03-02

### Added

- gRPC server


## 0.1.3 - 2017-02-28

### Changed

- Renamed "name" field to "server"


## 0.1.2 - 2017-02-28

### Changed

- Improved logging


## 0.1.1 - 2017-02-28

### Changed

- Improved logging
- Improved tests


## 0.1.0 - 2017-02-28

### Added

- Shutdown handler
- Server manager
