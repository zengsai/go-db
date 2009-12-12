// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Database API for Go.
//
// Terminology:
//
// Database systems are pieces of software (usually outside of Go)
// that allow storage and retrieval of data. We try not to imply
// "relational" at the level of this API.
//
// Database drivers are pieces of software (usually written in Go)
// that allow Go programs to interact with database systems through
// some query language. We try not to imply "SQL" at the level of
// this API.
//
// Goals:
//
// The API described here is a set of conventions that should be
// followed by database drivers. Obviously there are levels of
// compliance, but every database driver should at least implement
// the core of the API: the functions Version() and Open() as well
// as the interfaces Connection, Statement, and Result.
//
// Errors:
//
// The API uses os.Error to indicate errors as is customary in Go.
// However, it can be useful for clients to be able to distinguish
// errors reported by the database driver from errors reported by
// the database system. We therefore encourage drivers to implement
// at least two error types, DriverError and SystemError. Clients
// can then check the runtime type of an error if they wish to.
//
// Classic API:
//
// The "classic" API is completely optional and not all database
// drivers support it. It's sole purpose is to provide a faster
// way of accessing results while the Go runtime is catching up
// and the speed of channels is being improved. The "classic" API
// will disappear eventually, so it's better to stay away from it.
package db

// TODO: should we do this?
// type DriverError interface {
// 	os.Error;
// 	Driver() string;
// }
//
// type SystemError interface {
// 	os.Error;
// 	System() string;
// }
