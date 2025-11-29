package conceptflow

import "testing"

func TestBuildParams(t *testing.T) {
    fid, fields := buildParams(1)
    if fid != "f62" { t.Fatalf("fid 1") }
    if fields == "" { t.Fatalf("fields 1 empty") }
    fid, fields = buildParams(5)
    if fid != "f164" { t.Fatalf("fid 5") }
    fid, fields = buildParams(10)
    if fid != "f174" { t.Fatalf("fid 10") }
}

