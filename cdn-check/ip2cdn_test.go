package cdn_check

import (
	"testing"
)

func Test_CdnCheck(t *testing.T) {
	cdnChecker := NewCdnChecker()
	result := cdnChecker.CdnCheckFilter([]string{"120.52.22.96", "8.8.8.8"})
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
}

func Test_CdnCheckOversea(t *testing.T) {
	cdnChecker := NewCdnChecker()
	res := cdnChecker.CdnCheckOversea([]string{"120.52.22.96", "8.8.8.8"})
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
	if _, ok := res["120.52.22.96"]; !ok {
		t.Fatalf("expected result to include known test IP")
	}
}
