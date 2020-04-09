package memory

import (
	"testing"

	managerTest "github.com/andrewstucki/web-app-tools/go/security/testing"
)

func TestMemoryManager(t *testing.T) {
	managerTest.ManagerTest(t, NewNamespaceManager())
}
