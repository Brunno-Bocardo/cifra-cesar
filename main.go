// POSSÍVEL OPÇÃO: https://api.dicionario-aberto.net/word/passaro 


package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// /cifrar 
type CifrarCesar struct {
	TextoClaro     string `json:"textoClaro"`
	Deslocamento   int `json:"deslocamento"`
}
type CifrarCesarResponse struct {
	TextoCifrado   string `json:"textoCifrado"`
}

// /decifrar
type DecifrarCesar struct {
	TextoCifrado   string `json:"textoCifrado"`
	Deslocamento   int `json:"deslocamento"`
}
type DecifrarCesarResponse struct {
	TextoClaro string `json:"textoClaro"`
}

// /decifrarForcaBruta
type ForcaBrutaCesar struct {
	TextoCifrado string `json:"textoCifrado"`
}
type ForcaBrutaCesarResponse struct {
	TextoClaro string `json:"textoClaro"`
}


func main() {
	http.HandleFunc("/cifrar", cifrarCesar)
	http.HandleFunc("/decifrar", decifrarCesar)
	http.HandleFunc("/decifrarForcaBruta", decifrarCesarForcaBruta)

	fmt.Println("Programa rodando em: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// ---------------- ENDPOINT Cifrar ----------------
func cifrarCesar(writer http.ResponseWriter, response *http.Request) {
	if response.Method != http.MethodPost {
		http.Error(writer, "Método inválido", http.StatusMethodNotAllowed)
		return
	}

	var mensagem CifrarCesar
	if err := json.NewDecoder(response.Body).Decode(&mensagem); err != nil {
		http.Error(writer, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validações básicas
	if mensagem.Deslocamento == 0 {
		http.Error(writer, "Deslocamento não pode ser zero", http.StatusBadRequest)
		return
	}
	if mensagem.TextoClaro == "" {
		http.Error(writer, "Texto claro não pode ser vazio", http.StatusBadRequest)
		return
	}

	// aplicar cifra
	textoCifrado, err := aplicarCifraCesar(mensagem.TextoClaro, mensagem.Deslocamento)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	resp := CifrarCesarResponse{TextoCifrado: textoCifrado}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(resp)
}



// ---------------- ENDPOINT Decifrar ----------------
func decifrarCesar(writer http.ResponseWriter, response *http.Request) {
	if response.Method != http.MethodPost {
		http.Error(writer, "Método inválido", http.StatusMethodNotAllowed)
		return
	}

	var mensagem DecifrarCesar
	if err := json.NewDecoder(response.Body).Decode(&mensagem); err != nil {
		http.Error(writer, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validações básicas
	if mensagem.Deslocamento == 0 {
		http.Error(writer, "Deslocamento não pode ser zero", http.StatusBadRequest)
		return
	}
	if mensagem.TextoCifrado == "" {
		http.Error(writer, "Texto cifrado não pode ser vazio", http.StatusBadRequest)
		return
	}

	// decifrar Cesar
	textoClaro, err := aplicarDescifraCesar(mensagem.TextoCifrado, mensagem.Deslocamento)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	resp := DecifrarCesarResponse{TextoClaro: textoClaro}
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(resp)
}



// ---------------- ENDPOINT Decifrar Força Bruta ----------------
func decifrarCesarForcaBruta(writer http.ResponseWriter, response *http.Request) {
	if response.Method != http.MethodPost {
		http.Error(writer, "Método inválido", http.StatusMethodNotAllowed)
		return
	}

	var mensagem ForcaBrutaCesar
	if err := json.NewDecoder(response.Body).Decode(&mensagem); err != nil {
		http.Error(writer, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validações básicas
	if mensagem.TextoCifrado == "" {
		http.Error(writer, "Texto cifrado não pode ser vazio", http.StatusBadRequest)
		return
	}

	// aplicar descifra forçada
	resultado := "teste funciona 3"

	resp := ForcaBrutaCesarResponse{TextoClaro: resultado}
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(resp)
}



// ---------------- FUNÇÕES AUXILIARES ----------------
func aplicarCifraCesar(texto string, deslocamento int) (string, error) {
	resultado := ""

	// percorremos cada caractere do texto
	for _, char := range texto {
		
		// aqui, iremos tratar apenas letras maiúsculas e minúsculas do alfabeto inglês
		// no caso, na linguaguem Go, os caracteres são representados por seus códigos Unicode (runes), ou ASCII
		// sendo assim, podemos fazer operações matemáticas com eles, considerando que 'a' = 97, 'b' = 98, ..., 'z' = 122
		// dessa forma, o deslocamento também pode ser aplicado matematicamente
		// Exemplo: para cifrar 'a' com deslocamento 3, fazemos: 'a' = 97 -> (97 - 97 + 3) % 26 + 97 = 100 -> 'd'

		if char >= 'a' && char <= 'z' {
			novoChar := ((int(char-'a') + deslocamento) % 26) + int('a')
			resultado += string(rune(novoChar))
		} else if char >= 'A' && char <= 'Z' {
			novoChar := ((int(char-'A') + deslocamento) % 26) + int('A')
			resultado += string(rune(novoChar))
		} else if char == ' ' {
			// o único caractere fora do alfabeto permitido é o espaço: ' '
			resultado += " "
		} else {
			// se utilizado um caractere fora de A-Z ou espaço, retorna um erro no endpoint
			return "", fmt.Errorf("caractere inválido: %c", char)
		}
	}

	return resultado, nil
}


func aplicarDescifraCesar(textoCifrado string, deslocamento int) (string, error) {
	resultado := ""

	// percorremos cada caractere do texto cifrado
	for _, char := range textoCifrado {
		
		// aqui, iremos tratar apenas letras maiúsculas e minúsculas do alfabeto inglês
		// no caso, na linguaguem Go, os caracteres são representados por seus códigos Unicode (runes), ou ASCII
		// sendo assim, podemos fazer operações matemáticas com eles, considerando que 'a' = 97, 'b' = 98, ..., 'z' = 122
		// dessa forma, o deslocamento também pode ser aplicado matematicamente
		// Exemplo: para decifrar 'd' com deslocamento 3, fazemos: 'd' = 100 -> (100 - 97 - 3 + 26) % 26 + 97 = 97 -> 'a'

		if char >= 'a' && char <= 'z' {
			novoChar := ((int(char-'a') - deslocamento + 26) % 26) + int('a')
			resultado += string(rune(novoChar))
		} else if char >= 'A' && char <= 'Z' {
			novoChar := ((int(char-'A') - deslocamento + 26) % 26) + int('A')
			resultado += string(rune(novoChar))
		} else if char == ' ' {
			// o único caractere fora do alfabeto permitido é o espaço: ' '
			resultado += " "
		} else {
			// se utilizado um caractere fora de A-Z ou espaço, retorna um erro no endpoint
			return "", fmt.Errorf("caractere inválido: %c", char)
		}
	}

	return resultado, nil
}