package testutil

import (
	"testing"

	fsa "github.com/mpppk/lorca-fsa"
)

type mockDispatcher struct {
	t           *testing.T
	wantActions []*fsa.Action
	gotActions  []*fsa.Action
}

func NewMockDispatcher(t *testing.T, expectedActions []*fsa.Action) *mockDispatcher {
	return &mockDispatcher{t: t, wantActions: expectedActions}
}

func (m *mockDispatcher) Dispatch(action *fsa.Action) error {
	m.gotActions = append(m.gotActions, action)
	return nil
}

func (m *mockDispatcher) Finish() {
	m.t.Helper()
	Diff(m.t, m.wantActions, m.gotActions)
}
