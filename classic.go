// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import "os"

// ClassicConnections avoid the use of channels for results.
// They perform slightly better (at least while we wait for
// the Go runtime to get tuned) but are less convenient to
// use because the for/range construct doesn't apply.
//
// ExecuteClassic() is similar to Execute() except that it
// returns a ClassicResultSet (see below).
type ClassicConnection interface {
	Connection;
	ExecuteClassic(stat Statement, parameters ...) (Cursor, os.Error);
}

// ClassicResultSets offer the same functionality as regular
// ResultSets but without the use of channels.
//
// More() returns true if there are more results to fetch.
//
// Fetch() produces the next result.
//
// Close() frees all resources associated with the result
// set. After a call to close, no operations can be applied
// anymore.
type ClassicResultSet interface {
	More() bool;
	Fetch() Result;
	Close() os.Error;
}
