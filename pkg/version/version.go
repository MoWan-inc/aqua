package version

import "fmt"

var (
	GitCommit string
	BuildDate string
)

func String() string {
	return fmt.Sprintf("%s(BuiltAt %s)", GitCommit, BuildDate)
}
