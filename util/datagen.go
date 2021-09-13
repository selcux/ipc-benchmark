package util

import (
	"crypto/rand"

	"github.com/pkg/errors"
)

func GenRandomBytes(size int) ([]byte, error) {
	blk := make([]byte, size)
	_, err := rand.Read(blk)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read the allocated memory space")
	}

	return blk, nil
}
