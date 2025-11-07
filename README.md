# Alunos: 
- Brunno Perez Bocardo
- Stephany Tucunduva Santos Miguel
---
# API Cifra de Cesar

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
```
{
    "textoClaro": "HELLO WORLD",
    "deslocamento": 3
}
```
- Reponse esperada:
{
    "textoCifrado": "KHOOR ZRUOG"
}

2) [POST] /decifrar
- http://localhost:8080/decifrar
- Body JSON:
```
{
    "textoCifrado": "KHOOR ZRUOG",
    "deslocamento": 3
}
```
- Response esperada:
```
{
    "textoClaro": "HELLO WORLD"
}
```

3) [POST] /decifrarForcaBruta
- http://localhost:8080/decifrarForcaBruta
- Body JSON:
```
{
    "textoCifrado": "KHOOR ZRUOG"
}
```
- Response esperada:
```
{
    "textoClaro": "HELLO WORLD"
}
```
