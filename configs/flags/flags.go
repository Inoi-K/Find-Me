package flags

import "flag"

var (
	DatabaseURL = flag.String("db-url", "", "Database URL")
)

func init() {
	flag.Parse()
}
