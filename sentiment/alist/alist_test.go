package alist

import (
	"context"
	"fmt"
	"testing"
)

func TestGetAListInfo(t *testing.T) {
	c := New()
	rows, _ := c.Details(context.Background(), "600297", "2024-07-12", 0)
	if len(rows) == 0 {
		t.Fatalf("alist info not empty")
	}
	if len(rows) == 0 {
		t.Fatalf("stock info not empty")
	}
	for _, row := range rows {
		fmt.Printf("row: %+v\n", row)
	}
}

func TestListAListDaily(t *testing.T) {
	c := New()
	rows, _ := c.Daily(context.Background(), "2024-07-12", 0)
	if len(rows) == 0 {
		t.Fatalf("alist info not empty")
	}
	if len(rows) == 0 {
		t.Fatalf("alist daily not empty")
	}
	for _, row := range rows {
		fmt.Printf("row: %+v\n", row)
	}
}

