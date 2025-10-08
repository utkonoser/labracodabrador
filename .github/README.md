# GitHub Actions - Автоматический деплой

## 🚀 Что делает

При каждом пуше в ветку `main`:
1. ✅ Запускает тесты Go
2. 🔨 Собирает проект
3. 🚀 Деплоит на сервер через SSH

## ⚙️ Настройка

### 1. Добавьте Secrets в GitHub

Перейдите в Settings → Secrets and variables → Actions и добавьте:

| Secret | Описание | Пример |
|--------|----------|--------|
| `SERVER_HOST` | IP или домен сервера | `123.45.67.89` или `myserver.com` |
| `SERVER_USER` | Пользователь для SSH | `root` или `deploy` |
| `SSH_PRIVATE_KEY` | Приватный SSH ключ | Содержимое файла `~/.ssh/id_rsa` |
| `SERVER_PORT` | Порт SSH (опционально) | `22` (по умолчанию) |
| `PROJECT_PATH` | Путь к проекту на сервере (опционально) | `/opt/labracodabrador` (по умолчанию) |

### 2. Создайте SSH ключ (если нет)

На вашем компьютере:
```bash
ssh-keygen -t rsa -b 4096 -C "deploy@labracodabrador"
```

### 3. Добавьте публичный ключ на сервер

```bash
ssh-copy-id -i ~/.ssh/id_rsa.pub user@your-server
```

Или вручную:
```bash
cat ~/.ssh/id_rsa.pub | ssh user@your-server "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys"
```

### 4. Добавьте приватный ключ в GitHub Secrets

```bash
# Скопируйте содержимое приватного ключа
cat ~/.ssh/id_rsa
```

Скопируйте весь вывод (включая `-----BEGIN` и `-----END`) и вставьте в GitHub Secret `SSH_PRIVATE_KEY`.

### 5. Подготовьте сервер

На сервере должны быть установлены:
- Git
- Docker
- Docker Compose

Клонируйте репозиторий на сервер:
```bash
cd /opt
git clone https://github.com/nn-selin/labracodabrador.git
cd labracodabrador
```

### 6. Готово! 🎉

Теперь при каждом пуше в `main` проект автоматически задеплоится на сервер!

## 📝 Как использовать

```bash
# Внесите изменения
git add .
git commit -m "Update feature"
git push origin main

# Следите за процессом на GitHub:
# https://github.com/nn-selin/labracodabrador/actions
```

## 🔍 Проверка деплоя

После деплоя проверьте, что всё работает:
```bash
ssh user@your-server
cd /opt/labracodabrador
docker ps  # Должны работать контейнеры
curl http://localhost:8081/health  # API должно отвечать
```

## ⚠️ Troubleshooting

### SSH подключение не работает
```bash
# Проверьте ключ локально
ssh -i ~/.ssh/id_rsa user@your-server

# Проверьте права на сервере
chmod 700 ~/.ssh
chmod 600 ~/.ssh/authorized_keys
```

### Docker команды требуют sudo
Добавьте пользователя в группу docker:
```bash
sudo usermod -aG docker $USER
# Перелогиньтесь
```

### Git pull требует аутентификацию
Настройте SSH для GitHub на сервере:
```bash
ssh-keygen -t rsa
cat ~/.ssh/id_rsa.pub
# Добавьте ключ в GitHub Settings → SSH Keys
```

Или используйте HTTPS с токеном:
```bash
git remote set-url origin https://TOKEN@github.com/nn-selin/labracodabrador.git
```
