package services

import "fmt"

// Target represents a service whose group memberships can be mutated.
type Target interface {
	AddMembers(users []User) error
	RemoveMembers(users []User) error

	// Target implementors should also implement Service.
	GroupMembers(group string) ([]User, error)
	acquireIdentity(user *User) (Identity, error)
}

func Diff(srcGrp, tarGrp []User, tar string) (rem, add []User) {
	// Build hashmaps of identities for faster lookup.
	// This approach also takes care of duplicates for free.
	srcMap := make(map[string]User)
	tarMap := make(map[string]User)

	if len(srcGrp) < 1 {
		panic("sanity check failed: the source group is empty")
	}

	for _, u := range srcGrp {
		i, err := u.getIdentity(tar)
		if err != nil {
			fmt.Errorf(
				"error acquiring identity for a user\nuser: %s\nerror: %s",
				u,
				err,
			)
		} else if IdentityExists(i) {
			srcMap[i.uniqueID()] = u
		}
	}

	for _, u := range tarGrp {
		i, err := u.getIdentity(tar)
		if err != nil {
			fmt.Errorf(
				"error acquiring identity for a user\nuser: %s\nerror: %s",
				u,
				err,
			)
		} else if IdentityExists(i) {
			tarMap[i.uniqueID()] = u
		}
	}

	// Remove elements that exist in both the source and the target.
	for id := range srcMap {
		_, ok := tarMap[id]
		if ok {
			delete(srcMap, id)
			delete(tarMap, id)
		}
	}

	// What's left in srcIdentities and tarIdentities is what we have to
	// add/remove.
	for _, identity := range srcMap {
		add = append(add, identity)
	}

	for _, identity := range tarMap {
		rem = append(rem, identity)
	}

	return
}

func createIDMap(ids []Identity) map[string]Identity {
	result := make(map[string]Identity)

	for _, id := range ids {
		result[id.uniqueID()] = id
	}

	return result
}