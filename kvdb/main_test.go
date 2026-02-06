package main_test

import (
	"fmt"
	"sync"
	"testing"

	"kvdb/pkg"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	db := pkg.NewKVDB()

	t.Run("CreateCollection", func(t *testing.T) {
		assert.NoError(t, db.CreateCollection("collection1"))
		collection, err := db.GetCollection("collection1")
		assert.NoError(t, err)
		assert.NotNil(t, collection)
	})

	t.Run("CRUD operations", func(t *testing.T) {
		collection, err := db.GetCollection("collection1")
		assert.NoError(t, err)
		assert.NoError(t, collection.CreateFilterKey("filterKey"))
		assert.NoError(t, collection.CreateFilterKey("filterKey2"))

		t.Run("Set", func(t *testing.T) {
			assert.NoError(t, collection.Set("key", map[string]any{
				"data":       "value",
				"filterKey":  "filterValue",
				"filterKey2": "filterValue2",
			}))

			assert.NoError(t, collection.Set("key2", map[string]any{
				"data":       "value",
				"filterKey":  "filterValue",
				"filterKey2": "filterValue2",
			}))

		})

		t.Run("Scan", func(t *testing.T) {
			scan, _ := collection.Scan()
			assert.Len(t, scan, 2)
		})

		t.Run("FindOne", func(t *testing.T) {
			findOne, _ := collection.FindOne(pkg.Key("key"))
			assert.Equal(t, findOne, `{"data":"value","filterKey":"filterValue","filterKey2":"filterValue2"}`)
		})

		t.Run("FindAll", func(t *testing.T) {
			findAll, err := collection.FindAll(pkg.NewFilterKey("filterKey"), "filterValue")
			assert.NoError(t, err)
			assert.Len(t, findAll, 2)
		})

		t.Run("FindAll", func(t *testing.T) {
			findAll, err := collection.FindAll(pkg.NewFilterKey("filterKey2"), "filterValue2")
			assert.NoError(t, err)
			assert.Len(t, findAll, 2)
		})
	})
}

func TestConcurrentWrites(t *testing.T) {
	db := pkg.NewKVDB()
	assert.NoError(t, db.CreateCollection("concurrent_test"))
	collection, err := db.GetCollection("concurrent_test")
	assert.NoError(t, err)
	assert.NoError(t, collection.CreateFilterKey("category"))

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", idx)
			err := collection.Set(key, map[string]any{
				"data":     fmt.Sprintf("value_%d", idx),
				"category": "test",
			})
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Verificar se todos os dados foram escritos
	scan, err := collection.Scan()
	assert.NoError(t, err)
	assert.Len(t, scan, numGoroutines)
}

func TestConcurrentReads(t *testing.T) {
	db := pkg.NewKVDB()
	assert.NoError(t, db.CreateCollection("concurrent_read_test"))
	collection, err := db.GetCollection("concurrent_read_test")
	assert.NoError(t, err)

	// Inserir dados primeiro
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key_%d", i)
		err := collection.Set(key, map[string]any{
			"data": fmt.Sprintf("value_%d", i),
		})
		assert.NoError(t, err)
	}

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			key := pkg.Key(fmt.Sprintf("key_%d", idx%10))
			value, err := collection.FindOne(key)
			assert.NoError(t, err)
			assert.NotEmpty(t, value)
		}(i)
	}

	wg.Wait()
}

func TestConcurrentReadsAndWrites(t *testing.T) {
	db := pkg.NewKVDB()
	assert.NoError(t, db.CreateCollection("concurrent_rw_test"))
	collection, err := db.GetCollection("concurrent_rw_test")
	assert.NoError(t, err)
	assert.NoError(t, collection.CreateFilterKey("type"))

	// Inserir alguns dados iniciais
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("initial_key_%d", i)
		err := collection.Set(key, map[string]any{
			"data": fmt.Sprintf("initial_value_%d", i),
			"type": "initial",
		})
		assert.NoError(t, err)
	}

	const numWriters = 50
	const numReaders = 100
	var wg sync.WaitGroup
	wg.Add(numWriters + numReaders)

	// Writers
	for i := 0; i < numWriters; i++ {
		go func(idx int) {
			defer wg.Done()
			key := fmt.Sprintf("new_key_%d", idx)
			err := collection.Set(key, map[string]any{
				"data": fmt.Sprintf("new_value_%d", idx),
				"type": "new",
			})
			assert.NoError(t, err)
		}(i)
	}

	// Readers
	for i := 0; i < numReaders; i++ {
		go func(idx int) {
			defer wg.Done()
			key := pkg.Key(fmt.Sprintf("initial_key_%d", idx%10))
			_, _ = collection.FindOne(key)
			_, _ = collection.Scan()
		}(i)
	}

	wg.Wait()

	// Verificar se todos os dados foram escritos
	scan, err := collection.Scan()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(scan), 10+numWriters)
}

func TestConcurrentCollectionCreation(t *testing.T) {
	db := pkg.NewKVDB()

	const numGoroutines = 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	var successCount int
	var mu sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			// Tentar criar a mesma collection de várias goroutines
			err := db.CreateCollection("shared_collection")
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Apenas uma goroutine deve ter sucesso
	assert.Equal(t, 1, successCount)

	// A collection deve existir
	collection, err := db.GetCollection("shared_collection")
	assert.NoError(t, err)
	assert.NotNil(t, collection)
}

func TestConcurrentFilterKeyCreation(t *testing.T) {
	db := pkg.NewKVDB()
	assert.NoError(t, db.CreateCollection("filter_key_test"))
	collection, err := db.GetCollection("filter_key_test")
	assert.NoError(t, err)

	const numGoroutines = 20
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	var successCount int
	var mu sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			// Tentar criar a mesma filter key de várias goroutines
			err := collection.CreateFilterKey("shared_filter_key")
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Apenas uma goroutine deve ter sucesso
	assert.Equal(t, 1, successCount)
}

func TestConcurrentFindAll(t *testing.T) {
	db := pkg.NewKVDB()
	assert.NoError(t, db.CreateCollection("findall_test"))
	collection, err := db.GetCollection("findall_test")
	assert.NoError(t, err)
	assert.NoError(t, collection.CreateFilterKey("category"))

	// Inserir dados
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("key_%d", i)
		err := collection.Set(key, map[string]any{
			"data":     fmt.Sprintf("value_%d", i),
			"category": "test_category",
		})
		assert.NoError(t, err)
	}

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			results, err := collection.FindAll(pkg.NewFilterKey("category"), "test_category")
			assert.NoError(t, err)
			assert.Len(t, results, 50)
		}()
	}

	wg.Wait()
}
