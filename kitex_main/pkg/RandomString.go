package pkg

import "math/rand"

func RandomString(count int) string {
	var str = "qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM"
	var b = make([]byte, count)
	for i := 0; i < count; i++ {
		intn := rand.Intn(len(str))
		b[i] = str[intn]
	}
	return string(b)
}
