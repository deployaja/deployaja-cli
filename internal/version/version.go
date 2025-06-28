package version

// Version will be set during build time with -ldflags
var Version = "dev"

// GetVersion returns the current CLI version
func GetVersion() string {
	if Version == "" {
		return "dev"
	}
	return Version
}
