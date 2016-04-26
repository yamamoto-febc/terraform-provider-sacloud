package sacloud

import (
	libsacloud "github.com/yamamoto-febc/libsacloud/api"
)

type Config struct {
	AccessToken       string
	AccessTokenSecret string
	Region            string
}

func (c *Config) NewClient() *libsacloud.Client {
	return libsacloud.NewClient(c.AccessToken, c.AccessTokenSecret, c.Region)
}
