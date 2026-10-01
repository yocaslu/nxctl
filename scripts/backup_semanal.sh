#!/bin/bash

# ==========================================
# Configurações do Backup
# ==========================================

# Diretório onde os dados dos usuários estão montados
DATA_DIR="/mnt/nextcloud/"

# Nome do container Docker onde o Nextcloud estárodando (execute `docker ps` para encontrar esta informa��o) 
NEXTCLOUD_CONTAINER_NAME="nextcloud"

# Usuário do servidor web (ex: www-data no Ubuntu/Debian, apache no CentOS)
WEB_USER="www-data"

# Destino e nome do arquivo de backup (com data atual)
BACKUP_DEST="/mnt/backup/nextcloud/backup_semanal/nextcloud_user_data_$(date +'%d-%m-%Y').tar.gz"

# ==========================================
# Execução do Backup
# ==========================================

echo "Iniciando o backup dos dados do Nextcloud..."

# 1. Ativar modo de manutenção
# O comando occ deve sempre ser executado como o usuário HTTP para não quebrar permissões
echo "Colocando o Nextcloud em modo de manutenção..."
docker exec -u $WEB_USER -it $NEXTCLOUD_CONTAINER_NAME php occ maintenance:mode --on

# 2. Criar o arquivo tar.gz
# A flag 'c' cria o arquivo, 'z' compacta em gzip, 'p' preserva as permissões e 'f' especifica o nome do arquivo
echo "Compactando e copiando os arquivos de $DATA_DIR..."
sudo -u $USER tar -czpf $BACKUP_DEST $DATA_DIR

# 3. Desativar modo de manutenção
# Libera o acesso novamente para os usuários
echo "Desativando o modo de manutenção..."
docker exec -it $NEXTCLOUD_CONTAINER_NAME php occ maintenance:mode --off
echo "Backup finalizado com sucesso! Arquivo salvo em: $BACKUP_DEST"
