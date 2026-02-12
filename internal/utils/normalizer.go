package utils

import "strings"

// NormalizeMunicipio normaliza o nome do município: capitaliza palavras com 3+ caracteres
func NormalizeMunicipio(municipio string) string {
	splittedMunicipio := strings.Split(municipio, " ")

	for i, word := range splittedMunicipio {
		if len(word) < 3 {
			continue
		}

		word = strings.ToLower(word)

		runes := []rune(word)
		runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]

		splittedMunicipio[i] = string(runes)
	}

	return strings.Join(splittedMunicipio, " ")
}
