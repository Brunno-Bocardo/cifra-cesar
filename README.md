# Aluno: Brunno Perez Bocardo
## API Cifra de Vernam

Linguagem utilizada: Go

### Instruções:
1) Instale o Go seguindo as instruções conforme o seu sistema operacional: https://go.dev/dl/
2) Certifique-se que o Go está instalado usando `go version`
3) Na raíz do projeto, execute `go run main.go`. Aceite as permições do Go e você deverá ver a seguinte mensagem: `Programa rodando em: http://localhost:8080`

### Endpoints disponíveis:
O endpoins podem ser testados via CURL ou Postman.

1) [POST] /cifrar
- http://localhost:8080/cifrar
- Body JSON:
{
    "textoClaro": "malfeito feito",
    "chave": "CHGENPKTGLPCMW"
}
- Reponse esperada:
{
    "textoCifrado": "0010111000101001001010110010001100101011001110010011111100111011011001110010101000110101001010100011100100111000"
}

2) [POST] /decifrar
- http://localhost:8080/decifrar
- Body JSON:
{
    "textoCifrado": "0010111000101001001010110010001100101011001110010011111100111011011001110010101000110101001010100011100100111000",
    "chave": "CHGENPKTGLPCMW"
}
- Response esperada:
{
    "textoClaro": "MALFEITO FEITO"
}
