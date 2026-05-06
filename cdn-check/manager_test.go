package cdn_check

import (
	"reflect"
	"testing"
)

func TestChunkIPs(t *testing.T) {
	ips := []string{"1", "2", "3", "4", "5"}
	got := chunkIPs(ips, 2)
	want := [][]string{{"1", "2"}, {"3", "4"}, {"5"}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("chunkIPs() = %#v, want %#v", got, want)
	}
}
