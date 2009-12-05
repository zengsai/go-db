/*
	BRAINSTORMING ONLY! DON'T RELY ON THIS YET!

	Terminology:

	Database systems are pieces of software (usually outside of Go)
	that allow storage and retrieval of data. We try not to imply
	"relational" at the level of this API.

	Database bindings are pieces of software (usually written in
	Go) that allow Go programs to interact with database systems
	through some query language. We try not to imply "SQL" at the
	level of this API.

	Goals:

	The API described here is a set of conventions that should be
	followed by database bindings. Obviously there are levels of
	compliance, but every database binding should at least implement
	the core of the API: the functions Version() and Open() as well
	as the interfaces Connection, Statement, and Cursor.
*/

package db

import "os"

/*
	Each database binding must provide a Version() function to
	allow careful clients to configure themselves appropriately
	for the database system in question. There are a number of
	well-known keys in the map returned by Version():

	Key		Description

	version		generic version (if client/server doesn't apply)
	client		client version
	server		server version
	protocol	protocol version
	binding		database binding version

	Database bindings decide which of these keys to return. For
	example, sqlite3 returns "version" and "binding"; mysql
	should probably return all keys except "version" instead.

	Database bindings can also return additional keys, provided
	they prefix them with the package name of the binding in
	question. The sqlite3 binding, for example, returns
	"sqlite3.sourceid" as well.
*/
type VersionSignature func () (map[string]string, os.Error)

/*
	Each database binding must provide an Open() function to
	establish connections to a database system. Database systems
	require a wide variety of parameters for connections, which
	is why the parameters to Open() are passed as a map.

	TODO: use map[string]string instead? may be friendlier if we
	are sure we never need to pass anything complicated; or pass
	a URI instead?

	Each map entry consists of a string key and a generic value.
	There are a number of well-known keys that apply to many (if
	not all) database systems:

	Name		Type	Description

	name		string	the database to connect to
	host		string	the host to connect to
	port		int	the port to connect to
	username	string	the user to connect as
	password	string	the password for that user

	For example, the following piece of code tries to connect to
	a MySQL database on the local machine at the default port:

	c, e := mysql.Open(Arguments{
		"name": "mydb",
		"username": "phf",
		"password": "somepassword"}
	)

	Note that defaults for all keys are specific to the database
	binding in question and should be documented there.

	The Open() function is free to ignore entries that it has no
	use for. For example, the sqlite3 binding only understands
	"name" and ignores the other well-known keys.

	A database binding is free to introduce additional keys if
	necessary, however those keys have to start with the package
	name of the database binding in question. For example, the
	sqlite3 binding supports the key "sqlite3.vfs".
*/
type OpenSignature func (args map[string]interface{}) (connection Connection, error os.Error)

/*
	A successful call to Open() results in a connection to the
	database system. Specific database binding will return
	connection objects conforming to one or more of the following
	interfaces which represent different levels of functionality.

	Note that the choice to separate Prepare() and Execute() for
	the most basic connection interface is deliberate: It leaves
	the database binding the most flexibilty in achieving good
	performance without requiring it to implement additional
	caching schemes.
*/
type Connection interface {
	/*
		Prepare() accepts a query language string and returns
		a precompiled statement that can be executed after any
		remaining parameters have been bound. The format of
		parameters in the query string is dependent on the
		database binding in question.
	*/
	Prepare(query string) (Statement, os.Error);
	/*
		Execute() accepts a precompiled statement, binds the
		given parameters, and then executes the statement.
		If the statement produces results, Execute() returns
		a cursor; otherwise it returns nil.
	*/
	Execute(statement Statement, parameters ...) (Cursor, os.Error);
	/*
		Close() ends the connection to the database system
		and frees up all internal resources associated with
		it. Note that you must close all Statement and Cursor
		objects created through a connection before closing
		the connection itself. After a connection has been
		closed, no further operations are allowed on it.
	*/
	Close() os.Error
}

/*
	InformativeConnections supply useful but optional information.
	TODO: more operations?
*/
type InformativeConnection interface {
	Connection;
	/*
		If a query modified the database, Changes() returns the number
		of changes that took place. Note that the database binding
		has to explain what exactly constitutes a change for a given
		database system and query.
	*/
	Changes() (int, os.Error);
}

/*
	FancyConnections support additional convenience operations.
	TODO: more operations?
*/
type FancyConnection interface {
	Connection;
	/*
		ExecuteDirectly() is a wrapper around Prepare() and Execute().
	*/
	ExecuteDirectly(query string, parameters ...) (*Cursor, os.Error)
}

/*
	TransactionalConnections support transactions. Note that
	the database binding in question may be in "auto commit"
	mode by default. Once you call Begin(), "auto commit" will
	be disabled for that connection.
*/
type TransactionalConnection interface {
	Connection;
	/*
		Begin() starts a transaction.
	*/
	Begin() os.Error;
	/*
		Commit() tries to push all changes made as part
		of the current transaction to the database.
	*/
	Commit() os.Error;
	/*
		Rollback() tries to undo all changes made as
		part of the current transaction.
	*/
	Rollback() os.Error
}

/*
	Statements are precompiled queries, possibly with remaining
	parameter slots that need to be filled before execution.
	TODO: include parameter binding API? or subsume in Execute()?
	what about resetting the statement or clearing parameter
	bindings?
*/
type Statement interface {
	Close() os.Error;
}

/*
	A call to Execute() that generates results from the database
	system returns a Cursor to allow clients to iterate through
	the results. Specific database binding will return
	cursor objects conforming to one or more of the following
	interfaces which represent different levels of functionality.

	TODO: API based on iterable instead?
*/
type Cursor interface {
	/*
		Are there more results to be fetched?
	*/
	MoreResults() bool;
	/*
		Fetch the next result from the database. A result
		is returned as a array of generic objects, one for
		each field. The database binding in question has
		to define what concrete types are returned depending
		on the types used in the database system.
	*/
	FetchOne() ([]interface {}, os.Error);
	/*
		Fetch at most count results. MAY GO AWAY SOON.
	*/
	FetchMany(count int) ([][]interface {}, os.Error);
	/*
		Fetch all (remaining) results. MAY GO AWAY SOON.
	*/
	FetchAll() ([][]interface {}, os.Error);
	/*
		Close() frees the cursor. After a cursor has been
		closed, no further operations are allowed on it.
	*/
	Close() os.Error
}

type InformativeCursor interface {
	Cursor;
	/*
		Description() returns a map from (the name of) a field to
		(the name of) its type. The exact format of field and type
		names is specified by the database binding in question.
	*/
	Description() (map[string]string, os.Error);
	/*
		Results returns the number of results remaining to be
		Fetch()ed.
	*/
	Results() int;
};

type PythonicCursor interface {
	Cursor;
	/*
		FetchDict() is similar to FetchOne() but instead of
		returning a slice it returns a map from fields names
		to values instead.
	*/
        FetchDict() (data map[string]interface{}, error os.Error);
	/*
		Fetch at most count results. MAY GO AWAY SOON.
	*/
        FetchManyDicts(count int) (data []map[string]interface{}, error os.Error);
	/*
		Fetch all remaining results. MAY GO AWAY SOON.
	*/
        FetchAllDicts() (data []map[string]interface{}, error os.Error)
};
