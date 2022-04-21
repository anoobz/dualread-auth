package mockstore

import (
	"testing"
)

func CreateTestStore(t *testing.T) *MockStore {
	t.Helper()

	return NewMockStore()
}
