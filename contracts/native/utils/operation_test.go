package utils

import "testing"

func TestGetUint64Bytes(t *testing.T) {
	var a uint64 = 0x789
	t.Logf("uint64 %d bytes %x", a, GetUint64Bytes(a))
}
