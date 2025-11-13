// POSSÍVEL OPÇÃO: https://api.dicionario-aberto.net/word/passaro 


package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"bufio" //deixa a leitura do arquivo mais rapida
	"os" // Para abrir o arquivo e para sair do programa 
)
//guardar todas as palavras sem acento do dicionario 
var dicionarioLocal map[string]bool

//-----------Função para carregar o Dicionário Local------------

func carregarDicionarioLocal(caminhoDicionario string) (map[string]bool, error){
	mapa := make(map[string]bool)

	arquivo, err := os.Open(caminhoDicionario)
	if err != nil {
		return nil, fmt.Errorf("nao foi possivel abrir o arquivo %s: %w", caminhoDicionario, err)
	}
	defer arquivo.Close()

	// Lê o arquivo linha por linha
	scanner := bufio.NewScanner(arquivo)
	for scanner.Scan() {
		// Pega a palavra, remove espaços e converte para minúsculas
		palavra := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if palavra != "" {
			mapa[palavra] = true 
		}
	}
	//verifica se teve algum erro na leitura
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo %s: %w", caminhoDicionario, err)
	}

	return mapa, nil
} 
	

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

	var err error

	//define o nome do arquivo do dicionário
	arquivoDicionario := "palavras_sem_acento.txt"

	dicionarioLocal, err = carregarDicionarioLocal(arquivoDicionario)
	if err != nil {
		// Se o dicionário não puder ser carregado, o programa não deve rodar.
		fmt.Fprintf(os.Stderr, "Erro fatal ao carregar dicionario local (%s): %v\n", arquivoDicionario, err)
	}

	// Confirmar que funcionou
	fmt.Printf("Dicionario local carregado com %d palavras.\n", len(dicionarioLocal))

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
	if mensagem.TextoCifrado == "" {
		http.Error(writer, "Texto cifrado não pode ser vazio", http.StatusBadRequest)
		return
	}

	textoClaro, err := tentarForcaBruta(mensagem.TextoCifrado)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	resp := ForcaBrutaCesarResponse{TextoClaro: textoClaro}
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(resp)
}



// ---------------- FUNÇÕES AUXILIARES ----------------
func removerAcentos(texto string) string {
    acento := map[rune]rune{
        'á': 'a', 'à': 'a', 'ã': 'a', 'â': 'a', 'ä': 'a',
        'é': 'e', 'è': 'e', 'ê': 'e', 'ë': 'e',
        'í': 'i', 'ì': 'i', 'î': 'i', 'ï': 'i',
        'ó': 'o', 'ò': 'o', 'õ': 'o', 'ô': 'o', 'ö': 'o',
        'ú': 'u', 'ù': 'u', 'û': 'u', 'ü': 'u',
        'ç': 'c',
        'Á': 'A', 'À': 'A', 'Ã': 'A', 'Â': 'A', 'Ä': 'A',
        'É': 'E', 'È': 'E', 'Ê': 'E', 'Ë': 'E',
        'Í': 'I', 'Ì': 'I', 'Î': 'I', 'Ï': 'I',
        'Ó': 'O', 'Ò': 'O', 'Õ': 'O', 'Ô': 'O', 'Ö': 'O',
        'Ú': 'U', 'Ù': 'U', 'Û': 'U', 'Ü': 'U',
        'Ç': 'C',
    }
    resultado := ""
    for _, r := range texto {
        if novo, ok := acento[r]; ok {
            resultado += string(novo)
        } else {
            resultado += string(r)
        }
    }
    return resultado
}


func aplicarCifraCesar(texto string, deslocamento int) (string, error) {
	resultado := ""
	texto = removerAcentos(texto)

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
		} else if int(char) >= 32 || int(char) <= 57 {
			// consideramos alguns caracteres fora do alfabeto permitido, como o espaço, números e pontuações
			// nesse caso, apenas os adicionamos ao resultado sem alteração
			// fonte: https://www.ime.usp.br/~kellyrb/mac2166_2015/tabela_ascii.html
			resultado += string(char)
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
		} else if int(char) >= 32 || int(char) <= 57 {
			// consideramos alguns caracteres fora do alfabeto permitido, como o espaço, números e pontuações
			// nesse caso, apenas os adicionamos ao resultado sem alteração
			// fonte: https://www.ime.usp.br/~kellyrb/mac2166_2015/tabela_ascii.html
			resultado += string(char)
		} else {
			// se utilizado um caractere fora de A-Z ou espaço, retorna um erro no endpoint
			return "", fmt.Errorf("caractere inválido: %c", char)
		}
	}

	return resultado, nil
}


// ---------------- FUNÇÕES DE FORÇA BRUTA (simples) ----------------
func tentarForcaBruta(textoCifrado string) (string, error) {
	client := &http.Client{Timeout: 3 * time.Second}

	// tentar deslocamentos 1..25
	for desloc := 1; desloc <= 25; desloc++ {
		decifrado, err := aplicarDescifraCesar(textoCifrado, desloc)
		if err != nil {
			continue
		}
		// extrair palavras (apenas A-Z/a-z)
		words := extrairPalavras(decifrado)
		if len(words) == 0 {
			continue
		}

		// verificar cada palavra no dicionário (apenas status 200)
		palavrasCorretas := 0
		for _, w := range words {
			if w == "" {
				continue
			}
			// dicionário não lida com acento no seu fluxo atual; removemos acentos para a consulta
			wSemAcento := strings.ToLower(removerAcentos(w))
			ok := existeNoDicionario(client, wSemAcento)

			// Se a API falhar, tenta o dicionário local
			if !ok {
				if _, existeLocal := dicionarioLocal[wSemAcento]; existeLocal {
					ok = true 
				}
			}

			// aqui vemos se a palavra existe no dicionário
			if !ok {
				// no caso, se ela não existir, podemos consultar um arquivo exerno de palavras
				continue
			} else {
				palavrasCorretas++
			}
		}
		
		fmt.Println("- palavras:", len(words))
		fmt.Println("- palavras corretas:", palavrasCorretas)
		if len(words) == 2 && palavrasCorretas == 2 {
			return decifrado, nil
		}
		if len(words) != 2 && palavrasCorretas >= ((len(words) + 1) / 2) {
			return decifrado, nil
		}
	}

	return "", fmt.Errorf("não foi possível encontrar uma chave válida")
}


var rePalavra = regexp.MustCompile(`[A-Za-z]+`)
func extrairPalavras(s string) []string {
	return rePalavra.FindAllString(s, -1)
}


func existeNoDicionario(client *http.Client, palavra string) bool {
	if palavra == "" {
		return false
	}
	u := "https://api.dicionario-aberto.net/word/" + url.PathEscape(palavra)
	fmt.Println("Consultando dicionário para palavra:", palavra)
	fmt.Println("Consultando URL:", u)

	resp, err := client.Get(u)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	// A API retorna 200 com "[]" quando não encontra.
	// Não precisamos inspecionar o conteúdo; só checamos se há ao menos 1 item.
	var payload []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return false
	}
	fmt.Println("Payload length:", len(payload))
	return len(payload) > 0
}