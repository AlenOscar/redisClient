/**
 * Created by Alen on 2019-05-20 09:55
 */

package models

import (
	"testing"
)

func TestRedisCli(t *testing.T) {
	t.Logf("error = %v", RedisCli.SetHash("alen", "5", "100"))

	v, err := RedisCli.GetHash("alen", "5")
	t.Logf("alen[5] = %s, error = %v", v, err)

	var fields = []string{"1", "3", "4", "5"}
	results, err := RedisCli.GetHashMulti("alen", fields...)
	t.Logf("results = %s, error = %v", results, err)
}
