package djson

import "github.com/tidwall/sjson"

// Set sets a json value for the specified path.
// A path is in dot syntax, such as "name.last" or "age".
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return back unexpected results.
// An error is returned if the path is not valid.
//
// A path is a series of keys separated by a dot.
//
//  {
//    "name": {"first": "Tom", "last": "Anderson"},
//    "age":37,
//    "children": ["Sara","Alex","Jack"],
//    "friends": [
//      {"first": "James", "last": "Murphy"},
//      {"first": "Roger", "last": "Craig"}
//    ]
//  }
//  "name.last"          >> "Anderson"
//  "age"                >> 37
//  "children.1"         >> "Alex"
//
func Set(json, path string, value interface{}) (string, error) {
	return sjson.Set(json, path, value)
}

// SetBytes sets a json value for the specified path.
// If working with bytes, this method preferred over
// Set(string(data), path, value)
func SetBytes(json []byte, path string, value interface{}) ([]byte, error) {
	return sjson.SetBytes(json, path, value)
}

func Delete(s string, path string) (string, error) {
	return deleteJsonWithPaths(s, []string{path})
}

func Deletes(s string, paths []string) (string, error) {
	return deleteJsonWithPaths(s, paths)
}

func DeleteBytes(s []byte, path string) ([]byte, error) {
	j, err := deleteJsonWithPaths(string(s), []string{path})
	return []byte(j), err
}

func DeletesBytes(s []byte, paths []string) ([]byte, error) {
	j, err := deleteJsonWithPaths(string(s), paths)
	return []byte(j), err
}
