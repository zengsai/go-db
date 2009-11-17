/*
	Error for SQLite3.
*/

package sqlite3

type Error struct {
	code int;
	extended int; // TODO: unused so far
	message string;
}

func (self *Error) String() string {
	// TODO: better message that includes code(s)?
	return self.message;
}

func (self *Error) Code() int {
	return self.code;
}

func (self *Error) Extended() int {
	// TODO: always 0 for now
	return self.extended;
}
