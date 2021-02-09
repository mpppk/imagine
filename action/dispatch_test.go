package action

import (
	"testing"

	"github.com/mpppk/imagine/testutil"

	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
)

type mockDispatcher struct {
	t           *testing.T
	wantActions []*fsa.Action
	gotActions  []*fsa.Action
}

func newMockDispatcher(t *testing.T, expectedActions []*fsa.Action) *mockDispatcher {
	return &mockDispatcher{t: t, wantActions: expectedActions}
}

func (m *mockDispatcher) Dispatch(action *fsa.Action) error {
	m.gotActions = append(m.gotActions, action)
	return nil
}

func (m *mockDispatcher) Finish() {
	m.t.Helper()
	testutil.Diff(m.t, m.wantActions, m.gotActions)
}
