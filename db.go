// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import "os"

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
type OpenSignature func(url string) (Connection, os.Error)
