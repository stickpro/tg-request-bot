## Simple telegram bot for parse message and resend to any service

### Install

You can use docker or using systemd start

for start from docker you need start container

``` docker compose up -d```

or build app 

``` make run build ```

on dir ./.bin/go-tgbot building app for start him from systemd using sim,ple config

```
[Unit]
Description=Telegram bot service
After=network.target

[Service]
Restart=always
RestartSec=3
WorkingDirectory=/home/server/telegram_bot/
ExecStart=/home/server/telegram_bot/.bin/go-tgbot

[Install]
WantedBy=multi-user.target
```
