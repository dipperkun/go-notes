package bank

import (
	"testing"
)

func TestWallet(t *testing.T) {
	wallet := Wallet{}
	wallet.Deposit(Bitcoin(10))
	got := wallet.Balance()
	wanted := Bitcoin(10)
	if got != wanted {
		t.Errorf("got %s wanted %s", got, wanted)
	}
}

func TestWallet2(t *testing.T) {
	check := func(t testing.TB, wallet Wallet, wanted Bitcoin) {
		t.Helper()
		got := wallet.Balance()

		if got != wanted {
			t.Errorf("got %s want %s", got, wanted)
		}
	}

	errchk := func(t testing.TB, got, wanted error) {
		t.Helper()
		if got == nil {
			t.Fatal("didn't get an error but wanted one")
		}

		if got != wanted {
			t.Errorf("got %q, wanted %q", got, wanted)
		}
	}

	noerr := func(t testing.TB, got error) {
		t.Helper()
		if got != nil {
			t.Fatal("got an error but didn't want one")
		}
	}

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		check(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(10))
		noerr(t, err)
		check(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{startingBalance}
		err := wallet.Withdraw(Bitcoin(100))
	
		errchk(t, err, ErrInsufficientFunds)
		check(t, wallet, startingBalance)
		
	})
}