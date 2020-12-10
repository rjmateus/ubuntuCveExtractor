package main

type PackagesChannel struct {
	Channel  string
	Packages []string
}

type OsData struct {
	Name, Url string
	PackagesChannelList []PackagesChannel
}
type ReleaseStatus struct {
	FullName, shortName, versionName string
	Status string
	PackageReleasedVersion string
}

type PackageInfo struct {
	PackageName   string
	//ReleaseStatus map[string]ReleaseStatus
	ReleaseStatus map[string]string
	OsPackages    map[string]OsData
}


type CveData struct {
	CveID                string
	UbuntuCvePublishDate string
	Packages             map[string]PackageInfo
	//References           []string
}
