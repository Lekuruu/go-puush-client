package updater

import "fmt"

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) Compare(other Version) int {
	if v.Major != other.Major {
		return v.Major - other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor - other.Minor
	}
	return v.Patch - other.Patch
}

func (v Version) IsNewerThan(other Version) bool {
	return v.Compare(other) > 0
}

func (v Version) IsOlderThan(other Version) bool {
	return v.Compare(other) < 0
}

func (v Version) IsEqualTo(other Version) bool {
	return v.Compare(other) == 0
}

func NewVersion(major, minor, patch int) Version {
	return Version{Major: major, Minor: minor, Patch: patch}
}

func NewVersionFromString(versionString string) (Version, error) {
	var major, minor, patch int
	n, err := fmt.Sscanf(versionString, "%d.%d.%d", &major, &minor, &patch)
	if err != nil {
		return Version{}, err
	}
	if n != 3 {
		return Version{}, fmt.Errorf("invalid version string: %s", versionString)
	}
	return NewVersion(major, minor, patch), nil
}
