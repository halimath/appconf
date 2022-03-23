package appconf

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type AppConfig struct {
	n *Node
}

func (c *AppConfig) HasKey(key string) bool {
	_, err := c.get(key)
	return err == nil || errors.Is(err, ErrNotAScalar)
}

func (c *AppConfig) Sub(key string) *AppConfig {
	s, err := c.SubE(key)
	if err != nil {
		return &AppConfig{&Node{}}
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
	v, err := c.GetStringE(key)
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			return ""
		}
		panic(fmt.Sprintf("invalid duration value for key %s: %s", key, err))
	}

	return v
}

func (c *AppConfig) GetStringE(key string) (string, error) {
	n, err := c.get(key)
	if err != nil {
		return "", err
	}

	if len(n.Children) != 0 {
		return "", ErrNotAScalar
	}

	return n.Value, nil
}

func (c *AppConfig) GetInt(key string) int {
	v, err := c.GetIntE(key)
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			return 0
		}
		panic(fmt.Sprintf("invalid int value for key %s: %s", key, err))
	}

	return v
}

func (c *AppConfig) GetIntE(key string) (int, error) {
	v, err := c.GetInt64E(key)
	return int(v), err
}

func (c *AppConfig) GetInt64(key string) int64 {
	v, err := c.GetInt64E(key)
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			return 0
		}
		panic(fmt.Sprintf("invalid int64 value for key %s: %s", key, err))
	}

	return v
}

func (c *AppConfig) GetInt64E(key string) (int64, error) {
	v, err := c.GetStringE(key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(v, 10, 64)
}

func (c *AppConfig) GetUint(key string) uint {
	v, err := c.GetUintE(key)
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			return 0
		}
		panic(fmt.Sprintf("invalid uint value for key %s: %s", key, err))
	}

	return v
}

func (c *AppConfig) GetUintE(key string) (uint, error) {
	v, err := c.GetUint64E(key)
	return uint(v), err
}

func (c *AppConfig) GetUint64(key string) uint64 {
	v, err := c.GetUint64E(key)
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			return 0
		}
		panic(fmt.Sprintf("invalid uint64 value for key %s: %s", key, err))
	}

	return v
}

func (c *AppConfig) GetUint64E(key string) (uint64, error) {
	v, err := c.GetStringE(key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(v, 10, 64)
}

func (c *AppConfig) GetFloat32(key string) float32 {
	v, err := c.GetFloat32E(key)
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			return 0
		}
		panic(fmt.Sprintf("invalid float32 value for key %s: %s", key, err))
	}

	return v
}

func (c *AppConfig) GetFloat32E(key string) (float32, error) {
	v, err := c.GetStringE(key)
	if err != nil {
		return 0, err
	}

	f, err := strconv.ParseFloat(v, 32)
	return float32(f), err
}

func (c *AppConfig) GetFloat64(key string) float64 {
	v, err := c.GetFloat64E(key)
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			return 0
		}
		panic(fmt.Sprintf("invalid float64 value for key %s: %s", key, err))
	}

	return v
}

func (c *AppConfig) GetFloat64E(key string) (float64, error) {
	v, err := c.GetStringE(key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(v, 64)
}

func (c *AppConfig) GetComplex128(key string) complex128 {
	v, err := c.GetComplex128E(key)
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			return 0
		}
		panic(fmt.Sprintf("invalid complex128 value for key %s: %s", key, err))
	}

	return v
}

func (c *AppConfig) GetComplex128E(key string) (complex128, error) {
	v, err := c.GetStringE(key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseComplex(v, 128)
}

func (c *AppConfig) GetBool(key string) bool {
	v, err := c.GetBoolE(key)
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			return false
		}
		panic(fmt.Sprintf("invalid bool value for key %s: %s", key, err))
	}

	return v
}

func (c *AppConfig) GetBoolE(key string) (bool, error) {
	v, err := c.GetStringE(key)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(v)
}

func (c *AppConfig) GetDuration(key string) time.Duration {
	d, err := c.GetDurationE(key)
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			return 0
		}
		panic(fmt.Sprintf("invalid duration value for key %s: %s", key, err))
	}

	return d
}

func (c *AppConfig) GetDurationE(key string) (time.Duration, error) {
	v, err := c.GetStringE(key)
	if err != nil {
		return 0, err
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, err
	}
	return d, nil
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
		n: &Node{
			Children: make(map[Key]*Node),
		},
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
