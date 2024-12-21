package main

func handleLobby() {
	Println("Witaj! gdzie checsz sie udać?")
	PrintC("1. Ropucha")
	PrintC("2. Wyjdź z gry")
	inp := Prompt(">>> ")
	switch inp {
	case "1":
		goTo("ROPUCHA")
	case "2":
		Exit()
	}
}
