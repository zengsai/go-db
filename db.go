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
package db

import "os"
import "strings"

// Database drivers must provide the Version() function to allow
// careful clients to configure themselves appropriately for the
// database system in question. There are a number of well-known
// keys in the map returned by Version():
//
//	Key		Description
//
//	version		generic version (if client/server doesn't apply)
//	client		client version
//	server		server version
//	protocol	protocol version
//	driver		database driver version
//
// Database drivers decide which of these keys to return. For
// example, the sqlite3 driver returns "version" and "driver";
// the mysql driver should probably return all keys except
// "version" instead.
//
// Database drivers can also return additional keys provided
// they prefix them with the package name of the driver. The
// sqlite3 driver, for example, returns "sqlite3.sourceid" in
// addition to "version" and "driver".
type VersionSignature func() (map[string]string, os.Error)

// Database drivers must provide the Open() function to allow
// clients to establish connections to a database system. The
// parameter to Open() is a URL of the following form:
//
//	driver://username:password@host:port/database?key=value;key=value
//
// Most parts of this URL are optional. The sqlite3 database
// driver for example interprets "sqlite3://test.db" as the
// database "test.db" in the current directory. Actually, it
// also interprets "test.db" by itself that way. If a driver
// is specified in the URL, it has to match the driver whose
// Open() function is called. For example the sqlite3 driver
// will fail if asked to open "mysql://somedb". There can be
// as many key/value pairs as necessary to configure special
// features of the particular database driver. Here are more
// examples:
//
//	c, e := mysql.Open("mysql://phf:wow@example.com:7311/mydb");
//	c, e := sqlite3.Open("test.db?flags=0x00020001");
//
// Note that defaults for all optional components are specific
// to the database driver in question and should be documented
// there.
//
// The Open() function is free to ignore components that it
// has no use for. For example, the sqlite3 driver ignores
// username, password, host, and port.
//
// A successful call to Open() results in a connection to the
// database system. Specific database drivers will return
// connection objects conforming to one or more of the following
// interfaces which represent different levels of functionality.
type OpenSignature func(url string) (conn Connection, err os.Error)

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
// If the statement produces results, Execute() returns
// a cursor; otherwise it returns nil. Specific database
// driver will return cursor objects conforming to one
// or more of the following interfaces which represent
// different levels of functionality.
//
// Iterate() is an experimental variant of Execute()
// that returns a channel of Result objects instead
// of a Cursor. XXX: Is this any good?
//
// Close() ends the connection to the database system
// and frees up all internal resources associated with
// it. Note that you must close all Statement and Cursor
// objects created through a connection before closing
// the connection itself. After a connection has been
// closed, no further operations are allowed on it.
type Connection interface {
	Prepare(query string) (Statement, os.Error);
	Execute(statement Statement, parameters ...) (Cursor, os.Error);
	Iterate(statement Statement, parameters ...) (<-chan Result, os.Error);
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

// Statements are precompiled queries, possibly with remaining
// parameter slots that need to be filled before execution.
// TODO: include parameter binding API? or subsume in Execute()?
// what about resetting the statement or clearing parameter
// bindings?
type Statement interface {
	Close() os.Error;
}

// The most basic type of result.
//
// Data() returns the data for this result as an array
// of generic objects. The database driver in question
// defines what concrete types are returned depending
// on the types used by the database system.
//
// Error() returns the error that occurred when this
// result was fetched, or nil if no error occurred.
type Result interface {
	Data() []interface{};
	Error() os.Error;
}

// InformativeResults supply useful but optional information.
//
// Fields() returns the names of each item of data in the
// result.
//
// Types() returns the names of the types of each item in
// the result.
type InformativeResult interface {
	Result;
	Fields() []string;
	Types() []string;
}

// FancyResults provide an alternate way of processing results.
//
// DataMap() returns a map from item names to item values. As
// for Data() the concrete types have to be defined by the
// database driver in question.
//
// TypeMap() returns a map from item names to the names of the
// types of each item.
type FancyResult interface {
	Result;
	DataMap() map[string]interface{};
	TypeMap() map[string]string;
}

// The most basic type of database cursor.
// TODO: base on exp/iterable instead? Iter() <-chan interface{};
//
// MoreResults() returns true if there are more results
// to be fetched.
//
// FetchOne() returns the next result from the database.
// Each result is returned as an array of generic objects.
// The database driver in question has to define what
// concrete types are returned depending on the types
// used by the database system.
//
// FetchMany() returns at most count results.
// XXX: FetchMany() MAY GO AWAY SOON.
//
// FetchAll() returns all (remaining) results.
// XXX: FetchAll() MAY GO AWAY SOON.
//
// Close() frees the cursor. After a cursor has been
// closed, no further operations are allowed on it.
type Cursor interface {
	MoreResults() bool;
	FetchOne() ([]interface{}, os.Error);
	FetchMany(count int) ([][]interface{}, os.Error);
	FetchAll() ([][]interface{}, os.Error);
	Close() os.Error;
}

// InformativeCursors supply useful but optional information.
//
// Description() returns a map from (the name of) a field to
// (the name of) its type. The exact format of field and type
// names is specified by the database driver in question.
//
// Results() returns the number of results remaining to be
// fetched.
type InformativeCursor interface {
	Cursor;
	Description() (map[string]string, os.Error);
	Results() int;
}

// XXX: an experimental "classic" API for results (to replace Cursor)
type ResultSet interface {
	More() bool;
	Fetch() Result;
}

// ExecuteDirectly is a convenience function for "one-off" queries.
// It's particularly convenient for queries that don't produce any
// results.
//
// If you need more control, for example to rebind parameters over
// and over again, to get results one by one, or to access metadata
// about the results, you should use the Prepare() and Execute()
// methods explicitly instead.
func ExecuteDirectly(conn Connection, query string, params ...) (results [][]interface{}, err os.Error) {
	var s Statement;
	s, err = conn.Prepare(query);
	if err != nil || s == nil {
		return
	}
	defer s.Close();

	var c Cursor;
	c, err = conn.Execute(s, params);
	if err != nil || c == nil {
		return
	}
	defer c.Close();

	results, err = c.FetchAll();
	return;
}

// ParseQueryURL() helps database drivers parse URLs passed to
// Open(). It takes a string of the form
//
//	key=value{;key=value;...;key=value}]
//
// and returns a map from keys to values. Empty strings yield
// an empty map; malformed strings yield a nil map instead; a
// string is malformed, for example, if it contains duplicate
// keys.
func ParseQueryURL(str string) (opt map[string]string) {
	opt = make(map[string]string);
	if len(str) == 0 {
		return opt
	}

	pairs := strings.Split(str, ";", 0);
	if len(pairs) == 0 {
		return nil
	}

	for _, p := range pairs {
		pieces := strings.Split(p, "=", 0);
		if len(pieces) == 2 {
			if _, duplicate := opt[pieces[0]]; duplicate {
				return nil;
			}
			opt[pieces[0]] = pieces[1]
		} else {
			return nil
		}
	}

	return;
}
