package config

import (
	"os"
	"path"
)

var CachePath = path.Join(os.TempDir(), "podcast-cli/")
