<p align="center">
  <picture>
    <source height="125" media="(prefers-color-scheme: dark)" srcset="assets/admin.png">
    <img height="125" alt="Fiber" src="assets/light-admin.png">
  </picture>
</p>

<p align="center">
   <strong>MS-admin</strong> — это  <strong>микросервис</strong> для управления и администрирования MongoDB, предоставляющий  <strong>REST API</strong> с набором инструментов для безопасной работы с данными. Он предназначен для внутреннего использования администраторами и обладает гибкой системой доступа, основанной на ролях (администратор и ограниченный администратор).
</p>

# 💡Основные возможности

- **Управление данными**: Полноценный CRUD (создание, чтение, обновление, удаление записей в коллекциях).  
- **Работа с коллекциями**: Получение списка коллекций, массовые операции, группировка данных.  
- **Резервное копирование**: Автоматические бэкапы (`dump`), восстановление (`upload`) и очистка данных (`drop`).  
- **Логирование**: Просмотр системных логов для мониторинга действий.  
- **Безопасность**: CORS, Cookie, разделение прав доступа.  

# 🤖 Используемые технологии

- **Golang** — основной язык программирования
- **MS-database** — основая библиотека для микросервиса
- **Fiber** — фреймворк для написания REST API
- **MongoDB** — основная база данных
- **Redis** — для управления доступом
- **Docker** — развертывание проекта

# ⚠️Важно
Перед стартом необходимо перейти в [MS-database](https://github.com/Muraddddddddd9/ms-database) и поднять MongoDB, Redis, S3

# ⚡️ Быстрый старт
Перейти в env и поменять конфигурацию
```env
MONGO_NAME=diary
# MONGO_HOST=localhost # <- для локального запуска
MONGO_HOST=host.docker.internal # <- для запуска в Docker 
MONGO_PORT=27018 # <- ваш порт (27018 для Docker)
MONGO_USERNAME=college # <- username для MongoDB
MONGO_PASSWORD=BIM_LOCAL1 # <- пароль для MongoDB
MONGO_AUTH_SOURCE=admin

REDIS_DB=0
# REDIS_HOST=localhost # <- для локального запуска
REDIS_HOST=host.docker.internal # <- для запуска в Docker 
REDIS_PASSWORD=BIM_LOCAL1 # <- пароль для MongoDB
REDIS_PORT=6380 # <- ваш порт (6380 для Docker)

ORIGIN_URL=http://localhost:5173 # <- адрес сайта
PROJECT_PORT=:8080 # <- порт приложения

ADMIN_EMAIL=admin # <- начальный email админа 
ADMIN_PASSWORD=BIM_LOCAL1 # <- начальный пароль админа 
```

## CMD
Клонирование репозитория
```bash
git clone https://github.com/Muraddddddddd9/ms-admin.git
```
Установка всех пакетов
```bash
go get .
```
Запустить программу
```bash
go run .
```
## Docker
Клонирование репозитория
```bash
git clone https://github.com/Muraddddddddd9/ms-admin.git
```
Билд Docker container 
```bash
docker-compose build
```
Поднятие Docker container 
```bash
docker-compose up
```

# 🧬 API
- <strong>create_data</strong> — Post, создание данных
- <strong>get_data/:collection</strong> — Get, получение данных по коллекции
- <strong>get_logs</strong> — Get, получение логов
- <strong>update_data</strong> — Patch, обновление данных
- <strong>delete_data</strong> — Delete, удаление данных
- <strong>get_all_object/:group</strong> — Get, получение всех предметов по группе
- <strong>get_collections</strong> — Get,  получение всех коллекций(название)
- <strong>dump</strong> — Post, скачивание коллекций
- <strong>drop</strong> — Post, удаление коллекции
- <strong>up_group</strong> — Patch, обновление всех групп в конце года
- <strong>upload</strong> — Post, загрузка коллекций в базы данных

# 🧩 Остальные
- <strong>[MS-common](https://github.com/Muraddddddddd9/ms-common)</strong> - микросервис (необходимый)
- <strong>[MS-database](https://github.com/Muraddddddddd9/ms-database)</strong> - микросервис (необходимый)
- <strong>[MS-teacher](https://github.com/Muraddddddddd9/ms-teacher)</strong> - микросервис (необходимый)
- <strong>[MS-student](https://github.com/Muraddddddddd9/ms-student)</strong> - микросервис (необходимый)
- <strong>[MS-telegram](https://github.com/Muraddddddddd9/ms-telegram)</strong> - микросервис (необходимый)
- <strong>[MDiary](https://github.com/Muraddddddddd9/MDiary)</strong> - Вебсайт (необходимый)