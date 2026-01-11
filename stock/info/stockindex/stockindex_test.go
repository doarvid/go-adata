package stockindex

import (
	"testing"
)

func TestLoadIndexCodeRelTHS(t *testing.T) {
	m, err := LoadIndexCodeRelTHS()
	if err != nil {
		t.Fatalf("LoadIndexCodeRelTHS err: %v", err)
	}
	if m == nil {
		t.Fatalf("LoadIndexCodeRelTHS nil")
	}
}

func TestHTTPClientConfigSetters(t *testing.T) {
	SetHTTPClientConfig(HTTPClientConfig{UserAgent: "test/ua"})
	SetHTTPClient(nil)
}

func TestAllIndexCodeEast_Smoke(t *testing.T) {
	rows, err := AllIndexCodeEast()
	if err != nil || len(rows) == 0 {
		t.Skipf("AllIndexCodeEast skip: err=%v len=%d", err, len(rows))
		return
	}
}

func TestIndexConstituentBaidu_Smoke(t *testing.T) {
	rows, err := IndexConstituentBaidu("000001")
	if err != nil || len(rows) == 0 {
		t.Skipf("IndexConstituentBaidu skip: err=%v len=%d", err, len(rows))
		return
	}
}
