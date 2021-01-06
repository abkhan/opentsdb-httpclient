package tsdbh

import "time"

// Conf has http client conf
type Conf struct {
	Host            string
	Port            int
	DefaultBuffer   int
	DefaultInterval int
	HTTPConf        HTTPConf `mapstructure:"http"`
}

// HTTPConf keeps all http related configs
type HTTPConf struct {
	DialTimeout         time.Duration `mapstructure:"dialTimeout"`
	TLSHandshakeTimeout time.Duration `mapstructure:"tlsHandShakeTimeout"`
	MaxIdleConnsPerHost int           `mapstructure:"maxIdleConnsPerHost"`
	MaxIdleConns        int           `mapstructure:"maxIdleConns"`
	IdleConnTimeout     time.Duration `mapstructure:"idleConnTimeout"`
	ClientTimeout       time.Duration `mapstructure:"clientTimeout"`
}
