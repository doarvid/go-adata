package stock

import "testing"

func TestBatchEmpty(t *testing.T) {
    s := New()
    if out, _ := s.IndexCurrentBatch([]string{}, 0); len(out) != 0 { t.Fatalf("index batch not empty") }
    if out, _ := s.ConceptCurrentBatch([]string{}, 0); len(out) != 0 { t.Fatalf("concept batch not empty") }
}

