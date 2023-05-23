package buildkitskipper

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

var errInvalidHash = fmt.Errorf("invalid sha1 hash")

func NewLocal(path string) (*LocalBuildkitSkipper, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %w", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("builds"))
		if err != nil {
			return fmt.Errorf("could not create builds bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update database, %w", err)
	}
	return &LocalBuildkitSkipper{
		db: db,
	}, nil
}

type LocalBuildkitSkipper struct {
	db *bolt.DB
}

func (lbks *LocalBuildkitSkipper) Add(ctx context.Context, data []byte) error {
	if len(data) != sha1.Size {
		return errInvalidHash
	}
	return lbks.db.Update(func(tx *bolt.Tx) error {
		payload := []byte(time.Now().String()) // could be serialized into a structure; however LocalBuildkitSkipper is only meant for dev/testing
		err := tx.Bucket([]byte("builds")).Put(data, payload)
		if err != nil {
			return fmt.Errorf("could not set config: %w", err)
		}
		return nil
	})
}

func (lbks *LocalBuildkitSkipper) Exists(ctx context.Context, data []byte) (bool, error) {
	if len(data) != sha1.Size {
		return false, errInvalidHash
	}
	var found bool
	err := lbks.db.View(func(tx *bolt.Tx) error {
		payload := tx.Bucket([]byte("builds")).Get(data)
		if payload != nil {
			found = true
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return found, nil
}
