/*
	Error for SQLite3.
*/

package sqlite3

type Error struct {
	code int;
	extended int;
	message string;
}

func (self *Error) String() string {
	return self.message;
}

func (self *Error) Code() int {
	return self.code;
}

func (self *Error) Extended() int {
	return self.extended;
}
