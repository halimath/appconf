// Package appconf provides a flexible configuration loader for applications that wish to load configuration
// values from different sources, such as files or the environment. It supports nested data structures,
// type-safe conversions and binding using reflection.
package appconf

import (
	"errors"
	"fmt"
	"time"
)

// AppConf is the main data type used to interact with configuration values.
type AppConfig struct {
	n *Node
}

// HasKey returns whether c contains key which may be nested key.
func (c *AppConfig) HasKey(key string) bool {
	_, err := c.get(key)
	return err == nil || errors.Is(err, ErrNotAScalar)
}

// Sub returns the sub-structure of c rooted at key.
func (c *AppConfig) Sub(key string) *AppConfig {
	s, err := c.SubE(key)
	if err != nil {
		return &AppConfig{NewNode("")}
	}

	return s
}

// SubE returns the sub-structure of c rooted at key or an error if the key does not exist.
func (c *AppConfig) SubE(key string) (*AppConfig, error) {
	n, err := c.get(key)
	if err != nil {
		return nil, err
	}
	return &AppConfig{n: n}, nil
}

// GetString returns the string value stored under key.
func (c *AppConfig) GetString(key string) string {
	n, err := c.get(key)
	if err != nil {
		return ""
	}
	return n.GetString()
}

// GetStringE returns the string value stored under key or an error if the key is not defined.
func (c *AppConfig) GetStringE(key string) (string, error) {
	n, err := c.get(key)
	if err != nil {
		return "", err
	}

	return n.GetStringE()
}

// GetInt returns the int value stored under key.
func (c *AppConfig) GetInt(key string) int {
	n, err := c.get(key)
	if err != nil {
		return 0
	}

	return n.GetInt()
}

// GetIntE returns the int value stored under key or an error if the key is not defined or the value cannot
// be converted to an int.
func (c *AppConfig) GetIntE(key string) (int, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetIntE()
}

// GetInt64 returns the int64 value stored under key.
func (c *AppConfig) GetInt64(key string) int64 {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetInt64()
}

// GetInt64E returns the int64 value stored under key or an error if the key is not defined or the value
// cannot be converted to an int64.
func (c *AppConfig) GetInt64E(key string) (int64, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetInt64E()
}

// GetUint returns the uint value stored under key.
func (c *AppConfig) GetUint(key string) uint {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetUint()
}

// GetUintE returns the uint value stored under key or an error if the key is not defined or the value
// cannot be converted to an uint.
func (c *AppConfig) GetUintE(key string) (uint, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetUintE()
}

// GetUint64 returns the uint64 value stored under key.
func (c *AppConfig) GetUint64(key string) uint64 {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetUint64()
}

// GetUint64E returns the uint64 value stored under key or an error if the key is not defined or the value
// cannot be converted to an uint64.
func (c *AppConfig) GetUint64E(key string) (uint64, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetUint64E()
}

// GetFloat32 returns the float32 value stored under key.
func (c *AppConfig) GetFloat32(key string) float32 {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetFloat32()
}

// GetFloat32E returns the float32 value stored under key or an error if the key is not defined or the value
// cannot be converted to an float32.
func (c *AppConfig) GetFloat32E(key string) (float32, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetFloat32E()
}

// GetFloat64 returns the float64 value stored under key.
func (c *AppConfig) GetFloat64(key string) float64 {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetFloat64()
}

// GetFloat64E returns the float64 value stored under key or an error if the key is not defined or the value
// cannot be converted to an float64.
func (c *AppConfig) GetFloat64E(key string) (float64, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetFloat64E()
}

// GetComplex128 returns the complex128 value stored under key.
func (c *AppConfig) GetComplex128(key string) complex128 {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetComplex128()
}

// GetComplex128E returns the complex128 value stored under key or an error if the key is not defined or the value
// cannot be converted to an complex128.
func (c *AppConfig) GetComplex128E(key string) (complex128, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetComplex128E()
}

// GetBool returns the bool value stored under key.
func (c *AppConfig) GetBool(key string) bool {
	n, err := c.get(key)
	if err != nil {
		return false
	}
	return n.GetBool()
}

// GetBoolE returns the bool value stored under key or an error if the key is not defined or the value
// cannot be converted to an bool.
func (c *AppConfig) GetBoolE(key string) (bool, error) {
	n, err := c.get(key)
	if err != nil {
		return false, err
	}
	return n.GetBoolE()
}

// GetDuration returns the duration value stored under key.
func (c *AppConfig) GetDuration(key string) time.Duration {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetDuration()
}

// GetDurationE returns the duration value stored under key or an error if the key is not found or the string
// value cannot be parsed into a Duration.
func (c *AppConfig) GetDurationE(key string) (time.Duration, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetDurationE()
}

// Bind binds the configuration to the data structure v and returns any error that occured during binding.
// v must be a pointer to either a struct value or a map[string]interface{}. Other values are not supported
// and are rejected by an error. See the README for an explanation of how to use and customize the binding.
func (c *AppConfig) Bind(v interface{}) error {
	return bind(c.n, v)
}

func (c *AppConfig) get(key string) (*Node, error) {
	n := c.n.resolve(ParseKeyPath(key))
	if n == nil {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchKey, key)
	}
	return n, nil
}

// New creates a new AppConfig using the given loaders. The loaders are executed in given order with values
// from later loaders overwriting values from earlier ones (put most significant loaders last).
func New(loaders ...Loader) (*AppConfig, error) {
	c := &AppConfig{
		n: NewNode(""),
	}

	for _, l := range loaders {
		n, err := l.Load()
		if err != nil {
			return nil, err
		}
		c.n.OverwriteWith(n)
	}

	return c, nil
}
