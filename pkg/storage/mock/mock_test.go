//
// KULADO INC. CONFIDENTIAL
// __________________
//
// [2014] - [2019] KULADO INC.
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of KULADO INC. and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to KULADO INC.
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from KULADO INC..
//

package mock_test

import (
	"testing"

	"github.com/onedomain/lastbackend/pkg/storage"
	"github.com/onedomain/lastbackend/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
)

func TestStorage_Get(t *testing.T) {
	stg, err := mock.New()
	assert.NoError(t, err, "storage initialize err")
	storage.StorageGetAssets(t, stg)
}

func TestStorage_List(t *testing.T) {
	stg, err := mock.New()
	assert.NoError(t, err, "storage initialize err")
	storage.StorageListAssets(t, stg)
}

func TestStorage_Map(t *testing.T) {
	stg, err := mock.New()
	assert.NoError(t, err, "storage initialize err")
	storage.StorageMapAssets(t, stg)
}

func TestStorage_Put(t *testing.T) {
	stg, err := mock.New()
	assert.NoError(t, err, "storage initialize err")
	storage.StoragePutAssets(t, stg)
}

func TestStorage_Set(t *testing.T) {
	stg, err := mock.New()
	assert.NoError(t, err, "storage initialize err")
	storage.StorageSetAssets(t, stg)
}

func TestStorage_Del(t *testing.T) {
	stg, err := mock.New()
	assert.NoError(t, err, "storage initialize err")
	storage.StorageDelAssets(t, stg)
}
