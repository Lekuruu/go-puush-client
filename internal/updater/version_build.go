package updater

import "fmt"

type VersionBuild struct {
	CommitNumber int
}

func (v VersionBuild) String() string {
	return fmt.Sprintf("%d", v.CommitNumber)
}

func (v VersionBuild) CanCompare(other Version) bool {
	_, ok := other.(VersionBuild)
	return ok
}

func (v VersionBuild) Compare(anyOther Version) int {
	other, ok := anyOther.(VersionBuild)
	if !ok {
		panic(fmt.Sprintf("Cannot compare VersionBuild with %T", anyOther))
	}
	return v.CommitNumber - other.CommitNumber
}

func (v VersionBuild) IsNewerThan(other Version) bool {
	return v.Compare(other) > 0
}

func (v VersionBuild) IsOlderThan(other Version) bool {
	return v.Compare(other) < 0
}

func (v VersionBuild) IsEqualTo(other Version) bool {
	return v.Compare(other) == 0
}

func NewBuildVersion(commitNumber int) VersionBuild {
	return VersionBuild{CommitNumber: commitNumber}
}

func NewBuildVersionFromString(versionString string) (VersionBuild, error) {
	var commitNumber int
	n, err := fmt.Sscanf(versionString, "%d", &commitNumber)
	if err != nil {
		return VersionBuild{}, err
	}
	if n != 1 {
		return VersionBuild{}, fmt.Errorf("invalid version string: %s", versionString)
	}
	return NewBuildVersion(commitNumber), nil
}
