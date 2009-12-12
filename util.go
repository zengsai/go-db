// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import "os"
import "strings"

// ExecuteDirectly is a convenience function for "one-off" queries.
// It's particularly convenient for queries that don't produce any
// results.
//
// If you need more control, for example to rebind parameters over
// and over again, to get results one by one, or to access metadata
// about the results, you should use the Prepare() and Execute()
// methods explicitly instead.
//
// TODO: results should be returned some other way...
func ExecuteDirectly(conn Connection, query string, params ...) (results [][]interface{}, err os.Error) {
	var s Statement;
	s, err = conn.Prepare(query);
	if err != nil || s == nil {
		return
	}
	defer s.Close();

	var c ClassicResultSet;
	con := conn.(ClassicConnection);
	c, err = con.ExecuteClassic(s, params);
	if err != nil || c == nil {
		return
	}
	defer c.Close();

	results, err = ClassicFetchAll(c);
	return;
}

// ParseQueryURL() helps database drivers parse URLs passed
// to Open(). ParseQueryURL() takes a string of the form
//
//	key=value;key=value;...;key=value
//
// and returns a map from keys to values. The empty string
// yields an empty map. Format violations or duplicate keys
// yield an error and an incomplete map.
func ParseQueryURL(str string) (opt map[string]string, err os.Error) {
	opt = make(map[string]string);
	if len(str) > 0 {
		err = parseQueryHelper(str, opt);
	}
	return;
}

func parseQueryHelper(str string, opt map[string]string) (err os.Error) {
	pairs := strings.Split(str, ";", 0);
	if len(pairs) == 0 {
		err = os.NewError("ParseQueryURL: No pairs in "+str);
		return; // nothing left to do
	}

	for _, p := range pairs {
		pieces := strings.Split(p, "=", 0);
		// we keep going even if there was an error to fill the
		// map as much as possible; this means we'll return only
		// the last error, a tradeoff
		if len(pieces) == 2 {
			if _, duplicate := opt[pieces[0]]; duplicate {
				err = os.NewError("ParseQueryURL: Duplicate key "+pieces[0]);
			}
			else {
				opt[pieces[0]] = pieces[1]
			}
		}
		else {
			err = os.NewError("ParseQueryURL: One '=' expected, got "+p);
		}
	}

	return;
}
