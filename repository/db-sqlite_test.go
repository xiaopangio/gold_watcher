package repository

import (
	"testing"
	"time"
)

func TestSQLiteRepository_Migrate(t *testing.T) {
	err := testRepo.Migrate()
	if err != nil {
		t.Error("Migrate failed:", err)
	}
}
func TestSQLiteRepository_InsertHolding(t *testing.T) {
	h := Holdings{
		Amount:        1,
		PurchaseDate:  time.Now(),
		PurchasePrice: 1000,
	}
	holding, err := testRepo.InsertHolding(h)
	if err != nil {
		t.Error("insert failed:", err)
	}

	if holding.ID <= 0 {
		t.Error("invalid id sent back:", holding.ID)
	}
}

func TestSQLiteRepository_AllHoldings(t *testing.T) {
	holdings, err := testRepo.AllHoldings()
	if err != nil {
		t.Error("get all failed:", err)
	}
	if len(holdings) != 1 {
		t.Error("wrong number of rows returned;expected 1 but got", len(holdings))
	}
}
func TestSQLiteRepository_GetHoldingByID(t *testing.T) {
	h, err := testRepo.GetHoldingByID(1)
	if err != nil {
		t.Error("get holding by id failed:", err)
	}
	if h.PurchasePrice != 1000 {
		t.Error("wrong purchase price returned;expected 1000 but got", h.PurchasePrice)
	}
	_, err = testRepo.GetHoldingByID(2)
	if err == nil {
		t.Error("get one returned value for non-existent id")
	}
}
func TestSQLiteRepository_UpdateHolding(t *testing.T) {
	h, err := testRepo.GetHoldingByID(1)
	if err != nil {
		t.Error(err)
	}
	h.PurchasePrice = 1001
	err = testRepo.UpdateHolding(1, *h)
	if err != nil {
		t.Error("update failed:", err)
	}
}
func TestSQLiteRepository_DeleteHolding(t *testing.T) {
	err := testRepo.DeleteHolding(1)
	if err != nil {
		t.Error("delete failed", err)
		if err != errDeleteFailed {
			t.Error("wrong error returned")
		}
	}
	err = testRepo.DeleteHolding(2)
	if err == nil {
		t.Error("no error when trying to delete no-existent record")
	}
}
