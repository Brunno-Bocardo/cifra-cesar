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
	// textoCifrado := cifrarVernam(mensagem.TextoClaro, mensagem.Chave)
	textoCifrado := "teste funciona"

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

	// aplicar descifra
	resultado := "teste funciona 2"

	resp := DecifrarCesarResponse{TextoClaro: resultado}
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

