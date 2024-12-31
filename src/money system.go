package main

import (
	"errors"
	"strconv"
	"strings"
)

func walletAdd(x float64) {
	wallet += x
}

func bankAdd(x float64) {
	bank += x
}

func walletD(x float64) {
	wallet -= x
}

func bankD(x float64) {
	bank -= x
}

func cWithdraw(x float64) error {
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

func cDeposit(x float64) error {
	if x <= 0 {
		return errors.New("kwota wpłaty musi być większa od 0")
	} else if wallet < x {
		return errors.New("Wygląda na to, że nie masz wystarczającej liczby pieniędzy w portfelu! no cóż, zdaża się...")
	} else {
		bank += x
		wallet -= x
		return nil
	}
}

func cPay(price float64) bool {
	if wallet >= price {
		wallet -= price
		return true
	}
	return false
}

func handleBank() {
	Talk([][2]interface{}{
		{"Bankier", "Witaj w banku!"},
		{"Bankier", "Co chcesz zrobić? wypłacić/wpłacic pieniądze, a może coś innego?"},
		{"Ty", "Hm..."},
	}, map[string]string{
		"Bankier": "magenta",
		"Ty":      "green",
	})
	Println("\r")
	help := func() {
		PrintS("Info")
		PrintColor("<bold>- {kwota}</bold> - wypłaca pieniądze z banku")
		PrintColor("<bold>+ {kwota}</bold> - wpłaca pieniądze do banku")
		PrintColor("<bold>help</bold> - wyświetla pomoc")
		PrintColor("<bold>info</bold> - wyświetla aktualny stan portfela i banku")
		PrintColor("<bold><x>exit</x></bold> - opuszcza bank")
	}
	help()
	for {
		inp := Prompt(">>> ")
		if strings.HasPrefix(inp, "-") {
			x, err := strconv.Atoi(strings.TrimSpace(inp[1:]))
			if err != nil {
				Println("Podany tekst to nie liczba")
				continue
			}
			err = cWithdraw(float64(x))
			if err != nil {
				ShowError(err.Error())
				continue
			}
			PrintS("Pomyślnie wypłacono kwote " + strconv.Itoa(x) + "zł")
		} else if strings.HasPrefix(inp, "+") {
			x, err := strconv.Atoi(strings.TrimSpace(inp[1:]))
			if err != nil {
				Println("Podany tekst to nie liczba")
				continue
			}
			err = cDeposit(float64(x))
			if err != nil {
				ShowError(err.Error())
				continue
			}
			PrintS("Pomyślnie wpłacono kwote " + strconv.Itoa(x) + "zł do banku")
			if x >= 200 {
				if tutStep[0] == 2 {
					tutStep = []int{3, 0}
				}
			}
		} else if inp == "help" {
			help()
		} else if inp == "exit" {
			break
		} else if inp == "info" {
			PrintColor("<bold><x>Portfel:</bold></x> " + strconv.FormatFloat(wallet, 'f', 2, 64) + "zł")
			PrintColor("<bold><x>Bank:</bold></x> " + strconv.FormatFloat(bank, 'f', 2, 64) + "zł")
		}
	}
	Back()
}
