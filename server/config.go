package main

import (
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Root string
	TmpDir string
	Prefix string
	Host string
}

func (c *Config) Load() {
	c.Root = c.getEnv("FILES_ROOT", func() string {
		var here, _ = filepath.Abs(".");
		return here
	});
	c.TmpDir = "/tmp/zipfs" + c.Root;
	c.Prefix = c.getEnvS("FILES_PREFIX", "/")
	c.Host = c.getEnvS("FILES_HOST", "http://localhost:8080")
}

func (c *Config) Inspect() {
	log.Printf("Root: %s\nPrefix: %s", c.Root, c.Prefix)
}

func (c Config) getEnv(key string, def func() string) string {
	val, found := os.LookupEnv(key)
	if (found) { return val }
	return def()
}

func (c Config) getEnvS(key, def string) string {
	val, found := os.LookupEnv(key)
	if (found) { return val }
	return def
}
