# Деплой startup_back

Go/Fiber API с Postgres и S3-совместимым хранилищем. Схема БД создаётся
автоматически (GORM AutoMigrate) при старте приложения.

## 1. Переменные окружения

Все настройки читаются в таком порядке приоритета:
**env → `config.yaml` → значения по умолчанию в коде.**

```bash
cp .env.example .env
# сгенерировать секреты (должны быть разными)
openssl rand -hex 32   # -> ACCESS_SECRET
openssl rand -hex 32   # -> REFRESH_SECRET
```

Обязательные (приложение не стартует без них):
`DB_USER`, `DB_PASSWORD`, `DB_NAME`, `ACCESS_SECRET`, `REFRESH_SECRET`,
`S3_ENDPOINT`, `S3_REGION`, `S3_BUCKET`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`.

Важные для прода:

| Переменная | Значение в проде | Зачем |
|---|---|---|
| `APP_ENV` | `production` | выключает swagger и debug-логи БД |
| `SERVER_HOST` | `0.0.0.0` | иначе сервис недоступен снаружи контейнера |
| `CORS_ORIGINS` | домен фронтенда | `*` не работает вместе с cookie-авторизацией |
| `DB_SSLMODE` | `require` для managed-Postgres | `disable` только для БД в той же сети |

## 2. Запуск через docker compose

```bash
docker compose up -d --build
docker compose logs -f api
```

Поднимаются два сервиса: `db` (postgres:16, данные в volume `pgdata`) и `api`.
`api` стартует только после того, как healthcheck базы стал зелёным.

Порт `api` публикуется на `127.0.0.1:3001` — наружу его отдаёт реверс-прокси
(см. п. 4). Порт Postgres наружу не публикуется вовсе.

## 3. Запуск только образа (внешняя БД)

```bash
docker build --build-arg VERSION=$(git rev-parse --short HEAD) -t startup-back:latest .

docker run -d --name startup-api \
  --env-file .env \
  -e APP_ENV=production \
  -e SERVER_HOST=0.0.0.0 \
  -p 127.0.0.1:3001:3001 \
  --restart unless-stopped \
  startup-back:latest
```

Образ: multi-stage, статический бинарник (`CGO_ENABLED=0`), runtime `alpine`,
процесс работает от непривилегированного пользователя `app` (uid 10001).

## 4. Реверс-прокси и HTTPS

HTTPS обязателен: куки авторизации выставляются с флагом `Secure` и
без TLS браузер их не сохранит.

```nginx
server {
    listen 443 ssl http2;
    server_name api.example.com;

    ssl_certificate     /etc/letsencrypt/live/api.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;

    client_max_body_size 32m;   # совпадает с BodyLimit приложения

    location / {
        proxy_pass http://127.0.0.1:3001;
        proxy_set_header Host              $host;
        proxy_set_header X-Real-IP         $remote_addr;
        proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

> **Домен фронтенда.** Куки выставляются с `SameSite=Strict`. Если фронтенд живёт
> на другом сайте (например, `app.example.com` против `api.other.com`), браузер
> не отправит куки на API и авторизация сломается. Рабочие варианты: держать API
> поддоменом того же сайта (`api.example.com` + `app.example.com`) или поменять
> куки на `SameSite=None` в `internal/auth/handler.go`.

## 5. Проверки состояния

| Endpoint | Смысл |
|---|---|
| `GET /health` | процесс жив (liveness) |
| `GET /ready` | есть соединение с БД (readiness) |

В образе прописан `HEALTHCHECK`, дёргающий `/health`. Для балансировщика и
Kubernetes используйте `/ready` — он не пустит трафик, пока база недоступна.

```bash
curl -f http://127.0.0.1:3001/health
curl -f http://127.0.0.1:3001/ready
```

## 6. CI/CD (GitHub Actions)

Workflow `.github/workflows/ci-cd.yml` есть и здесь, и во фронтенде
(`IdeaCrafter_Front`). Схема одинаковая:

| Событие | test / check | image | deploy |
|---|---|---|---|
| Pull request | ✅ | сборка без push | — |
| Push в `main` | ✅ | push в GHCR: `sha-<commit>` и `latest` | если `DEPLOY_ENABLED=true` |
| Ручной запуск | ✅ | push в GHCR | если `DEPLOY_ENABLED=true` и ветка `main` |

Образы: `ghcr.io/zkuncry/ideacrafter_back` и `ghcr.io/zkuncry/ideacrafter_front`.
Деплой выключен, пока не задана переменная `DEPLOY_ENABLED`, поэтому до
появления сервера пайплайн только проверяет код и публикует образы.

### Подготовка сервера (один раз)

```bash
# Docker + compose plugin
curl -fsSL https://get.docker.com | sh

# отдельный пользователь для деплоя
sudo useradd -m -s /bin/bash deploy
sudo usermod -aG docker deploy
sudo mkdir -p /opt/ideacrafter/back /opt/ideacrafter/front
sudo chown -R deploy:deploy /opt/ideacrafter

# .env бэкенда кладётся руками, в GitHub он не хранится
sudo -u deploy nano /opt/ideacrafter/back/.env   # по образцу .env.example
```

Ключ для GitHub (на своей машине):

```bash
ssh-keygen -t ed25519 -C github-deploy -f deploy_key -N ""
ssh-copy-id -i deploy_key.pub deploy@SERVER_IP
ssh-keyscan -H SERVER_IP                         # -> SSH_KNOWN_HOSTS
```

### Настройки в GitHub (в обоих репозиториях)

Settings → Environments → создать `production`, в нём:

| Тип | Имя | Значение |
|---|---|---|
| Secret | `SSH_HOST` | IP или домен сервера |
| Secret | `SSH_USER` | `deploy` |
| Secret | `SSH_PRIVATE_KEY` | содержимое `deploy_key` |
| Secret | `SSH_KNOWN_HOSTS` | вывод `ssh-keyscan` |
| Variable | `DEPLOY_PATH` | `/opt/ideacrafter/back` или `/opt/ideacrafter/front` |
| Variable | `SSH_PORT` | необязательно, по умолчанию `22` |

Settings → Secrets and variables → Actions → Variables: `DEPLOY_ENABLED` = `true`.

В environment `production` можно включить Required reviewers, чтобы каждый
деплой подтверждался вручную.

### Как проходит деплой

1. `docker-compose.yml` копируется на сервер в `DEPLOY_PATH`.
2. Сервер логинится в GHCR временным `GITHUB_TOKEN`, который живёт только пока
   идёт job, так что постоянный токен на сервере не нужен.
3. `docker compose pull` и `up -d --wait` на образ с тегом `sha-<commit>`.
   Если контейнер не стал healthy за 120 секунд, job падает.

**Порядок при первом деплое:** сначала бэкенд, потом фронтенд. Фронт
подключается к docker-сети `startup_back_default`, которую создаёт бэкенд.

**Откат:** Actions → старый успешный запуск → Re-run jobs. Каждый деплой
привязан к неизменяемому тегу `sha-…`.

**Ручные команды на сервере.** Чтобы `docker compose` на сервере использовал
образ из GHCR, а не пытался собрать его из исходников, добавьте в
`/opt/ideacrafter/back/.env` строку
`API_IMAGE=ghcr.io/zkuncry/ideacrafter_back:latest`
(для фронта в `/opt/ideacrafter/front/.env`:
`WEB_IMAGE=ghcr.io/zkuncry/ideacrafter_front:latest`) и один раз выполните
`docker login ghcr.io` с PAT, у которого есть право `read:packages`.

Приложение обрабатывает `SIGTERM`: перестаёт принимать запросы, дожидается
активных (до 20 секунд) и закрывает пул соединений с БД.

## 7. Бэкап базы

```bash
docker compose exec -T db pg_dump -U "$DB_USER" "$DB_NAME" | gzip > backup-$(date +%F).sql.gz
# восстановление
gunzip -c backup-2026-09-21.sql.gz | docker compose exec -T db psql -U "$DB_USER" -d "$DB_NAME"
```

## Чеклист перед первым деплоем

- [ ] `.env` заполнен, `ACCESS_SECRET` и `REFRESH_SECRET` разные и сгенерированы случайно
- [ ] `APP_ENV=production`
- [ ] `CORS_ORIGINS` указывает на реальный домен фронтенда
- [ ] `DB_SSLMODE=require`, если база вне docker-сети
- [ ] Бакет S3 создан, ключи имеют права на запись
- [ ] TLS-сертификат выпущен, прокси проксирует на `127.0.0.1:3001`
- [ ] `curl -f https://api.example.com/ready` отвечает 200
- [ ] Настроен бэкап volume `pgdata`
- [ ] В обоих репозиториях заполнен environment `production` и `DEPLOY_ENABLED=true`
