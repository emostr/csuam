#!/usr/bin/env bash
#
# Автоматическое развёртывание на Debian
# Ставит Docker, nginx и certbot, собирает и запускает стек, настраивает HTTPS
# Все параметры задаются в deploy.conf (см. deploy.conf.example)
#
# Запуск:  sudo ./deploy.sh [--config FILE] [--check] [--skip-packages]

set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="$SCRIPT_DIR/deploy.conf"
CHECK_ONLY=no
SKIP_PACKAGES=no

NGINX_SITE_FILE=/etc/nginx/conf.d/csuam.conf
ERROR_PAGES_DIR=/var/www/csuam-errors
ACME_WEBROOT=/var/www/letsencrypt
SELF_SIGNED_DIR=/etc/ssl/csuam

if [[ -t 1 ]]; then
	C_RESET=$'\033[0m'; C_INFO=$'\033[36m'; C_OK=$'\033[32m'; C_WARN=$'\033[33m'; C_ERR=$'\033[31m'
else
	C_RESET=''; C_INFO=''; C_OK=''; C_WARN=''; C_ERR=''
fi

step() { printf '\n%s==>%s %s\n' "$C_INFO" "$C_RESET" "$*"; }
info() { printf '    %s\n' "$*"; }
ok() { printf '    %sготово%s %s\n' "$C_OK" "$C_RESET" "$*"; }
warn() { printf '%sвнимание:%s %s\n' "$C_WARN" "$C_RESET" "$*" >&2; }
die() { printf '%sошибка:%s %s\n' "$C_ERR" "$C_RESET" "$*" >&2; exit 1; }

trap 'die "сбой на строке $LINENO (команда: $BASH_COMMAND)"' ERR

usage() {
	cat <<'EOF'
Развёртывание ЦСУАМ «Архивли».

  sudo ./deploy.sh                  полное развёртывание (повторный запуск = обновление)
  sudo ./deploy.sh --check          проверить конфигурацию и окружение, ничего не менять
  sudo ./deploy.sh --skip-packages  не трогать apt-пакеты (Docker, nginx и certbot уже стоят)
  ./deploy.sh --help                эта справка

  --config FILE                     путь к конфигурации (по умолчанию ./deploy.conf)
EOF
}

parse_args() {
	while [[ $# -gt 0 ]]; do
		case $1 in
			--config) [[ $# -ge 2 ]] || die "--config требует путь к файлу"; CONFIG_FILE=$2; shift 2 ;;
			--check) CHECK_ONLY=yes; shift ;;
			--skip-packages) SKIP_PACKAGES=yes; shift ;;
			-h|--help) usage; exit 0 ;;
			*) die "неизвестный аргумент: $1 (--help для справки)" ;;
		esac
	done
}

load_config() {
	[[ -f $CONFIG_FILE ]] || die "нет файла конфигурации $CONFIG_FILE — скопируйте deploy.conf.example в deploy.conf и заполните его"
	# shellcheck disable=SC1090
	source "$CONFIG_FILE"

	: "${DOMAIN:=}"
	: "${ACME_EMAIL:=}"
	: "${TLS_MODE:=auto}"
	: "${ACME_STAGING:=${CADDY_STAGING:-no}}"
	: "${APP_DIR:=$SCRIPT_DIR}"
	: "${HTTP_PORT:=8000}"
	: "${HTTP_BIND:=127.0.0.1}"
	: "${INSTALL_DOCKER:=yes}"
	: "${INSTALL_NGINX:=${INSTALL_CADDY:-yes}}"
	: "${DISABLE_CADDY:=yes}"
	: "${FIREWALL:=no}"
	: "${SSH_PORT:=22}"
	: "${POSTGRES_USER:=csuam}"
	: "${POSTGRES_DB:=csuam}"
	: "${S3_BUCKET:=csuam}"
	: "${BACKUP_BEFORE_DEPLOY:=yes}"
	: "${BACKUP_KEEP:=10}"

	ENV_FILE="$APP_DIR/.env"
	BACKUP_DIR="$APP_DIR/backups"
}

validate_config() {
	[[ -n $DOMAIN ]] || die "в конфигурации не задан DOMAIN"
	[[ $DOMAIN =~ ^[A-Za-z0-9.-]+$ ]] || die "DOMAIN содержит недопустимые символы: $DOMAIN"
	[[ $HTTP_PORT =~ ^[0-9]+$ ]] || die "HTTP_PORT должен быть числом"
	[[ $SSH_PORT =~ ^[0-9]+$ ]] || die "SSH_PORT должен быть числом"
	[[ -z ${GARAGE_RPC_SECRET:-} || $GARAGE_RPC_SECRET =~ ^[0-9a-fA-F]{64}$ ]] \
		|| die "GARAGE_RPC_SECRET должен состоять из 64 шестнадцатеричных символов (openssl rand -hex 32)"
	[[ -z ${S3_ACCESS_KEY:-} || $S3_ACCESS_KEY =~ ^[A-Za-z0-9._-]{8,}$ ]] \
		|| die "S3_ACCESS_KEY: не короче 8 символов, только латиница, цифры и . _ -"
	[[ -z ${S3_SECRET_KEY:-} || ${#S3_SECRET_KEY} -ge 16 ]] || die "S3_SECRET_KEY должен быть не короче 16 символов"
	[[ $S3_BUCKET =~ ^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$ ]] || die "S3_BUCKET не подходит как имя бакета: $S3_BUCKET"

	case $TLS_MODE in
		auto)
			[[ -n $ACME_EMAIL ]] || die "при TLS_MODE=auto нужно указать ACME_EMAIL"
			[[ $ACME_EMAIL == *@*.* ]] || die "ACME_EMAIL не похож на адрес почты: $ACME_EMAIL"
			;;
		internal|off) ;;
		*) die "TLS_MODE может быть auto, internal или off (сейчас: $TLS_MODE)" ;;
	esac

	[[ -f $SCRIPT_DIR/docker-compose.yml ]] || die "рядом со скриптом нет docker-compose.yml — запускайте его из каталога проекта"
	ok "конфигурация корректна: $DOMAIN, TLS_MODE=$TLS_MODE, каталог $APP_DIR"
}

require_root() {
	[[ $EUID -eq 0 ]] || die "нужны права root: запустите через sudo"
}

detect_os() {
	[[ -f /etc/os-release ]] || die "не удалось определить дистрибутив (нет /etc/os-release)"
	# shellcheck disable=SC1091
	source /etc/os-release
	OS_ID=${ID:-unknown}
	OS_CODENAME=${VERSION_CODENAME:-}
	case $OS_ID in
		debian|ubuntu) ;;
		*) die "скрипт рассчитан на Debian или Ubuntu, обнаружено: $OS_ID" ;;
	esac
	[[ -n $OS_CODENAME ]] || die "не удалось определить кодовое имя выпуска ${OS_ID}"
	command -v systemctl >/dev/null || die "нужен systemd (systemctl не найден)"
	ok "система: ${PRETTY_NAME:-$OS_ID $OS_CODENAME}"
}

check_dns() {
	local server_ip domain_ip
	server_ip=$(curl -fsS --max-time 5 https://api.ipify.org 2>/dev/null || true)
	domain_ip=$(getent ahostsv4 "$DOMAIN" 2>/dev/null | awk 'NR==1 {print $1}' || true)
	if [[ -z $domain_ip ]]; then
		warn "домен $DOMAIN сейчас не резолвится — Let's Encrypt не сможет выдать сертификат, пока не появится A-запись"
		return
	fi
	if [[ -n $server_ip && $server_ip != "$domain_ip" ]]; then
		warn "домен $DOMAIN указывает на $domain_ip, а внешний адрес сервера — $server_ip; проверьте DNS"
	else
		ok "DNS: $DOMAIN → $domain_ip"
	fi
}

port_owner() {
	local port=$1
	ss -lntp "sport = :$port" 2>/dev/null | awk -F'"' 'NR > 1 && /users:/ {print $2; exit}' || true
}

check_ports() {
	local port owner
	for port in 80 443; do
		owner=$(port_owner "$port")
		[[ -z $owner ]] && continue
		case $owner in
			nginx) info "порт $port уже занят nginx — он будет перенастроен" ;;
			caddy) info "порт $port занят caddy — потребуется его отключить (DISABLE_CADDY=$DISABLE_CADDY)" ;;
			*) warn "порт $port занят процессом $owner — nginx не сможет его открыть" ;;
		esac
	done
}

apt_install() {
	local missing=()
	local pkg
	for pkg in "$@"; do
		dpkg-query -W -f='${Status}' "$pkg" 2>/dev/null | grep -q "ok installed" || missing+=("$pkg")
	done
	[[ ${#missing[@]} -eq 0 ]] && return 0
	info "устанавливаю: ${missing[*]}"
	DEBIAN_FRONTEND=noninteractive apt-get install -y -qq "${missing[@]}" >/dev/null
}

install_base_packages() {
	step "Базовые пакеты"
	info "обновляю списки пакетов"
	apt-get update -qq
	apt_install ca-certificates curl gnupg openssl git rsync iproute2 tar gzip
	ok "базовые пакеты на месте"
}

install_docker() {
	step "Docker"
	if command -v docker >/dev/null && docker compose version >/dev/null 2>&1; then
		ok "Docker и compose-плагин уже установлены ($(docker --version | awk '{print $3}' | tr -d ,))"
		return
	fi
	if [[ $INSTALL_DOCKER != yes ]]; then
		die "Docker не установлен, а INSTALL_DOCKER=no"
	fi
	install -m 0755 -d /etc/apt/keyrings
	curl -fsSL "https://download.docker.com/linux/$OS_ID/gpg" -o /etc/apt/keyrings/docker.asc
	chmod a+r /etc/apt/keyrings/docker.asc
	printf 'deb [arch=%s signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/%s %s stable\n' \
		"$(dpkg --print-architecture)" "$OS_ID" "$OS_CODENAME" >/etc/apt/sources.list.d/docker.list
	apt-get update -qq
	apt_install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
	systemctl enable --now docker >/dev/null
	ok "Docker установлен"
}

nginx_version() {
	nginx -v 2>&1 | sed -n 's#^nginx version: nginx/\([0-9.]*\).*#\1#p'
}

install_nginx() {
	step "nginx и certbot"
	local packages=(nginx) pkg missing=()
	[[ $TLS_MODE == auto ]] && packages+=(certbot)
	for pkg in "${packages[@]}"; do
		command -v "$pkg" >/dev/null || missing+=("$pkg")
	done
	if [[ ${#missing[@]} -eq 0 ]]; then
		ok "уже установлены: ${packages[*]} (nginx $(nginx_version))"
		return
	fi
	if [[ $INSTALL_NGINX != yes ]]; then
		die "не установлены ${missing[*]}, а INSTALL_NGINX=no"
	fi
	apt_install "${packages[@]}"
	systemctl enable nginx >/dev/null
	ok "установлены: ${packages[*]} (nginx $(nginx_version))"
}

backup_path() {
	local src=$1 name
	name=$(basename "$src")
	name=${name#.}
	mkdir -p "$BACKUP_DIR"
	cp -a "$src" "$BACKUP_DIR/${name}.$(date +%Y%m%d-%H%M%S).bak"
}

disable_caddy() {
	systemctl is-enabled caddy >/dev/null 2>&1 || systemctl is-active --quiet caddy || return 0
	step "Миграция с Caddy"
	if [[ $DISABLE_CADDY != yes ]]; then
		warn "Caddy оставлен включённым (DISABLE_CADDY=no) — nginx не сможет занять порты 80/443"
		return
	fi
	if [[ -d /etc/caddy ]]; then
		mkdir -p "$BACKUP_DIR"
		local archive
		archive="$BACKUP_DIR/caddy-config-$(date +%Y%m%d-%H%M%S).tar.gz"
		tar -czf "$archive" -C /etc caddy
		info "конфигурация Caddy сохранена: $archive"
	fi
	systemctl disable --now caddy >/dev/null 2>&1 || true
	ok "Caddy остановлен и убран из автозапуска: сертификаты теперь ведёт certbot"
}

sync_sources() {
	[[ $APP_DIR == "$SCRIPT_DIR" ]] && return 0
	step "Копирование исходников в $APP_DIR"
	mkdir -p "$APP_DIR"
	rsync -a --delete \
		--exclude '.git/' --exclude 'node_modules/' --exclude 'dist/' \
		--exclude '.env' --exclude 'deploy.conf' --exclude 'backups/' \
		"$SCRIPT_DIR"/ "$APP_DIR"/
	cp -f "$CONFIG_FILE" "$APP_DIR/deploy.conf"
	ok "исходники синхронизированы"
}

gen_secret() {
	local raw
	raw=$(openssl rand -base64 64 | tr -dc 'A-Za-z0-9')
	printf '%s' "${raw:0:${1:-32}}"
}

gen_hex() {
	openssl rand -hex 32
}

gen_access_key() {
	printf 'GK%s' "$(openssl rand -hex 12)"
}

env_value() {
	local key=$1 line
	[[ -f $ENV_FILE ]] || return 1
	line=$(grep -E "^${key}=" "$ENV_FILE" 2>/dev/null | tail -n1) || return 1
	[[ -n $line ]] || return 1
	printf '%s' "${line#*=}"
}

resolve_secret() {
	local name=$1 generator=${2:-gen_secret}
	local configured="${!name:-}"
	local existing
	existing=$(env_value "$name" || true)
	if [[ -n $configured ]]; then
		if [[ -n $existing && $existing != "$configured" ]]; then
			warn "$name в deploy.conf отличается от значения в .env — будет записано значение из deploy.conf"
		fi
		printf '%s' "$configured"
	elif [[ -n $existing ]]; then
		printf '%s' "$existing"
	else
		"$generator"
	fi
}

write_env() {
	step "Файл окружения"
	local scheme=https
	[[ $TLS_MODE == off ]] && scheme=http

	local db_pass rpc_secret s3_access s3_secret jwt seed_head seed_lib seed_teacher
	db_pass=$(resolve_secret POSTGRES_PASSWORD)
	rpc_secret=$(resolve_secret GARAGE_RPC_SECRET gen_hex)
	s3_access=$(resolve_secret S3_ACCESS_KEY gen_access_key)
	s3_secret=$(resolve_secret S3_SECRET_KEY gen_hex)
	jwt=$(resolve_secret JWT_SECRET)
	seed_head=$(resolve_secret SEED_HEAD_TEACHER_PASSWORD)
	seed_lib=$(resolve_secret SEED_LIBRARIAN_PASSWORD)
	seed_teacher=$(resolve_secret SEED_TEACHER_PASSWORD)

	[[ -f $ENV_FILE ]] && backup_path "$ENV_FILE"

	cat >"$ENV_FILE" <<EOF
# Создано deploy.sh $(date '+%Y-%m-%d %H:%M:%S'). Правки перезапишутся при следующем запуске:
# меняйте значения в deploy.conf.

POSTGRES_USER=$POSTGRES_USER
POSTGRES_PASSWORD=$db_pass
POSTGRES_DB=$POSTGRES_DB

GARAGE_RPC_SECRET=$rpc_secret
S3_ACCESS_KEY=$s3_access
S3_SECRET_KEY=$s3_secret
S3_BUCKET=$S3_BUCKET

JWT_SECRET=$jwt
PUBLIC_URL=$scheme://$DOMAIN
HTTP_PORT=$HTTP_PORT
HTTP_BIND=$HTTP_BIND

SEED_HEAD_TEACHER_PASSWORD=$seed_head
SEED_LIBRARIAN_PASSWORD=$seed_lib
SEED_TEACHER_PASSWORD=$seed_teacher
EOF
	chmod 600 "$ENV_FILE"
	ok "записан $ENV_FILE"
}

backup_database() {
	[[ $BACKUP_BEFORE_DEPLOY == yes ]] || return 0
	local container
	container=$(cd "$APP_DIR" && docker compose ps -q postgres 2>/dev/null || true)
	[[ -n $container ]] || return 0
	docker inspect -f '{{.State.Running}}' "$container" 2>/dev/null | grep -q true || return 0

	step "Резервная копия базы"
	mkdir -p "$BACKUP_DIR"
	local dump
	dump="$BACKUP_DIR/db-$(date +%Y%m%d-%H%M%S).sql.gz"
	if docker exec "$container" \
		pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" --clean --if-exists 2>/dev/null | gzip >"$dump"; then
		ok "дамп сохранён: $dump ($(du -h "$dump" | cut -f1))"
	else
		rm -f "$dump"
		warn "не удалось снять дамп базы — развёртывание продолжается"
	fi
	# shellcheck disable=SC2012
	ls -1t "$BACKUP_DIR"/db-*.sql.gz 2>/dev/null | tail -n +$((BACKUP_KEEP + 1)) | xargs -r rm -f || true
}

deploy_stack() {
	step "Сборка и запуск контейнеров"
	cd "$APP_DIR"
	docker compose build --pull
	docker compose up -d --remove-orphans
	ok "стек запущен"
}

wait_for_app() {
	step "Проверка приложения"
	local host=$HTTP_BIND
	[[ $host == 0.0.0.0 || $host == '::' ]] && host=127.0.0.1
	local code
	for _ in $(seq 1 60); do
		code=$(curl -fsS -o /dev/null -w '%{http_code}' --max-time 5 "http://$host:$HTTP_PORT/" 2>/dev/null || true)
		if [[ $code == 200 ]]; then
			ok "фронтенд отвечает на http://$host:$HTTP_PORT/"
			return 0
		fi
		sleep 2
	done
	cd "$APP_DIR" && docker compose ps
	die "приложение не поднялось за 2 минуты — смотрите: cd $APP_DIR && docker compose logs"
}

install_error_pages() {
	local src="$APP_DIR/frontend/public/errors"
	local pages=("$src"/*.html)
	if [[ ! -e ${pages[0]} ]]; then
		warn "не найдены страницы ошибок в $src — nginx будет отдавать стандартные ответы"
		return 0
	fi
	mkdir -p "$ERROR_PAGES_DIR"
	install -m 0644 "${pages[@]}" "$ERROR_PAGES_DIR"/
	info "страницы ошибок: $ERROR_PAGES_DIR (${#pages[@]} шт.)"
}

nginx_listen() {
	local port=$1 ssl=${2:-}
	printf '\tlisten %s%s;\n' "$port" "${ssl:+ ssl}"
	[[ -f /proc/net/if_inet6 ]] && printf '\tlisten [::]:%s%s;\n' "$port" "${ssl:+ ssl}"
	return 0
}

nginx_http2() {
	local current
	current=$(nginx_version)
	if [[ -n $current && $(printf '1.25.1\n%s\n' "$current" | sort -V | head -n1) == 1.25.1 ]]; then
		printf '\thttp2 on;\n'
	fi
}

nginx_app_locations() {
	cat <<EOF
	client_max_body_size 512m;

	access_log /var/log/nginx/csuam.access.log;
	error_log /var/log/nginx/csuam.error.log;

	error_page 502 /csuam-errors/502.html;
	error_page 500 503 504 /csuam-errors/500.html;

	location ^~ /csuam-errors/ {
		alias $ERROR_PAGES_DIR/;
		internal;
	}

	location / {
		proxy_pass http://127.0.0.1:$HTTP_PORT;
		proxy_http_version 1.1;
		proxy_set_header Host \$host;
		proxy_set_header X-Real-IP \$remote_addr;
		proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
		proxy_set_header X-Forwarded-Proto \$scheme;
		proxy_request_buffering off;
		proxy_buffering off;
		proxy_read_timeout 300s;
		proxy_send_timeout 300s;
	}
EOF
}

nginx_acme_location() {
	cat <<EOF
	location /.well-known/acme-challenge/ {
		root $ACME_WEBROOT;
	}
EOF
}

render_nginx_http() {
	printf '# Управляется csuam deploy.sh — правки перезапишутся при следующем запуске.\n'
	printf 'server {\n'
	nginx_listen 80
	printf '\tserver_name %s;\n\n' "$DOMAIN"
	if [[ $TLS_MODE == auto ]]; then
		nginx_acme_location
		printf '\n'
	fi
	nginx_app_locations
	printf '}\n'
}

render_nginx_https() {
	local cert=$1 key=$2
	printf '# Управляется csuam deploy.sh — правки перезапишутся при следующем запуске.\n'
	printf 'server {\n'
	nginx_listen 80
	printf '\tserver_name %s;\n\n' "$DOMAIN"
	nginx_acme_location
	printf '\n\tlocation / {\n\t\treturn 301 https://$host$request_uri;\n\t}\n'
	printf '}\n\n'
	printf 'server {\n'
	nginx_listen 443 ssl
	nginx_http2
	printf '\tserver_name %s;\n\n' "$DOMAIN"
	printf '\tssl_certificate %s;\n' "$cert"
	printf '\tssl_certificate_key %s;\n\n' "$key"
	nginx_app_locations
	printf '}\n'
}

reload_nginx() {
	nginx -t >/dev/null 2>&1 || die "nginx не принял конфигурацию: nginx -t"
	systemctl enable nginx >/dev/null 2>&1 || true
	if systemctl is-active --quiet nginx; then
		systemctl reload nginx
	else
		systemctl start nginx
	fi
}

ensure_self_signed() {
	local crt="$SELF_SIGNED_DIR/$DOMAIN.crt" key="$SELF_SIGNED_DIR/$DOMAIN.key"
	if [[ -f $crt && -f $key ]] && openssl x509 -checkend 2592000 -noout -in "$crt" >/dev/null 2>&1; then
		return 0
	fi
	install -d -m 0700 "$SELF_SIGNED_DIR"
	openssl req -x509 -newkey rsa:2048 -nodes -days 825 -subj "/CN=$DOMAIN" \
		-addext "subjectAltName=DNS:$DOMAIN" -keyout "$key" -out "$crt" >/dev/null 2>&1 \
		|| die "не удалось выпустить самоподписанный сертификат"
	chmod 600 "$key"
	info "выпущен самоподписанный сертификат: $crt"
}

obtain_certificate() {
	local live="/etc/letsencrypt/live/$DOMAIN"
	local args=(certonly --webroot -w "$ACME_WEBROOT" -d "$DOMAIN" --cert-name "$DOMAIN"
		--email "$ACME_EMAIL" --agree-tos --non-interactive --keep-until-expiring
		--deploy-hook 'systemctl reload nginx')
	if [[ $ACME_STAGING == yes ]]; then
		args+=(--staging)
	elif [[ -f $live/cert.pem ]] && openssl x509 -in "$live/cert.pem" -noout -issuer | grep -qi staging; then
		args+=(--force-renewal)
	fi
	if certbot "${args[@]}" >/dev/null 2>&1; then
		ok "сертификат Let's Encrypt для $DOMAIN на месте"
	else
		warn "certbot не смог получить сертификат для $DOMAIN — сайт пока доступен по HTTP; проверьте DNS и запустите deploy.sh ещё раз (журнал: /var/log/letsencrypt/letsencrypt.log)"
	fi
	systemctl enable --now certbot.timer >/dev/null 2>&1 || true
}

warn_nginx_conflicts() {
	local file
	while IFS= read -r file; do
		[[ $file == "$NGINX_SITE_FILE" ]] && continue
		warn "в $file тоже объявлен $DOMAIN — уберите старый сайт, иначе nginx выберет один из них"
	done < <(grep -rlsE "server_name[^;]*[[:space:]]$DOMAIN([[:space:]]|;)" /etc/nginx/conf.d /etc/nginx/sites-enabled 2>/dev/null || true)
}

write_nginx_site() {
	step "Настройка nginx"
	install_error_pages
	mkdir -p "$(dirname "$NGINX_SITE_FILE")"
	[[ -f $NGINX_SITE_FILE ]] && backup_path "$NGINX_SITE_FILE"

	case $TLS_MODE in
		off)
			render_nginx_http >"$NGINX_SITE_FILE"
			;;
		internal)
			ensure_self_signed
			render_nginx_https "$SELF_SIGNED_DIR/$DOMAIN.crt" "$SELF_SIGNED_DIR/$DOMAIN.key" >"$NGINX_SITE_FILE"
			;;
		auto)
			local live="/etc/letsencrypt/live/$DOMAIN"
			mkdir -p "$ACME_WEBROOT"
			if [[ ! -f $live/fullchain.pem ]]; then
				render_nginx_http >"$NGINX_SITE_FILE"
				reload_nginx
			fi
			obtain_certificate
			if [[ -f $live/fullchain.pem ]]; then
				render_nginx_https "$live/fullchain.pem" "$live/privkey.pem" >"$NGINX_SITE_FILE"
			fi
			;;
	esac

	warn_nginx_conflicts
	reload_nginx
	ok "записан $NGINX_SITE_FILE"
}

configure_firewall() {
	[[ $FIREWALL == ufw ]] || return 0
	step "Фаервол"
	apt_install ufw
	ufw allow "$SSH_PORT/tcp" >/dev/null
	ufw allow 80/tcp >/dev/null
	ufw allow 443/tcp >/dev/null
	ufw --force enable >/dev/null
	ok "ufw включён: открыты $SSH_PORT, 80, 443"
}

wait_for_https() {
	[[ $TLS_MODE == off ]] && return 0
	step "Проверка HTTPS"
	local curl_opts=(-fsS -o /dev/null --max-time 10)
	[[ $TLS_MODE == internal || $ACME_STAGING == yes ]] && curl_opts+=(-k)
	local code
	for _ in $(seq 1 24); do
		code=$(curl "${curl_opts[@]}" -w '%{http_code}' "https://$DOMAIN/" 2>/dev/null || true)
		if [[ $code == 200 ]]; then
			ok "https://$DOMAIN/ отвечает"
			return 0
		fi
		sleep 5
	done
	warn "https://$DOMAIN/ пока не отвечает. Проверьте DNS и журналы: /var/log/nginx/csuam.error.log, /var/log/letsencrypt/letsencrypt.log"
}

summary() {
	local seed_head
	seed_head=$(env_value SEED_HEAD_TEACHER_PASSWORD || true)
	local scheme=https
	[[ $TLS_MODE == off ]] && scheme=http

	printf '\n%s==> Развёртывание завершено%s\n\n' "$C_OK" "$C_RESET"
	printf '  Адрес:            %s://%s\n' "$scheme" "$DOMAIN"
	printf '  Каталог:          %s\n' "$APP_DIR"
	printf '  Секреты:          %s (права 600)\n' "$ENV_FILE"
	printf '  Конфиг nginx:     %s\n' "$NGINX_SITE_FILE"
	printf '  Резервные копии:  %s\n\n' "$BACKUP_DIR"
	printf '  Вход завуча:      zavuch / %s\n' "${seed_head:-см. .env}"
	printf '  Пароли учителя и библиотекаря — в %s\n\n' "$ENV_FILE"
	printf '  Логи приложения:  cd %s && docker compose logs -f\n' "$APP_DIR"
	printf '  Логи nginx:       tail -f /var/log/nginx/csuam.access.log\n'
	printf '  Обновление:       sudo %s/deploy.sh\n\n' "$APP_DIR"
}

main() {
	parse_args "$@"
	load_config
	validate_config

	if [[ $CHECK_ONLY == yes ]]; then
		step "Проверка окружения"
		[[ $EUID -eq 0 ]] || warn "скрипт запущен не от root: для развёртывания понадобится sudo"
		detect_os
		check_dns
		[[ $EUID -eq 0 ]] && check_ports
		if command -v docker >/dev/null; then ok "docker: $(docker --version)"; else info "docker будет установлен"; fi
		if command -v nginx >/dev/null; then ok "nginx: $(nginx_version)"; else info "nginx будет установлен"; fi
		if [[ $TLS_MODE == auto ]]; then
			if command -v certbot >/dev/null; then ok "$(certbot --version 2>&1)"; else info "certbot будет установлен"; fi
		fi
		printf '\n%sПроверка пройдена, изменений не вносилось.%s\n\n' "$C_OK" "$C_RESET"
		exit 0
	fi

	require_root
	detect_os
	check_dns
	check_ports

	if [[ $SKIP_PACKAGES == no ]]; then
		install_base_packages
		install_docker
		disable_caddy
		install_nginx
	else
		info "установка пакетов пропущена (--skip-packages)"
		command -v docker >/dev/null || die "Docker не установлен"
		command -v nginx >/dev/null || die "nginx не установлен"
		[[ $TLS_MODE != auto ]] || command -v certbot >/dev/null || die "certbot не установлен"
		disable_caddy
	fi

	sync_sources
	mkdir -p "$BACKUP_DIR"
	backup_database
	write_env
	deploy_stack
	wait_for_app
	write_nginx_site
	configure_firewall
	wait_for_https
	summary
}

main "$@"
