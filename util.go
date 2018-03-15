package weeb

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/securecookie"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func init() {
	gob.Register(&J{})
	gob.Register(&Flash{})

	sqlx.NameMapper = ToSnakeCase

	godotenv.Load()
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

func mergeStringMaps(a, b map[string]string) map[string]string {
	c := map[string]string{}
	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}
	return c
}

func containsString(values []string, search string) bool {
	for _, value := range values {
		if value == search {
			return true
		}
	}
	return false
}

func title(value string) string {
	if len(value) == 0 {
		return value
	}
	return strings.ToUpper(string(value[0])) + value[1:]
}

// OrString returns the first non-empty string
func OrString(options ...string) string {
	for _, o := range options {
		if len(o) > 0 {
			return o
		}
	}
	return ""
}

const randomKeyDict = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func generateRandomKey(length int) string {
	key := securecookie.GenerateRandomKey(length)
	if key == nil {
		panic("securecookie.GenerateRandomKey returned nil")
	}
	for i := range key {
		key[i] = randomKeyDict[key[i]%byte(len(randomKeyDict))]
	}
	return string(key)
}

// from golang/src/net/http/http.go:62
func hexEscapeNonASCII(s string) string {
	newLen := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			newLen += 3
		} else {
			newLen++
		}
	}
	if newLen == len(s) {
		return s
	}
	b := make([]byte, 0, newLen)
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			b = append(b, '%')
			b = strconv.AppendInt(b, int64(s[i]), 16)
		} else {
			b = append(b, s[i])
		}
	}
	return string(b)
}
