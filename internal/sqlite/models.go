// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlite

import ()

type Hook struct {
	Name   string
	Script string
}

type KeyHook struct {
	Key  string
	Hook string
}

type Kv struct {
	Key string
	Val string
}
