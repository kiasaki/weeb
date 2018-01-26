package weeb

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

func init() {
	gob.Register(&J{})
	gob.Register(&Flash{})
}

type contextKey int

// J is a shorthand used to build JSON values
type J map[string]interface{}

// UUID represents a UUID value
type UUID [16]byte

// NewUUID Creates a new UUID v4, panicing if it fails to get entropy
func NewUUID() *UUID {
	u := &UUID{}
	_, err := rand.Read(u[:16])
	if err != nil {
		panic(err)
	}

	u[8] = (u[8] | 0x80) & 0xBf
	u[6] = (u[6] | 0x40) & 0x4f
	return u
}

func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func dirExists(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// ToCamelCase converts a string to camel case.
func ToCamelCase(givenValue string) string {
	value := strings.ToLower(givenValue)
	valueParts := strings.Split(value, "_")
	out := ""
	for i, part := range valueParts {
		firstChar := strings.ToUpper(string(part[0]))
		if i == 0 {
			// First part is lowercase
			firstChar = string(part[0])
		}
		out += firstChar + part[1:]
	}
	return out
}

// ToSnakeCase converts a string to snake case, words separated with underscores.
func ToSnakeCase(src string) string {
	thisUpper := false
	prevUpper := false

	buf := bytes.NewBufferString("")
	for i, v := range src {
		if v >= 'A' && v <= 'Z' {
			thisUpper = true
		} else {
			thisUpper = false
		}
		if i > 0 && thisUpper && !prevUpper {
			buf.WriteRune('_')
		}
		prevUpper = thisUpper
		buf.WriteRune(v)
	}
	return strings.ToLower(buf.String())
}

func displayMap(valuesMap map[string]string, leftPadding, keyPadding int) string {
	out := ""
	longestKey := 0
	for k := range valuesMap {
		if len(k) > longestKey {
			longestKey = len(k)
		}
	}
	for k, v := range valuesMap {
		for i := 0; i < leftPadding; i++ {
			out += " "
		}
		out += k
		for i := 0; i < longestKey-len(k)+keyPadding; i++ {
			out += " "
		}
		out += v
		out += "\n"
	}
	return out
}
