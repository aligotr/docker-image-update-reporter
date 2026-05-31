#!/usr/bin/env sh
set -eu

# Переменные
DEVELOPER="aligotr"
PROJECT_NAME="docker-check"
VERSION="1.0.0"

# Функция сборки
build() {
  docker build -t ${DEVELOPER}/${PROJECT_NAME}:latest -t ${DEVELOPER}/${PROJECT_NAME}:${VERSION} ./
}

# Вспомогательные функции
default() {
  build
}

logs() {
  get_container_id
  if [ ! -z ${CONTAINER_ID} ]; then
    docker logs ${CONTAINER_ID} --follow --tail 50
  else
    echo "Контейнер не найден"
  fi
}

rm() {
  get_container_id
  if [ ! -z ${CONTAINER_ID} ]; then
    echo -n "Контейнер остановлен: " && docker stop ${CONTAINER_ID}
    echo -n "Контейнер удалён: " && docker rm ${CONTAINER_ID} --force
  else
    echo "Контейнер не найден"
  fi
}

clean() {
  echo "ВНИМАНИЕ! Команда удалит не используемые docker-ресурсы:\n${clean_details}"
  prompt
  docker image prune --all --force
  docker system prune --volumes --force
}

remove() {
  echo "ВНИМАНИЕ! Команда удалит контейнер, его образ, тома, сеть и слои сборки"
  prompt
  IMAGE_ID=$(docker images --quiet --filter="reference=${DEVELOPER}/${PROJECT_NAME}:latest")
  if [ ! -z $IMAGE_ID ]; then
    rm
    echo -n "Образ удалён: " && docker rmi ${DEVELOPER}/${PROJECT_NAME}:latest --force
    docker system prune --volumes --force
  else
    echo "Образ не найден"
  fi
}

# Утилиты
prompt() {
  read -p "Продолжить выполнение? [y/n] " choice
  case $choice in
  [Yy]*)
    break
    ;;
  *)
    exit
    ;;
  esac
}

get_container_id() {
  CONTAINER_ID=$(docker ps --all --quiet --filter "label=project-name=$PROJECT_NAME")
}

help() {
  cat <<EOF
────────────────────────────────────
                             O o
                                o
 ______   ______   ______  [O]__ST
|""""""|_|""""""|_|""""""|_|======}
'-0--0-'"'-0--0-'"'-0--0-'"'000--o\.
────────────────────────────────────
Проект: ${PROJECT_NAME}
Тип: docker
────────────────────────────────────
Сборка образа: build
Вывести логи контейнера: logs
Удалить контейнер: rm
Удалить контейнер, его образ, тома, сеть и слои сборки: remove
Удалить не используемые docker-ресурсы: clean
${clean_details}
────────────────────────────────────
EOF
}

clean_details=$(
  cat <<EOF
- Остановленные контейнеры;
- Неиспользуемые образы;
- Неиспользуемые тома;
- Неиспользуемые сети;
- Кэш сборки.
EOF
)

# Сопоставление: Аргумент-Функция
if [ $# -eq 0 ]; then
  default
else
  for arg in "$@"; do
    case $arg in
    "build")
      build
      ;;
    "logs")
      logs
      ;;
    "rm")
      rm
      ;;
    "clean")
      clean
      ;;
    "remove")
      clean
      ;;
    "help")
      help
      ;;
    *)
      echo "────────────────────────────────────"
      echo "Неизвестный аргумент: $arg"
      help
      ;;
    esac
  done
fi
