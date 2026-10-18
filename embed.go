package fibera

import "embed"

//go:embed configs/*.yaml
var ConfigFile embed.FS

//go:embed .env
var EnvFile embed.FS
