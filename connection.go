// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import "os"

// The most basic type of database connection.
//
// The choice to separate Prepare() and Execute() is deliberate:
// It leaves the database driver the most flexibilty for achieving
// good performance without requiring additional caching schemes.
//
// Prepare() accepts a query language string and returns
// a precompiled statement that can be executed after any
// remaining parameters have been bound. The format of
// parameters in the query string is dependent on the
// database driver in question.
//
// Execute() accepts a precompiled statement, binds the
// given parameters, and then executes the statement.
// Execute() returns a channel of Result objects which
// can be examined one at a time (if the query produced
// results to begin with). Specific database drivers
// will return result objects conforming to one or more
// of the following interfaces which represent different
// levels of functionality.
//
// Close() ends the connection to the database system
// and frees up all internal resources associated with
// it. Note that you must close all objects created on
// the connection before closing the connection itself.
// After a connection has been closed, no further
// operations are allowed on it.
type Connection interface {
	Prepare(query string) (Statement, os.Error);
	Execute(stat Statement, parameters ...) (<-chan Result, os.Error);
	Close() os.Error;
}

// InformativeConnections supply useful but optional information.
//
// Changes() returns the number of changes the last query made
// to the database. Note that the database driver has to explain
// what exactly constitutes a "change" for a given database system
// and query.
//
// LastId() returns the id of the last successful insertion into
// the database. The database driver has to explain the exact
// meaning of the id and the conditions under which it changes.
type InformativeConnection interface {
	Connection;
	Changes() (int, os.Error);
	LastId() (int, os.Error);
}

// TransactionalConnections support transactions. Note that
// the database driver in question may be in "auto commit"
// mode by default. Once you call Begin(), "auto commit" will
// be disabled for that connection until you either Commit()
// or Rollback() successfully.
//
// Begin() starts a transaction.
//
// Commit() tries to push all changes made as part of the
// current transaction to the database.
//
// Rollback() tries to undo all changes made as part of the
// current transaction.
type TransactionalConnection interface {
	Connection;
	Begin() os.Error;
	Commit() os.Error;
	Rollback() os.Error;
}
