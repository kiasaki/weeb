package main

import (
	"crypto/rand"
	"fmt"
	"strings"
)

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

func camelCase(givenValue string) string {
	value := strings.ToLower(givenValue)
	valueParts := strings.Split(value, "_")
	out := ""
	for _, part := range valueParts {
		out += strings.ToUpper(string(part[0])) + part[1:]
	}
	return out
}
