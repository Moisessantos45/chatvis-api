package pkg

import (
	"fmt"
	"math/big"
	"strings"
)

const BASE58 = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// UuidToBase58 convierte un UUID a string corto Base58
func UuidToBase58(uuid string) string {
	// 1. Limpiar guiones
	hexStr := strings.ReplaceAll(uuid, "-", "")

	// 2. Convertir hex a BigInt (el "0x" es clave para que Go entienda la base 16)
	n := new(big.Int)
	n.SetString(hexStr, 16)

	// 3. Lógica de división por 58
	var result string
	base := big.NewInt(58)
	mod := new(big.Int)

	if n.Sign() == 0 {
		return string(BASE58[0])
	}

	for n.Sign() > 0 {
		n.DivMod(n, base, mod)
		result = string(BASE58[mod.Int64()]) + result
	}

	return result
}

// Base58ToUuid convierte el string corto de vuelta al UUID original
func Base58ToUuid(short string) string {
	n := big.NewInt(0)
	base := big.NewInt(58)

	for _, char := range short {
		index := strings.IndexRune(BASE58, char)
		n.Mul(n, base)
		n.Add(n, big.NewInt(int64(index)))
	}

	// 4. Formatear a hex con 32 caracteres (relleno con ceros a la izquierda)
	hex := fmt.Sprintf("%032x", n)

	// 5. Reconstruir formato UUID: 8-4-4-4-12
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hex[:8], hex[8:12], hex[12:16], hex[16:20], hex[20:])
}
