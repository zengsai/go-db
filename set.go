// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import "os"

// The simplest way of processing query results.
//
// Iter() returns a channel of Result objects which can be
// examined one at a time. Calling Iter() repeatedly will
// return the same channel over and over again, ResultSet
// does not cache results. Once Close() has been called,
// Iter() will return nil.
//
// Close() shuts down the iteration mechanism and frees all
// associated resources.
type ResultSet interface {
	Iter() <-chan Result;
	Close() os.Error;
}

// InformativeResultSets supply useful but optional information.
//
// Names() returns the names of each item of data in the
// result.
//
// Types() returns the names of the types of each item in
// the result.
type InformativeResultSet interface {
	Names() []string;
	Types() []string;
}
