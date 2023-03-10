package utility

func CaesarCipher(input string, delta int) string {
	var ciphered []rune
	for _, ch := range input {
		ciphered = append(ciphered, cipher(ch, delta))
	}
	return string(ciphered)
}

func cipher(r rune, delta int) rune {
	if r >= 'A' && r <= 'Z' {
		return rotate(r, 'A', delta)
	}
	if r >= 'a' && r <= 'z' {
		return rotate(r, 'a', delta)
	}
	return r
}

func rotate(r rune, base, delta int) rune {
	tmp := int(r) - base
	tmp = (tmp + delta) % 26 //Modulo operator to see what is remaining at end of the alphabet
	return rune(tmp + base)
}

//
// func rotate(s rune, delta int, key []rune) rune {
// 	idx := strings.IndexRune(string(key), s)
// 	if idx < 0 {
// 		panic("idx < 0")
// 	}
// 	idx = (idx + delta) % len(key) //Modulo operator to see what is remaining at end of the alphabet
// 	return key[idx]
// }
