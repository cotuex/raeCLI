package main

import (
	"github.com/gocolly/colly"
	"strings"
	"unicode"
)

// Article is the main structure in this program. It represents a RAE page
type Article struct {
	Title                string
	Variants             string
	Etimology            string
	Location             string
	Ortography           string
	ExtraInfo            []string
	Definitions          []Definition
	PossibleAlternatives []string
}

// Definition represents the multiple definitions of a word and the respective abbreviations
type Definition struct {
	Definition    string
	Abbreviations []string
}

// Scrape queries the RAE webpage and extracts the information
func Scrape(word string) Article {
	output := Article{}
	c := colly.NewCollector()

	c.OnHTML("div#resultados > .item-list", func(e *colly.HTMLElement) {
		alternatives := map[string]struct{}{}

		for _, d := range e.ChildAttrs("div > a", "data-eti") {
			alternatives[d] = struct{}{}
		}

		for k := range alternatives {
			output.PossibleAlternatives = append(output.PossibleAlternatives, k)
		}
	})

	c.OnHTML("div#resultados > article:nth-of-type(1)", func(e *colly.HTMLElement) {
		// Titulo
		output.Title = e.ChildText("header")

		// Definiciones
		for _, i := range e.ChildAttrs("p[class^='j']", "id") {
			char := e.ChildAttrs("p[class^='j'][id='"+i+"'] > abbr", "title")

			def := e.ChildText("p[class^='j'][id='" + i + "']")
			def = strings.TrimLeftFunc(def, func(r rune) bool {
				return unicode.IsDigit(r) || r == '.' || r == ' '
			})

			output.Definitions = append(output.Definitions, Definition{Definition: def, Abbreviations: char})
		}

		// Variantes
		output.Variants = e.ChildText("p.n1 > a")
		output.Variants = strings.TrimFunc(output.Variants, func(r rune) bool {
			return unicode.IsPunct(r)
		})

		// Informacion extra
		extraInfo := e.ChildTexts("p.n2, p.n3, p.n4, p.n5")
		for _, d := range extraInfo {
			if strings.HasPrefix(d, "Del ") || strings.HasPrefix(d, "De ") {
				output.Etimology = d
			} else if strings.HasPrefix(d, "Loc.") {
				output.Location = d
			} else if strings.HasPrefix(d, "Escr.") {
				output.Ortography = d
			} else {
				output.ExtraInfo = append(output.ExtraInfo, d)
			}

		}
	})

	c.Visit("https://dle.rae.es/" + word)
	return output
}
