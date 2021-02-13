package testutil

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mpppk/imagine/infra"

	bolt "go.etcd.io/bbolt"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

func ExecuteCommand(t *testing.T, cmd *cobra.Command, command, in string) string {
	if in != "" {
		cmd.SetIn(strings.NewReader(in))
	}

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmdArgs := strings.Split(command, " ")
	cmd.SetArgs(cmdArgs)
	if err := cmd.Execute(); err != nil {
		t.Errorf("failed to execute box add command: %v", err)
	}

	return buf.String()
}

// Diff tests want and got are equals.
// If there are difference, fail test and print diff.
// See https://pkg.go.dev/github.com/google/go-cmp/cmp#Diff
func Diff(t *testing.T, want, got interface{}) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want +got)\n%s", diff)
	}
}

func DeepEqual(t *testing.T, want, got interface{}) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %#v, got: %#v", want, got)
	}
}

func FatalIfErrIsUnexpected(t *testing.T, wantErr bool, gotErr error) bool {
	if (gotErr != nil) != wantErr {
		t.Fatalf("error = %v, wantErr %v", gotErr, wantErr)
	}
	return wantErr
}

func NewTempDBFile(t *testing.T) (file *os.File, closeF func(), removeF func()) {
	file, err := ioutil.TempFile("", "imagine-test_*.imagine")
	if err != nil {
		t.Fatalf("failed to create temp db file: %v", err)
	}
	closeF = func() {
		t.Helper()
		if err := file.Close(); err != nil {
			t.Fatalf("failed to close temp file: %v", err)
		}
	}
	removeF = func() {
		t.Helper()
		if err := os.Remove(file.Name()); err != nil {
			t.Fatalf("failed to remove db file from %s: %v", file.Name(), err)
		}
	}
	return file, closeF, removeF
}

func NewTempBoltDB(t *testing.T) (db *bolt.DB, closeF func() error, removeF func()) {
	file, closeFile, removeFile := NewTempDBFile(t)
	defer closeFile()
	db, err := infra.NewBoltDB(file.Name())
	if err != nil {
		t.Fatalf("failed to create bolt DB: %v", err)
	}
	return db, db.Close, removeFile
}

func UseTempBoltDB(t *testing.T, f func(db *bolt.DB) error) (err error) {
	db, closeF, removeF := NewTempBoltDB(t)
	defer removeF()
	defer func() {
		err = closeF()
	}()
	if e := f(db); e != nil {
		return e
	}
	return
}
