###################################################################
# Создание промежуточного образа
###################################################################
FROM golang:alpine AS builder

# Переменные
ARG PROJECT_NAME="docker-check"
WORKDIR /srv/app

COPY ./app/ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o build/${PROJECT_NAME} .

###################################################################
# Создание финального образа
###################################################################
FROM alpine:latest

# Метки
ARG PROJECT_NAME="docker-check"
ENV PROJECT_NAME=$PROJECT_NAME
LABEL project-name=$PROJECT_NAME
LABEL maintainer="Aligotr"

# Переменные
ARG \
  UID=1000 \
  USERNAME=appuser \
  GID_DOCKER=997
ENV \
  TZ="Europe/Moscow" \
  HOME="/srv/app"
WORKDIR $HOME

##########      Настройка системы      ##########
#
RUN \
  echo "--- Установка зависимостей и дополнительных программ ---" && \
  apk add --no-cache \
    tzdata \
    nano && \
  echo "--- Очистка системы ---" && \
  rm -rf \
    /tmp/*

##########   Настройка и запуск основного процесса   #########
#
RUN addgroup -g ${GID_DOCKER} docker && \
    adduser -u ${UID} -D -H ${USERNAME} && \
    addgroup ${USERNAME} docker


COPY --from=builder /srv/app/build/${PROJECT_NAME} ./${PROJECT_NAME}

CMD ["sh", "-c", "./${PROJECT_NAME}"]
