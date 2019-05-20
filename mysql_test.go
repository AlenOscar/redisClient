/**
 * Created by Alen on 2019-05-20 14:38
 */

package models

import "testing"

func TestDbWorker_Insert(t *testing.T) {
	Db.Insert()
}

func TestDbWorker_Delete(t *testing.T) {
	Db.Delete()
}

func TestDbWorker_Update(t *testing.T) {
	Db.Update()
}

func TestDbWorker_Query(t *testing.T) {
	Db.Query()
}
