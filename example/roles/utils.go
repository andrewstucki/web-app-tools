package roles

import "github.com/andrewstucki/web-app-tools/go/security"

// Lifted from: https://stackoverflow.com/questions/18879109/subset-check-with-slices-in-go
// subset returns true if the first array is completely
// contained in the second array. There must be at least
// the same number of duplicate values in second as there
// are in first.
func subset(first, second []security.Policy) bool {
	set := make(map[security.Policy]int)
	for _, value := range second {
		set[value]++
	}

	for _, value := range first {
		if count, found := set[value]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			set[value] = count - 1
		}
	}

	return true
}

// IsAdmin checks if the policies contain the super-admin policies
func IsAdmin(policies []security.Policy) bool {
	return subset(SuperAdminRole.Policies, policies)
}
