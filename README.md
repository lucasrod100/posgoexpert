# PosGoExpert

## Desafio: Cotação do Dólar em Go

### 🖥️ Servidor (server/server.go)
    - 📤 Retorna a cotação do dólar no formato **JSON**.
    - 🗃️ Armazena a cotação no banco de dados **SQLite**.
### 🌐 Cliente (client/client.go)
    - 🔗 Faz a requisição HTTP para o servidor: **http://localhost:8080/cotacao**.
    - 📄 Recebe a cotação e salva no arquivo: **cotacao.txt**.