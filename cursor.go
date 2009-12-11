// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import "os"

// DEPRECATED
type Cursor interface {
	MoreResults() bool;
	FetchOne() ([]interface{}, os.Error);
	FetchMany(count int) ([][]interface{}, os.Error);
	FetchAll() ([][]interface{}, os.Error);
	Close() os.Error;
}

// DEPRECATED
type InformativeCursor interface {
	Cursor;
	Description() (map[string]string, os.Error);
	Results() int;
}
