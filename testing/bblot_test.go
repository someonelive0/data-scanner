package testing

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"go.etcd.io/bbolt"
)

// RUN: go test -v ./testing/
func Test_bblot(t *testing.T) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bbolt.Open("mybblot.db", 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		fmt.Println(bucket)

		if err = bucket.Put([]byte("answer"), []byte("42")); err != nil {
			return err
		}

		if err = bucket.Put([]byte("zero"), []byte("")); err != nil {
			return err
		}

		// if err = tx.DeleteBucket([]byte("MyBucket")); err != nil {
		// 	return err
		// }

		return nil
	})

	db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		v := b.Get([]byte("noexists"))
		fmt.Println(reflect.DeepEqual(v, nil)) // false
		fmt.Println(v == nil)                  // true

		v = b.Get([]byte("zero"))
		fmt.Println(reflect.DeepEqual(v, nil)) // false
		fmt.Println(v == nil)                  // true
		return nil
	})

	db.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("MyBucket"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})
}
