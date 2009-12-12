// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import "os"

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
