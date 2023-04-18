package app_version_manager

type appVersionManager struct {
	VersionStatus versionStatus
}

var AppVersionManager = appVersionManager{
	VersionStatus: Status,
}
