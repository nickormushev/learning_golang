package bitcoin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWallet(t *testing.T) {

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}

		wallet.Deposit(10)

		assertBalance(t, Bitcoin(10.0), wallet)
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{20}

		e := wallet.Withdraw(10)

		assertBalance(t, Bitcoin(10), wallet)
		assert.Nil(t, e)
	})

	t.Run("Withdraw with no money in wallet", func(t *testing.T) {
		wallet := Wallet{}

		e := wallet.Withdraw(10)

		assertError(t, e, ErrInsufficientFunds)
	})
}

func assertBalance(t *testing.T, amount Bitcoin, wallet Wallet) {
	t.Helper()
	got := wallet.Balance()

	if got != amount {
		t.Errorf("Got %s want %s", got, amount)
	}
}

func assertError(t *testing.T, got error, want error) {
	t.Helper()
	assert.Error(t, got)

	if got != want {
		t.Errorf("Wanted this error: %q, but got this one instead: %q", want, got)
	}
}
