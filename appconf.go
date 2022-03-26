package appconf

import (
	"errors"
	"fmt"
	"time"
)

type AppConfig struct {
	n *Node
}

func (c *AppConfig) Keys() []Key {
	keys := make([]Key, 0, len(c.n.Children))
	for k := range c.n.Children {
		keys = append(keys, k)
	}
	return keys
}

func (c *AppConfig) HasKey(key string) bool {
	_, err := c.get(key)
	return err == nil || errors.Is(err, ErrNotAScalar)
}

func (c *AppConfig) Sub(key string) *AppConfig {
	s, err := c.SubE(key)
	if err != nil {
		return &AppConfig{NewNode("")}
	}

	return s
}

func (c *AppConfig) SubE(key string) (*AppConfig, error) {
	n, err := c.get(key)
	if err != nil {
		return nil, err
	}
	return &AppConfig{n: n}, nil
}

func (c *AppConfig) GetString(key string) string {
	n, err := c.get(key)
	if err != nil {
		return ""
	}
	return n.GetString()
}

func (c *AppConfig) GetStringE(key string) (string, error) {
	n, err := c.get(key)
	if err != nil {
		return "", err
	}

	return n.GetStringE()
}

func (c *AppConfig) GetInt(key string) int {
	n, err := c.get(key)
	if err != nil {
		return 0
	}

	return n.GetInt()
}

func (c *AppConfig) GetIntE(key string) (int, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetIntE()
}

func (c *AppConfig) GetInt64(key string) int64 {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetInt64()
}

func (c *AppConfig) GetInt64E(key string) (int64, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetInt64E()
}

func (c *AppConfig) GetUint(key string) uint {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetUint()
}

func (c *AppConfig) GetUintE(key string) (uint, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetUintE()
}

func (c *AppConfig) GetUint64(key string) uint64 {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetUint64()
}

func (c *AppConfig) GetUint64E(key string) (uint64, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetUint64E()
}

func (c *AppConfig) GetFloat32(key string) float32 {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetFloat32()
}

func (c *AppConfig) GetFloat32E(key string) (float32, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetFloat32E()
}

func (c *AppConfig) GetFloat64(key string) float64 {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetFloat64()
}

func (c *AppConfig) GetFloat64E(key string) (float64, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetFloat64E()
}

func (c *AppConfig) GetComplex128(key string) complex128 {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetComplex128()
}

func (c *AppConfig) GetComplex128E(key string) (complex128, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetComplex128E()
}

func (c *AppConfig) GetBool(key string) bool {
	n, err := c.get(key)
	if err != nil {
		return false
	}
	return n.GetBool()
}

func (c *AppConfig) GetBoolE(key string) (bool, error) {
	n, err := c.get(key)
	if err != nil {
		return false, err
	}
	return n.GetBoolE()
}

func (c *AppConfig) GetDuration(key string) time.Duration {
	n, err := c.get(key)
	if err != nil {
		return 0
	}
	return n.GetDuration()
}

func (c *AppConfig) GetDurationE(key string) (time.Duration, error) {
	n, err := c.get(key)
	if err != nil {
		return 0, err
	}
	return n.GetDurationE()
}

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
