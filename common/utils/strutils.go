package utils

import (
	"bytes"
	"io"
	"strings"
)

//StrInSliceIgnoreCase whether a string in a string slice ignore cases
func StrInSliceIgnoreCase(str string, sli []string) bool {

	for _, s := range sli {
		if strings.EqualFold(s, str) {
			return true
		}
	}

	return false
}

//ValidateStr return true if lengthof str != 0
func ValidateStr(str string) bool {

	if (len(str)) == 0 {
		return false
	}
	return true
}

//Path2dbName convert a path string like '/home/user/dir' to string like 'home_user_dir' as mongoDB name
func Path2dbName(spath string) string {

	spath = strings.TrimPrefix(spath, "/")

	return strings.Replace(spath, "/", "_", -1)

}

//StreamToByte io.Reader to []byte
func StreamToByte(stream io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	return buf.Bytes(), err
}

// StreamToString io.Reader to string
func StreamToString(stream io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	return buf.String(), err
}
