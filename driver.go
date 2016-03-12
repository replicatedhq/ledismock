package ledismock

import (
	"fmt"

	"github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
	"github.com/siddontang/ledisdb/store/driver"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type Iterator struct {
	it *iterator.Iterator
}

func (it Iterator) Close() error {
	return nil
}

func (it Iterator) First() {

}

func (it Iterator) Last() {

}

func (it Iterator) Seek(key []byte) {

}

func (it *Iterator) Next() {

}

func (it *Iterator) Prev() {
}

func (it *Iterator) Valid() bool {
	return false
}

func (it *Iterator) Key() []byte {
	return nil
}
func (it *Iterator) Value() []byte {
	return nil
}

// WriteBatch is a mock implementation of the ledis IWriteBatch interface.
type WriteBatch struct {
}

func (b WriteBatch) Close() {

}

func (b WriteBatch) Commit() error {
	return nil
}

func (b WriteBatch) Data() []byte {
	return nil
}

func (b WriteBatch) Delete(key []byte) {

}

func (b WriteBatch) Put(key []byte, value []byte) {

}

func (b WriteBatch) Rollback() error {
	return nil
}

func (b WriteBatch) SyncCommit() error {
	return nil
}

// MockDB is a mock implementation of a storage driver for ledis.
type MockDB struct {
	mock *LedisMock
}

func (db MockDB) Begin() (driver.Tx, error) {
	fmt.Printf("\n\nBegin\n\n")
	return nil, nil
}

func (db MockDB) Close() error {
	fmt.Printf("\n\nClose\n\n")
	return nil
}

func (db MockDB) Compact() error {
	fmt.Printf("\n\nCompact\n\n")
	return nil
}

func (db MockDB) Delete(key []byte) error {
	fmt.Printf("\n\nDelete\n\n")
	return nil
}

func (db MockDB) Get(key []byte) ([]byte, error) {
	db.mock.receivedGet(key)
	return nil, nil
}

func (db MockDB) Put(key, value []byte) error {
	fmt.Printf("\n\nPut\n\n")
	return nil
}

func (db MockDB) SyncPut(key []byte, value []byte) error {
	fmt.Printf("\n\nSyncPut\n\n")
	return nil
}

func (db MockDB) SyncDelete(key []byte) error {
	fmt.Printf("\n\nSyncDelete\n\n")
	return nil
}

func (db MockDB) NewIterator() driver.IIterator {
	i := &Iterator{}
	return i
}

func (db MockDB) NewSnapshot() (driver.ISnapshot, error) {
	fmt.Printf("\n\nNewSnapshot\n\n")
	return nil, nil
}

func (db MockDB) NewWriteBatch() driver.IWriteBatch {
	b := WriteBatch{}
	return b
}

// Store is an implementation of a ledis store used for mocking.
type Store struct {
	DBName string
	Mock   *LedisMock
}

// String implements the required interface on the store object.
func (s Store) String() string {
	return s.DBName
}

// Open implements the required interface function to create the database object.
func (s Store) Open(path string, cfg *config.Config) (driver.IDB, error) {
	mockDB := MockDB{mock: s.Mock}
	return &mockDB, nil
}

// Repair is a required function on the store interface, but not implemented for ths mock.
func (s Store) Repair(path string, cfg *config.Config) error {
	// Not implemented
	return nil
}

// New creates a ledismock database connection
// and a mock to manage expectations.
func New() (*ledis.DB, *LedisMock, error) {
	// Create a mock object
	mock := LedisMock{}

	driver.Register(Store{DBName: "mock", Mock: &mock})

	mockcfg := config.NewConfigDefault()
	mockcfg.DBName = "mock"

	conn, err := ledis.Open(mockcfg)
	if err != nil {
		return nil, nil, err
	}

	db, err := conn.Select(0)
	if err != nil {
		return nil, nil, err
	}

	return db, &mock, nil
}
