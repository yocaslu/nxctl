# nxctl (Nextcloud Control CLI) 🚀

> Utilitário de linha de comando (CLI) em Go para automação, gerenciamento operacional, snapshots e rotinas de backup de instâncias Nextcloud conteinerizadas via Docker.

---

## 📋 Sumário
- [Sobre o Projeto](#-sobre-o-projeto)
- [Arquitetura e Recursos](#-arquitetura-e-recursos)
- [Estrutura do Repositório](#-estrutura-do-repositório)
- [Pré-requisitos](#-pré-requisitos)
- [Configuração (.env.example)](#-configuração-envexample)
- [Instalação e Compilação](#-instalação-e-compilação)
- [Comandos Principais](#-comandos-principais)
  - [`nxctl snapshot`](#nxctl-snapshot)
  - [`nxctl backup`](#nxctl-backup)
- [Automação e Scripts](#-automação-e-scripts)
- [Desenvolvimento e Contribuição](#-desenvolvimento-e-contribuição)
- [Licença](#-licença)

---

## 📖 Sobre o Projeto

O **nxctl** é uma ferramenta desenvolvida em Go projetada para simplificar a administração de ecossistemas Nextcloud em produção ou desenvolvimento local. Ele orquestra ações no container do Nextcloud (habilitando/desabilitando modo de manutenção com `occ`), manipula dumps consistentes no banco PostgreSQL e cria snapshots locais (arquivamento `.tar`) ou backups incrementais/deduplicados com **Restic**.

---

## 🛠 Arquitetura e Recursos

- **Orquestração Docker**: Comunicação direta com a API/CLI do Docker para gerenciar serviços Nextcloud e PostgreSQL.
- **Segurança e Consistência**: Ativação automática do modo de manutenção do Nextcloud durante operações de cópia para garantir integridade referencial dos dados.
- **Snapshots Locais Rápidos**: Empacotamento direto do volume de dados/configurações e dump SQL em arquivos comprimidos (`.tar.gz`).
- **Backups Incrementais com Restic**: Integração nativa para sincronização de repositórios Restic com criptografia, deduplicação e versionamento seguro.
- **Templates de Proxy Nginx**: Configurações modulares e templates para proxy reverso com suporte a SSL, WebSockets, upload de grandes arquivos, Collabora Office e pgAdmin.

---

## 📂 Estrutura do Repositório

```text
.
├── .env.example                # Modelo de variáveis de ambiente do projeto
├── compose.yml                 # Definição dos containers Docker (Nextcloud, DB, Proxy, etc.)
├── main.go                     # Ponto de entrada da aplicação Go
├── cmd/                        # Definições de comandos CLI (Cobra/Viper)
│   ├── root.go                 # Comando raiz e flags globais
│   ├── snapshot.go             # Implementação do comando 'snapshot'
│   └── backup.go               # Implementação do comando 'backup'
├── internal/                   # Pacotes internos da aplicação
│   ├── config/                 # Carregamento e parse das configurações (.env / YAML)
│   ├── archive/                # Módulos de compactação
│   │   ├── tar/                # Criação e extração de tarballs
│   │   └── restic/             # Wrapper e chamadas para Restic
│   ├── docker/                 # Integração Docker
│   │   ├── nextcloud/          # Comandos específicos do Nextcloud (occ, etc.)
│   │   └── postgres/           # Rotinas de dump e restore do PostgreSQL
│   ├── proc/                   # Execução e monitoramento de subprocessos do SO
│   └── utils/                  # Utilitários (datas, diretórios e logging nxlog)
├── proxy/                      # Templates e snippets de configuração do Nginx
│   ├── snippets/               # Headers SSL, Proxy, WebSockets e uploads
│   └── templates/              # Configurações dinâmicas para Nextcloud, Collabora, pgAdmin
└── scripts/                    # Scripts complementares para automação (ex: cron jobs)
    └── backup_semanal.sh       # Script shell para agendamento periódico
```

---

## ⚙️ Pré-requisitos

- **Go**: Versão 1.21 ou superior (para compilação a partir do código-fonte).
- **Docker & Docker Compose**: Para rodar a stack de containers.
- **Restic** (opcional/recomendado): Necessário caso utilize o comando `nxctl backup` com repositórios Restic.
- **Ferramentas padrão UNIX**: `tar`, `gzip`, `bash`.

---

## 🔐 Configuração (`.env.example`)

Antes de executar o **nxctl** ou subir o ambiente via Docker Compose, clone o arquivo de exemplo `.env.example` para `.env` e defina os parâmetros correspondentes ao seu ambiente.

```bash
cp .env.example .env
```

### Variáveis Críticas de Configuração

| Variável | Descrição | Exemplo |
| :--- | :--- | :--- |
| `POSTGRES_DB` | Nome do banco de dados do Nextcloud | `nextcloud` |
| `POSTGRES_USER` | Usuário do banco de dados | `nc_user` |
| `POSTGRES_PASSWORD` | Senha segura do banco PostgreSQL | `sua_senha_secreta_aqui` |
| `NEXTCLOUD_DATA_DIR` | Caminho do volume de dados do Nextcloud no host | `/var/lib/docker/volumes/nextcloud_data/_data` |
| `SNAPSHOT_OUTPUT_DIR` | Diretório de destino para arquivos de snapshot (`.tar`) | `/var/backups/nextcloud/snapshots` |
| `RESTIC_REPOSITORY` | Caminho local ou URI remota (S3, B2, SFTP) do Restic | `/var/backups/restic-repo` ou `s3:s3.amazonaws.com/meu-bucket` |
| `RESTIC_PASSWORD` | Senha mestra de criptografia do repositório Restic | `chave_criptografia_restic` |
| `NEXTCLOUD_CONTAINER` | Nome do container Nextcloud no Docker Compose | `nextcloud_app` |
| `POSTGRES_CONTAINER` | Nome do container PostgreSQL no Docker Compose | `nextcloud_db` |

> ⚠️ **Atenção:** Nunca versione o arquivo `.env` com senhas reais no Git. Certifique-se de que ele permaneça listado no `.gitignore`.

---

## 📦 Instalação e Compilação

Para compilar o binário localmente a partir da raiz do repositório:

```bash
# Baixar dependências
go mod download

# Compilar o binário nxctl
go build -o nxctl main.go

# (Opcional) Mover para o PATH do sistema
sudo mv nxctl /usr/local/bin/
```

---

## ⚡ Comandos Principais

O CLI disponibiliza comandos dedicados para as estratégias de preservação de dados:

### `nxctl backup`

O comando `backup` cria uma fotografia local, pontual e autocontida do estado atual da instância.

#### Como funciona:
1. Conecta-se ao container do Nextcloud e aciona o **modo de manutenção** via `occ maintenance:mode --on`.
2. Executa um dump seguro da base de dados PostgreSQL (`pg_dump`).
3. Empacota a base exportada junto aos arquivos de configuração e dados essenciais em um arquivo `.tar.gz` datado.
4. Desativa o modo de manutenção do Nextcloud (`occ maintenance:mode --off`).

#### Exemplo de uso:
```bash
# Executa o backup com configurações padrões do .env
./nxctl backup 

# Modo não interativo / debug para logs detalhados
./nxctl backup --debug
```

---

### `nxctl backup`

O comando `snapshot` integra-se ao **Restic** para gerar backups incrementais, deduplicados e criptografados, ideais para armazenamento offsite e retenção de longo prazo.

#### Como funciona:
1. Carrega as credenciais e localização do repositório a partir de `RESTIC_REPOSITORY` e `RESTIC_PASSWORD`.
2. Habilita o modo de manutenção do Nextcloud.
3. Gera o dump consistente do PostgreSQL para uma área intermediária.
4. Invoca o Restic para calcular deltas e enviar apenas os blocos modificados para o repositório.
5. Restaura a operação normal do Nextcloud.
6. (Opcional) Aplica políticas de retenção (`forget` / `prune`).

#### Exemplo de uso:
```bash
# Execução padrão do backup incremental
./nxctl snapshot 
```

---

## 🤖 Automação e Scripts

Na pasta `scripts/`, encontra-se o script utilitário `backup_semanal.sh`, pronto para ser integrado a um agendador como `cron` ou `systemd timers`:

```bash
# Tornar executável
chmod +x scripts/backup_semanal.sh

# Exemplo de entrada no crontab (todos os domingos às 02h00)
# 0 2 * * 0 /opt/nxctl/scripts/backup_semanal.sh >> /var/log/nxctl_backup.log 2>&1
```

---

## 🤝 Desenvolvimento e Contribuição

1. Faça um Fork do projeto
2. Crie uma branch para a funcionalidade: `git checkout -b feature/minha-feature`
3. Commit suas alterações: `git commit -m 'feat: adiciona suporte a retenção customizada'`
4. Push para a branch: `git push origin feature/minha-feature`
5. Abra um Pull Request

---

## 📄 Licença

Este projeto está distribuído sob a licença definida no arquivo [LICENSE](LICENSE).
