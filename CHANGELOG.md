# Change Log


## Unreleased

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
