package main

import "math/rand"

var towns = [][2]string{
	{"\033[92;1mLoresphread\033[0m",
		"\033[92mLoresphread to niewielkie miasteczko położone w zachodniej\n" +
			"części kraju. Jest tam umiarkowany klimat, ale śnieg prawie\n" +
			"nigdy nie pada. Na ogół społeczność jest przyjazna i skierowana\n" +
			"na poznawanie nowych ludzi, również spoza okolicy.\033[0m"},

	{"\033[94;1mDorthaven\033[0m",
		"\033[94mDorthaven to portowe miasto na wybrzeżu, znane z pięknych\n" +
			"widoków na ocean i dużej liczby rybaków. Wiatr od morza\n" +
			"często niesie świeżą bryzę. Mieszkańcy żyją w zgodzie z\n" +
			"naturą i słyną z organizowania barwnych świąt na plaży.\033[0m"},

	{"\033[93;1mHarvendale\033[0m",
		"\033[93mHarvendale leży w górach na północy i słynie z zapierających\n" +
			"dech w piersiach krajobrazów. Zimą miasto zmienia się w centrum\n" +
			"zimowych sportów, przyciągając turystów. Ludzie tam są hartowani\n" +
			"przez surowy klimat, ale zawsze gotowi pomóc przybyszom.\033[0m"},
}

func getTown() [2]string {
	return towns[rand.Intn(len(towns))]
}
