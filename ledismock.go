package ledismock

import (
	"fmt"
)

type HashValue struct {
	key   string
	value interface{}
}

func NewHashValue(key string, value interface{}) *HashValue {
	return &HashValue{
		key:   key,
		value: value,
	}
}

func NewHashValues(hashValues ...*HashValue) []*HashValue {
	result := make([]*HashValue, len(hashValues), len(hashValues))
	for i, val := range hashValues {
		result[i] = val
	}

	return result
}

// LedisMock is the top level mock object to simulate interactions with a Ledis database.
type LedisMock struct {
	expectations []*Command
	received     []*internalCommand
}

// Command is one expected query.
type Command struct {
	name  string
	key   string
	value string

	internal []*internalCommand

	result interface{}
}

type internalCommand struct {
	funcName string
	args     []string
}

// ExpectSAdd will add a SADD command to the expectations
func (m *LedisMock) ExpectSAdd() *Command {
	c := Command{
		name: "sadd",
	}
	m.expectations = append(m.expectations, &c)

	return &c
}

// ExpectHMset will ad a HMSET command to the expectations
func (m *LedisMock) ExpectHMset() *Command {
	c := Command{
		name: "hmset",
	}
	m.expectations = append(m.expectations, &c)

	return &c
}

// ExpectationsWereMet will return an error if the are unfufilled expectated queries.
func (m *LedisMock) ExpectationsWereMet() error {
	// build a flattened list of expectations
	expect := make([]*internalCommand, 0, 0)
	for _, expectation := range m.expectations {
		for _, internal := range expectation.internal {
			expect = append(expect, internal)
		}
	}

	// TODO all of these return errors showing an internal command, which will be
	// very hard to debug with.  These should get mapped back to the expected
	// commands that were passed in originally...

	for i, exp := range expect {
		if i >= len(m.received) {
			return fmt.Errorf("Expectation not met: %#v", expect)
		}

		//fmt.Printf("Expect to see: %q with args %#v, saw %q with args %#v, validating...", exp.funcName, exp.args, m.received[i].funcName, m.received[i].args)

		if m.received[i].funcName != exp.funcName {
			return fmt.Errorf("Expected to see %q, but saw %q.", exp.funcName, m.received[i].funcName)
		}

		if len(m.received[i].args) != len(exp.args) {
			return fmt.Errorf("Expected to see %d args, but saw %d instead.", len(exp.args), len(m.received[i].args))
		}

		// This is save from referencing objects out of range because of the length check immediately above
		for j, arg := range exp.args {
			if arg != m.received[i].args[j] {
				return fmt.Errorf("Expected to see arg %q, but saw %q instead.", arg, m.received[i].args[j])
			}
		}
	}

	return nil
}

func (m *LedisMock) receivedGet(key []byte) {
	m.received = append(m.received, &internalCommand{funcName: "GET", args: []string{string(key)}})
}

// WithKey sets the expectation that the command will operate on a specific key.
func (c *Command) WithKey(key string) *Command {
	c.key = key

	return c
}

// WithValue sets the expectation that the specific value will be set on the command.
func (c *Command) WithValue(value string) *Command {
	c.value = value

	internal := make([]*internalCommand, 0, 0)

	switch c.name {
	case "sadd":
		internal = append(internal, &internalCommand{funcName: "GET", args: []string{fmt.Sprintf("\x00\v\x00)%s:%s", c.key, value)}})
		internal = append(internal, &internalCommand{funcName: "GET", args: []string{fmt.Sprintf("\x00\f%s", c.key)}})
	}

	c.internal = internal

	return c
}

// WithValues sets the expectation that the specific values will be set on the command, but order is not important.
// This should be extended into an interface to get full support.
func (c *Command) WithValues(values []*HashValue) *Command {

	internal := make([]*internalCommand, 0, 0)

	switch c.name {
	case "hmset":
		for _, val := range values {
			internal = append(internal, &internalCommand{funcName: "GET", args: []string{fmt.Sprintf("\x00\x02\x00*%s:%s", c.key, val.key)}})
		}
	}

	c.internal = internal

	return c

}

// WillReturnResult sets the mock value to be returned from the mock query.
func (c *Command) WillReturnResult(result interface{}) *Command {
	c.result = result
	return c
}
