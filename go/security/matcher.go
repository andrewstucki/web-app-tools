package security

import (
	"context"
	"regexp"
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"
)

var (
	regexCache sync.Map
	tokenizer  *regexp.Regexp
)

func init() {
	tokenizer = regexp.MustCompile(`(.*):[^/]+(.*)`)
}

// Evaluator implements a permissions evaluation engine
type Evaluator struct {
	namespaceManager NamespaceManager
	namespace        uuid.UUID
	user             uuid.UUID
}

func newEvaluator(namespaceManager NamespaceManager, namespace uuid.UUID, user uuid.UUID) *Evaluator {
	return &Evaluator{
		namespaceManager: namespaceManager,
		namespace:        namespace,
		user:             user,
	}
}

// Can evaluates whether or not a user has permission to do something
func (e *Evaluator) Can(ctx context.Context, action Action, resource Resource) (bool, error) {
	if e.namespaceManager == nil {
		return false, nil
	}

	roles, err := e.namespaceManager.RolesFor(ctx, globalNamespace, e.namespace, e.user)
	if err != nil {
		return false, err
	}
	policies := globalRoleManager.getPolicies(roles...)
	for _, policy := range policies {
		actionMatch := policy.Action == ActionAll || policy.Action == action
		resourceMatch := policy.Resource == ResourceAll || pathMatch(resource.String(), policy.Resource.String())

		if actionMatch && resourceMatch {
			return true, nil
		}
	}
	return false, nil
}

// Policies returns policies for the user
func (e *Evaluator) Policies(ctx context.Context) ([]Policy, error) {
	if e.namespaceManager == nil {
		return nil, nil
	}
	roles, err := e.namespaceManager.RolesFor(ctx, globalNamespace, e.namespace, e.user)
	if err != nil {
		return nil, err
	}
	return globalRoleManager.getPolicies(roles...), nil
}

// pathMatch determines whether path matches the pattern, it matches
// on placeholders
func pathMatch(path, pattern string) bool {
	var regex *regexp.Regexp
	stored, ok := regexCache.Load(pattern)
	if ok {
		regex = stored.(*regexp.Regexp)
	} else {
		regexPattern := strings.Replace(pattern, "/*", "/.*", -1)
		for {
			if !strings.Contains(regexPattern, "/:") {
				break
			}
			regexPattern = tokenizer.ReplaceAllString(regexPattern, "$1[^/]+$2")
		}
		regex = regexp.MustCompile("^" + regexPattern + "$")
		regexCache.Store(pattern, regex)
	}

	return regex.MatchString(path)
}
