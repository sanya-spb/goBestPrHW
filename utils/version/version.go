package version

type AppVersion struct {
	Version   string
	Commit    string
	BuildTime string
	Copyright string
}

var Version *AppVersion = new(AppVersion)

// set vars from Makefile via go build -ldflags "-s -w -X ..."
var (
	version   string
	commit    string
	buildTime string
	copyright string
)

func init() {
	Version.Version = version
	Version.Commit = commit
	Version.BuildTime = buildTime
	Version.Copyright = copyright
}
