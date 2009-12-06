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
// Database bindings are pieces of software (usually written in
// Go) that allow Go programs to interact with database systems
// through some query language. We try not to imply "SQL" at the
// level of this API.
//
// Goals:
//
// The API described here is a set of conventions that should be
// followed by database bindings. Obviously there are levels of
// compliance, but every database binding should at least implement
// the core of the API: the functions Version() and Open() as well
// as the interfaces Connection, Statement, and Cursor.
package db

import "os"

// Each database binding must provide a Version() function to
// allow careful clients to configure themselves appropriately
// for the database system in question. There are a number of
// well-known keys in the map returned by Version():
//
//	Key		Description
//
//	version		generic version (if client/server doesn't apply)
//	client		client version
//	server		server version
//	protocol	protocol version
//	binding		database binding version
//
// Database bindings decide which of these keys to return. For
// example, db/sqlite3 returns "version" and "binding"; db/mysql
// should probably return all keys except "version" instead.
//
// Database bindings can also return additional keys, provided
// they prefix them with the package name of the binding in
// question. The db/sqlite3 binding, for example, returns
// "sqlite3.sourceid" as well.
type VersionSignature func() (map[string]string, os.Error)

// Each database binding must provide an Open() function to
// establish connections to a database system. Database systems
// require a wide variety of parameters for connections, which
// is why the parameters to Open() are passed as a map.
//
// XXX: THE MAP WILL BE REPLACED WITH SOME FORM OF URL IN THE
// NEAR FUTURE. http://golang.org/pkg/http/#URL
//
// Each map entry consists of a string key and a generic value.
// There are a number of well-known keys that apply to many (if
// not all) database systems:
//
//	Name		Type	Description
//
//	name		string	the database to connect to
//	host		string	the host to connect to
//	port		int	the port to connect to
//	username	string	the user to connect as
//	password	string	the password for that user
//
// For example, the following piece of code tries to connect to
// a MySQL database on the local machine at the default port:
//
//	c, e := mysql.Open(
//		Arguments{
//			"name": "mydb",
//			"username": "phf",
//			"password": "somepassword"
//		}
//	)
//
// Note that defaults for all keys are specific to the database
// binding in question and should be documented there.
//
// The Open() function is free to ignore entries that it has no
// use for. For example, the sqlite3 binding only understands
// "name" and ignores the other well-known keys.
//
// A database binding is free to introduce additional keys if
// necessary, however those keys have to start with the package
// name of the database binding in question. For example, the
// sqlite3 binding supports the key "sqlite3.vfs".
//
// A successful call to Open() results in a connection to the
// database system. Specific database binding will return
// connection objects conforming to one or more of the following
// interfaces which represent different levels of functionality.
type OpenSignature func(args map[string]interface{}) (conn Connection, err os.Error)

// The most basic type of database connection.
//
// The choice to separate Prepare() and Execute() is deliberate:
// It leaves the database binding the most flexibilty for achieving
// good performance without requiring additional caching schemes.
//
// Prepare() accepts a query language string and returns
// a precompiled statement that can be executed after any
// remaining parameters have been bound. The format of
// parameters in the query string is dependent on the
// database binding in question.
//
// Execute() accepts a precompiled statement, binds the
// given parameters, and then executes the statement.
// If the statement produces results, Execute() returns
// a cursor; otherwise it returns nil. Specific database
// bindings will return cursor objects conforming to one
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

// The iterator approach to execute returns these things
type Result struct {
	Data	[]interface{};
	Error	os.Error;
}

// InformativeConnections supply useful but optional information.
//
// Changes() returns the number of changes the last query made
// to the database. Note that the database binding has to explain
// what exactly constitutes a "change" for a given database system
// and query.
type InformativeConnection interface {
	Connection;
	Changes() (int, os.Error);
}

// TransactionalConnections support transactions. Note that
// the database binding in question may be in "auto commit"
// mode by default. Once you call Begin(), "auto commit" will
// be disabled for that connection.
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

// The most basic type of database cursor.
// TODO: base on exp/iterable instead? Iter() <-chan interface{};
//
// MoreResults() returns true if there are more results
// to be fetched.
//
// FetchOne() returns the next result from the database.
// Each result is returned as an array of generic objects.
// The database binding in question has to define what
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
// names is specified by the database binding in question.
//
// Results() returns the number of results remaining to be
// fetched.
type InformativeCursor interface {
	Cursor;
	Description() (map[string]string, os.Error);
	Results() int;
}

// PythonicCursors fetch results as maps from field names to
// values instead of just slices of values.
//
// TODO: find a better name for this!
//
// FetchDict() is similar to FetchOne().
// FetchDictMany() is similar to FetchMany().
// FetchDictAll() is similar to FetchAll().
type PythonicCursor interface {
	Cursor;
	FetchDict() (data map[string]interface{}, error os.Error);
	FetchManyDicts(count int) (data []map[string]interface{}, error os.Error);
	FetchAllDicts() (data []map[string]interface{}, error os.Error);
}

// ExecuteDirectly is a convenience function for "one-off" queries.
// If you need more control, for example to rebind parameters over
// and over again, or to get results one by one, you should use the
// Prepare() and Execute() methods explicitly instead.
func ExecuteDirectly(conn Connection, query string, params ...) (results [][]interface{}, err os.Error) {
	var s Statement;
	s, err = conn.Prepare(query);
	defer s.Close();
	if err != nil {
		return
	}

	var c Cursor;
	c, err = conn.Execute(s, params);
	if err != nil || c == nil {
		return
	}
	defer c.Close();

	results, err = c.FetchAll();
	return;
}
