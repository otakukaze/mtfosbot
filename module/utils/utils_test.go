package utils

import "testing"

func TestCalcPage(t *testing.T) {
	page := CalcPage(10, 1, 10)
	if page.Page != 1 {
		t.Error("Page Calc fail")
	}
	if page.Total != 1 {
		t.Error("Page Calc fail")
	}
	if page.Limit != 10 {
		t.Error("limit calc fail")
	}
	if page.Offset != 0 {
		t.Error("offset calc fail")
	}
}

func BenchmarkCalcPage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalcPage(10000, 30, 10)
	}
}
