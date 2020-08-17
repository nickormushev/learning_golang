package bitcoin

import (
	"errors"
	"fmt"
)

//Bitcoin type used in wallet
type Bitcoin float64

//Wallet is used to deposit and keep track of money
type Wallet struct {
	Money Bitcoin
}

//Deposit adds money to the wallet
func (w *Wallet) Deposit(amount Bitcoin) {
	w.Money += amount
}

//Balance returns the amount of money in the wallet
func (w *Wallet) Balance() Bitcoin {
	return w.Money
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%0.2f BTC", b)
}

//ErrInsufficientFunds tells the user they have insufficient funds to withdraw
var ErrInsufficientFunds = errors.New("Cannot withdraw! Insufficient funds")

//Withdraw removes bitcoins from wallet
func (w *Wallet) Withdraw(amount Bitcoin) error {
	if amount > w.Money {
		return ErrInsufficientFunds
	}

	w.Money -= amount
	return nil
}
