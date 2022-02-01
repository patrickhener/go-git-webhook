# go-git-webhook

This little tool will run a webhook listener and then run a system command if the webhook triggers.

It will first check against the HMAC256 signature using the secret and then fire.

## Usage

```bash
WEBHOOK_IP=127.0.0.1 WEBHOOK_PORT=8282 WEBHOOK_SECRET=somesecret WEBHOOK_CMD=doathing.sh go-git-webhook
```

## systemd service template

```bash
[Unit]
Description=Webhook service
After=network.target

[Service]
Type=simple
User=patrick
Group=patrick
Restart=always
RestartSec=5s

Environment=WEBHOOK_IP=127.0.0.1
Environment=WEBHOOK_PORT=8282
Environment=WEBHOOK_SECRET=somesecret
Environment=WEBHOOK_CMD="/home/patrick/doathing.sh"

WorkingDirectory=/home/patrick/hook
ExecStart=/home/patrick/hook/go-git-webhook
SyslogIdentifier=git-webhook

[Install]
WantedBy=multi-user.target
```
