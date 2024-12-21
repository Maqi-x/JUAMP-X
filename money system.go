package main

import (
    "errors"
)

func walletAdd(x int64) {
    wallet += x
}

func bankAdd(x int64) {
    bank += x
}

func walletD(x int64) {
    wallet -= x
}

func bankD(x int64) {
    bank -= x
}

func cWithdraw(x int64) error {
	if x <= 0 {
		return errors.New("kwota wypłaty musi być większa od 0")
	} else if bank < x {
		return errors.New("Wygląda na to, że nie masz wystarczającej liczby pieniędzy w banku! no cóż, zdaża się...")
	} else {
	    bank -= x
	    wallet += x

	    return nil
	}
}

func cDeposit(x int64) error {
	if x <= 0 {
		return errors.New("kwota wpłaty musi być większa od 0")
	} else if wallet < x {
		return errors.New("Wygląda na to, że nie masz wystarczającej liczby pieniędzy w portfelu! no cóż, zdaża się...")
	} else {
	    bank += int64(x)
	    wallet -= int64(x)
	    return nil
	}
}
