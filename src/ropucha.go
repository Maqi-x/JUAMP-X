package main

import (
	"math/rand"
	"strings"
	"time"
)

func handleRopucha() {
	Talk([][2]interface{}{
		{"Kasjer", "Witaj w najlepszym sklepie w całej dzielnicy!"},
		{"Kasjer", "Co podać?"},
		{"Ty", "Hm..."},
	}, map[string]string{
		"Kasjer": "magenta",
		"Ty":     "green",
	})

	list := func() {
		PrintLine("Dostępne produkty:")
		for name, info := range products {
			price := info.Price
			Println(Sprintf("%.2f - %s: %s", price, name, func() string {
				if price < wallet {
					return colorCodes["green"] + "✔" + colorCodes["reset"]
				} else {
					return colorCodes["red"] + "❌" + colorCodes["reset"]
				}
			}()))
		}
		Sep()
	}

	Println("Możesz uzyskać więcej info o danym przedmiocie wpisując `info {nazwa przedmiotu}`, ")
	Println("lub zakupić przedmiot wpisując `buy {nazwa przedmiotu}`. możesz też wpisać `list`, aby ponownie wyświetlić listę produktów")
	Println("Wpisz `exit`, aby wyjść ze sklepu")

	products = make(map[string]struct {
		Price       float64
		Description string
	})
	for key, v := range productS {
		products[normStr(key)] = v
	}
	// ---------------------------------------------- main loop ----------------------------------------- \\
	for {
		inp := Prompt(">>> ")

		if inp == "list" {
			list()
			continue
		} else if inp == "exit" {
			break
		} else if strings.HasPrefix(inp, "info") {
			if len(inp) < 6 {
				Println("Podaj prawidłową nazwę przedmiotu po komendzie `info`")
				continue
			}

			key := normStr(strings.TrimSpace(inp[5:]))
			v, exists := products[key]
			if !exists {
				Println("Nie ma takiego przedmiotu!")
				continue
			}
			if key == "gazeta" {
				v.Description = Sprintf(v.Description, TOWN)
			}
			PrintClr(Sprintf("Info o przedmiocie: %s", key), "cyan")
			PrintClr(Sprintf("Opis: %s", v.Description), "cyan")
			PrintClr(Sprintf("Cena: %.2f", v.Price), "cyan")
		} else if strings.HasPrefix(inp, "buy") {
			if len(inp) < 5 {
				Println("Podaj prawidłową nazwę przedmiotu po komendzie `buy`")
				continue
			}

			key := strings.TrimSpace(inp[4:])
			v, exists := products[normStr(key)]
			if !exists {
				Println("Nie ma takiego przedmiotu!")
				continue
			}

			ok := cPay(v.Price)
			if !ok {
				Println("Wygląda na to, że nie stać cię na ten produkt...")
				Println("Spróbuj ponownie później")
				if bank > v.Price {
					PrintClr("\033[1mWskazówka:", "blue")
					PrintClr("W banku masz wystarczającą ilość pieniędzy, by zakupić ten produkt!", "blue")
					PrintClr("Możesz udać się do banku i wypłacić pieniądze, a następnie tu wrócić", "blue")
				}
			} else {
				buy(key)
			}
		} else {
			Println("Nieznana komenda.")
		}
	}
	// ----------------------------------------------------------------------------------------------- \\
	Back()
}

func buy(name string) {
	name = normStr(name)
	switch name {
	case "bułka":
		hungryAdd(20)
		loading(2, "Jedzienie")
		chm("ROPUCHA")
		time.Sleep(500 * ms)
	case "pizza":
		hungryAdd(40)
		loading(3, "Jedzienie")
		chm("ROPUCHA")
		time.Sleep(500 * ms)
	case "gazeta":
		Sep()
		Println("Gazeta:")
		for _, line := range strings.Split(papers[rand.Intn(len(papers))], "\n") {
			Println(line)
			time.Sleep(300 * ms)
		}
		Sep()
		if tutStep[0] == 1 {
			tutStep = []int{2, 0}
		}
	}
}
