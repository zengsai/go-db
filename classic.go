// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import "os"
import "container/vector"

// ClassicConnections avoid the use of channels for results.
// They perform slightly better (at least while we wait for
// the Go runtime to get tuned more) but are less convenient
// to use because the for/range construct doesn't apply.
//
// ExecuteClassic() is similar to Execute() except that it
// returns a ClassicResultSet (see below).
type ClassicConnection interface {
	Connection;
	ExecuteClassic(stat Statement, parameters ...) (ClassicResultSet, os.Error);
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

// TODO
func ClassicFetchAll(rs ClassicResultSet) (data [][]interface{}, error os.Error) {
	var v vector.Vector;
	var d interface{}
	var e os.Error;

	for rs.More() {
		r := rs.Fetch();
		d = r.Data();
		if d != nil {
			v.Push(d);
		}
		e = r.Error();
		if e != nil {
			break;
		}
	}

	l := v.Len();

	if l > 0 {
		// TODO: how can this be done better?
		data = make([][]interface{}, l);
		for i := 0; i < l; i++ {
			data[i] = v.At(i).([]interface{})
		}
	} else {
		// no results at all, return the error
		error = e
	}

	return;
}

// TODO
func ClassicFetchMany(rs ClassicResultSet, count int) (data [][]interface{}, error os.Error) {
	d := make([][]interface{}, count);
	l := 0;
	var e os.Error;

	// grab at most count results
	for l < count {
		r := rs.Fetch();
		d[l] = r.Data();
		e = r.Error();
		if e == nil {
			l += 1
		} else {
			break
		}
	}

	if l > 0 {
		// there were results
		if l < count {
			// but fewer than expected, need fresh copy
			data = make([][]interface{}, l);
			for i := 0; i < l; i++ {
				data[i] = d[i]
			}
		} else {
			data = d
		}
	} else {
		// no results at all, return the error
		error = e
	}

	return;
}
