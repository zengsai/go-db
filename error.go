/*
	THIS IS NOT DONE AT ALL! USE AT YOUR OWN RISK!
*/

package sqlite3

/*
	SQLite has basic and extended error codes
	in addition to textual messages.
*/
type Error struct {
	basic int;
	extended int;
	message string;
}

/*
	Textual description of an error. Makes us
	compatible with the os.Error interface.
*/
func (self *Error) String() string {
	return self.message;
}

/*
	Basic SQLite error code. These are plain
	integers.
*/
func (self *Error) Basic() int {
	return self.basic;
}

/*
	Extended SQLite error code. These are OR'd
	together from various bits and pieces on top
	of basic error codes.
*/
func (self *Error) Extended() int {
	return self.extended;
}
