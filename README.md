# PosGoExpert

## Desafio: Cotação do Dólar em Go

### 🖥️ Servidor (server/server.go)
- 📤 Retorna a cotação do dólar no formato **JSON**.
- 🗃️ Armazena a cotação no banco de dados **SQLite**.
### 🌐 Cliente (client/client.go)
- 🔗 Faz a requisição HTTP para o servidor: **http://localhost:8080/cotacao**.
- 📄 Recebe a cotação e salva no arquivo: **cotacao.txt**.

## Desafio: Multithreading
- Busca o CEP fazendo requisições simultânias para as APIs da **BrasilApi** e **ViaCep**
- O resultado é exibido da API que entregou a resposta mais rápida
- Caso não ocorra resposta no tempo limit de 1 segundo, será retornado o resultado de erro
- A exebição do resultado é no command line com os dados do endereço e especificando qual API que entregou o resultado
- A execução é realizada via command line: **go run main.go <CEP>**