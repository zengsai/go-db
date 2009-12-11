// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import "os"

// The simplest way of processing query results.
//
// Iter() returns a channel of Result objects which can be
// examined one at a time. Note that called Iter() several
// times will return the same channel over and over again.
//
// Close() shuts down the iteration mechanism and frees all
// associated resources. After a result set has been closed,
// Iter() will return nil.
type ResultSet interface {
	Iter() <-chan Result;
	Close() os.Error;
}
