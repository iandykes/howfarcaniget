package main

import "fmt"

// VersionInfo encapsulates the build version details
type VersionInfo struct {
	Version    string
	BuildDate  string
	CommitHash string
}

// String returns the string value of the VersionInfo
func (v *VersionInfo) String() string {
	return fmt.Sprintf("%v %v %v", v.Version, v.BuildDate, v.CommitHash)
}
