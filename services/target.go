package services

import "fmt"

// Target represents a service whose group memberships can be mutated.
type Target interface {
	AddMembers(team string, users []User) error
	RemoveMembers(team string, users []User) error
	acquireIdentity(user *User) (Identity, error)

	// Target implementors should also implement Service.
	GroupMembers(group string) ([]User, error)
}

func TargetFromString(name string) (Target, error) {
	switch name {
	case "github":
		return githubSvc, nil
	default:
		return nil, fmt.Errorf(
			"no target %s defined",
			name,
		)
	}
}

type TargetNotDefined struct {
	serviceName string
}

func newTargetNotDefined(serviceName string) TargetNotDefined {
	return TargetNotDefined{
		serviceName: serviceName,
	}
}

func (e TargetNotDefined) Error() string {
	return fmt.Sprintf("target `%s` not defined", e.serviceName)
}
