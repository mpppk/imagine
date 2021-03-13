package testutil

import (
	"testing"

	fsa "github.com/mpppk/lorca-fsa"
)

type mockDispatcher struct {
	t          *testing.T
	gotActions []*fsa.Action
}

func NewMockDispatcher(t *testing.T) *mockDispatcher {
	return &mockDispatcher{t: t}
}

func (m *mockDispatcher) Dispatch(action *fsa.Action) error {
	m.gotActions = append(m.gotActions, action)
	return nil
}

func (m *mockDispatcher) Test(wantActions []*fsa.Action) {
	m.t.Helper()
	Diff(m.t, wantActions, m.gotActions)
}

func (m *mockDispatcher) Clean() {
	m.t.Helper()
	m.gotActions = nil
}

func (m *mockDispatcher) TestAndClean(wantActions []*fsa.Action) {
	m.Test(wantActions)
	m.Clean()
}
