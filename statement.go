// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import "os"

// Statements are precompiled queries, possibly with remaining
// parameter slots that need to be filled before execution.
// TODO: include parameter binding API? or subsume in Execute()?
// what about resetting the statement or clearing parameter
// bindings?
type Statement interface {
	Close() os.Error;
}
