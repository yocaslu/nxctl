# ☁️ Nextcloud - Docker Compose & Ferramentas de Backup

Este diretório contém as configurações de infraestrutura para implantar o **Nextcloud** via Docker Compose, além de um utilitário customizado para a realização de backups desenvolvido em **Go (Golang)**.

## 📁 Estrutura de Arquivos

Abaixo está a descrição dos arquivos contidos neste diretório:

- **`compose.yml`**: Arquivo de definição do Docker Compose responsável por orquestrar os contêineres do Nextcloud e seus serviços dependentes.
- **`.env`**: Arquivo para definição de variáveis de ambiente confidenciais (como senhas de banco de dados, credenciais e configurações de rede).
- **`.gitignore`**: Regras para ignorar arquivos locais e sensíveis no controle de versão Git.
- **`scripts/`**: Diretório que centraliza os scripts de automação e manutenção do serviço.
    - **`backup.go`**: Script automatizado de backup escrito em Go.
    - **`go.mod`**: Arquivo que define o módulo Go e gerencia as dependências do script de backup.

## 🚀 Como Iniciar o Serviço

Antes de iniciar, certifique-se de que o arquivo `.env` está configurado corretamente com as variáveis necessárias.

Para iniciar a stack do Nextcloud em segundo plano, execute:

```bash
docker compose up -d
```

## 📦 Como Executar o Backup

O utilitário de backup foi construído utilizando a linguagem Go. Para executá-lo, você precisará ter o [Go instalado](https://go.dev/) em seu ambiente ou executá-lo de dentro de um contêiner que possua o Go.

Para rodar o script de backup manualmente, navegue até a pasta de scripts e execute:

```bash
cd scripts
sudo go run backup.go
```
