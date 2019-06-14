package db

import (
	"github.com/danskeren/database/kv"
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
)

var BadgerDB kv.KV

func init() {
	opts := badger.DefaultOptions
	opts.Dir = "./badger.db"
	opts.ValueDir = "./badger.db"
	opts.SyncWrites = true
	opts.ValueLogLoadingMode = options.FileIO
	var err error
	BadgerDB, err = kv.Open(opts)
	if err != nil {
		panic(err)
	}
}
