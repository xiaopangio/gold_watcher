package main

import (
	"testing"
)

func TestConfig_getHoldings(t *testing.T) {
	holdings, err := testApp.currentHoldings()
	if err != nil {
		t.Error("failed to get current holdings from database:", err)
	}
	if len(holdings) != 2 {
		t.Error("wrong holdings length returned;expected 2 but got", len(holdings))
	}
}
func TestConfig_getHoldingSlice(t *testing.T) {
	slice := testApp.getHoldingsSlice()
	if len(slice) != 3 {
		t.Error("wrong number of rows returned;expected 3 but got", len(slice))
	}
}
