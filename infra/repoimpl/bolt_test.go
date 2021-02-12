package repoimpl

import (
	"encoding/json"
	"testing"

	"github.com/mpppk/imagine/testutil"

	bolt "go.etcd.io/bbolt"
)

type TestData struct {
	ID  uint64
	Msg string
}

func (t *TestData) GetID() uint64 {
	return t.ID
}

func (t *TestData) SetID(id uint64) {
	t.ID = id
}

func newTestData(id uint64, msg string) *TestData {
	return &TestData{ID: id, Msg: msg}
}

func newTestDataFromJson(t *testing.T, data []byte) *TestData {
	var testData TestData
	if err := json.Unmarshal(data, &testData); err != nil {
		t.Errorf("failed to unmarshal json to TestData. data:%v", string(data))
	}
	return &testData
}

//func getTestData(t *testing.T, r *boltRepository, bucketNames []string, id uint64) *TestData {
//	bytes, ok, err := r.get(bucketNames, id)
//	if !ok {
//		t.Errorf("data not found. ID(%d)", id)
//		return nil
//	}
//
//	if err != nil {
//		t.Errorf("failed to get data by ID(%d) from bolt.", id)
//	}
//	return newTestDataFromJson(t, bytes)
//}

func listTestData(t *testing.T, r *boltRepository, bucketNames []string) (testDataList []boltData) {
	bytesList, err := r.list(bucketNames)
	if err != nil {
		t.Errorf("failed to list data from bolt.")
	}
	for _, bytes := range bytesList {
		testDataList = append(testDataList, newTestDataFromJson(t, bytes))
	}
	return
}

func addDataList(t *testing.T, b *boltRepository, bucketNames []string, dataList []boltData) {
	for _, data := range dataList {
		if _, err := b.addWithID(bucketNames, data); err != nil {
			t.Fatalf("failed to add data to bolt: %v", err)
		}
	}
}

func useBoltRepository(t *testing.T, bucketNames []string, f func(r *boltRepository)) {
	err := testutil.UseTempBoltDB(t, func(db *bolt.DB) error {
		r := newBoltRepository(db)
		if err := r.createBucketIfNotExist(bucketNames); err != nil {
			t.Fatalf("failed to create buckets: %v", bucketNames)
		}
		f(r)
		return nil
	})
	if err != nil {
		t.Fatalf("Unexpected DB Error: %v", err)
	}
}

func Test_itob(t *testing.T) {
	tests := []struct {
		name string
		v    uint64
	}{
		{v: 0},
		{v: 1},
		{v: 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := btoi(itob(tt.v)); result != tt.v {
				t.Errorf("provided: %v, got: %v", tt.v, result)
			}
		})
	}
}

func Test_boltRepository_addWithID(t *testing.T) {
	type args struct {
		bucketNames []string
		data        boltData
	}
	tests := []struct {
		name          string
		args          args
		existDataList []boltData
		want          uint64
		wantDataList  []boltData
		wantErr       bool
	}{
		{
			name: "add data to empty bucket",
			args: args{
				bucketNames: []string{"test-bucket"},
				data:        newTestData(99, "new-data"), // ID will be ignored
			},
			existDataList: []boltData{},
			want:          1,
			wantDataList:  []boltData{newTestData(1, "new-data")},
		},
		{
			name: "add data to bucket",
			args: args{
				bucketNames: []string{"test-bucket"},
				data:        newTestData(99, "new-data"), // ID will be ignored
			},
			existDataList: []boltData{newTestData(1, "data1")},
			want:          2,
			wantDataList:  []boltData{newTestData(1, "data1"), newTestData(2, "new-data")},
		},
		{
			name: "add data to nested bucket",
			args: args{
				bucketNames: []string{"test-bucket", "nested-bucket"},
				data:        newTestData(99, "new-data"), // ID will be ignored
			},
			existDataList: []boltData{newTestData(1, "data1")},
			want:          2,
			wantDataList:  []boltData{newTestData(1, "data1"), newTestData(2, "new-data")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useBoltRepository(t, tt.args.bucketNames, func(b *boltRepository) {
				addDataList(t, b, tt.args.bucketNames, tt.existDataList)

				id, err := b.addWithID(tt.args.bucketNames, tt.args.data)
				testutil.FailIfErrIsUnexpected(t, tt.wantErr, err)
				testutil.Diff(t, tt.want, id)
				dataList := listTestData(t, b, tt.args.bucketNames)
				testutil.Diff(t, tt.wantDataList, dataList)
			})
		})
	}
}
