$(boot2docker shellinit)
export GOPATH=~/.go
export PATH=~/.go/bin:$PATH

export DOORBOT_DATABASE_DRIVER=postgres
export DOORBOT_DATABASE_URL=postgres://postgres:dev@$(boot2docker ip)/?sslmode=disable
export DOORBOT_DATABASE_TRACE=true

export DOORBOT_NOTIFICATOR_ENABLED=true
export DOORBOT_NOTIFICATOR_EMAIL_ENABLED=true
export DOORBOT_NOTIFICATOR_POSTMARK_TOKEN=""
export DOORBOT_USER_ACCOUNTS_DOMAIN=.doorbot.dev
