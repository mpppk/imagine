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

func newTestDataFromMsg(msg string) *TestData {
	return &TestData{Msg: msg}
}

func newExistDataList(messages ...string) (dataList []boltData) {
	for _, message := range messages {
		dataList = append(dataList, newTestDataFromMsg(message))
	}
	return dataList
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
			existDataList: newExistDataList("data1"),
			want:          2,
			wantDataList:  []boltData{newTestData(1, "data1"), newTestData(2, "new-data")},
		},
		{
			name: "add data to nested bucket",
			args: args{
				bucketNames: []string{"test-bucket", "nested-bucket"},
				data:        newTestData(99, "new-data"), // ID will be ignored
			},
			existDataList: newExistDataList("data1"),
			want:          2,
			wantDataList:  []boltData{newTestData(1, "data1"), newTestData(2, "new-data")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useBoltRepository(t, tt.args.bucketNames, func(b *boltRepository) {
				addDataList(t, b, tt.args.bucketNames, tt.existDataList)

				id, err := b.addWithID(tt.args.bucketNames, tt.args.data)
				if testutil.FatalIfErrIsUnexpected(t, tt.wantErr, err) {
					return
				}
				testutil.Diff(t, tt.want, id)
				dataList := listTestData(t, b, tt.args.bucketNames)
				testutil.Diff(t, tt.wantDataList, dataList)
			})
		})
	}
}

func Test_boltRepository_updateByID(t *testing.T) {
	type args struct {
		bucketNames []string
		data        boltData
	}
	tests := []struct {
		name          string
		args          args
		existDataList []boltData
		wantDataList  []boltData
		wantErr       bool
	}{
		{
			name: "update data",
			args: args{
				bucketNames: []string{"test-bucket"},
				data:        newTestData(1, "new-data"),
			},
			existDataList: newExistDataList("old-data1", "data2"),
			wantDataList:  []boltData{newTestData(1, "new-data"), newTestData(2, "data2")},
		},
		{
			name: "update data on nested bucket",
			args: args{
				bucketNames: []string{"test-bucket", "nested-bucket"},
				data:        newTestData(1, "new-data"),
			},
			existDataList: newExistDataList("old-data1", "data2"),
			wantDataList:  []boltData{newTestData(1, "new-data"), newTestData(2, "data2")},
		},
		{
			name: "fail if data which have same ID does not exist",
			args: args{
				bucketNames: []string{"test-bucket"},
				data:        newTestData(2, "data2"),
			},
			existDataList: newExistDataList("data1"),
			wantErr:       true,
		},
		{
			name: "fail if data does not have ID",
			args: args{
				bucketNames: []string{"test-bucket"},
				data:        newTestData(0, "invalid-id-data"), // ID will be ignored
			},
			existDataList: newExistDataList("data1"),
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useBoltRepository(t, tt.args.bucketNames, func(b *boltRepository) {
				addDataList(t, b, tt.args.bucketNames, tt.existDataList)

				err := b.updateByID(tt.args.bucketNames, tt.args.data)
				if testutil.FatalIfErrIsUnexpected(t, tt.wantErr, err) {
					return
				}
				dataList := listTestData(t, b, tt.args.bucketNames)
				testutil.Diff(t, tt.wantDataList, dataList)
			})
		})
	}
}
