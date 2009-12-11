// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import "os"

// TODO: the new way of returning results from Execute()
type ResultSet interface {
	Iter() <-chan Result;
	Close() os.Error;
}
