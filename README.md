# Open Transaction Notifier
SEC form 4 transaction notifier for insiders trading in the US stock market.

Uses only data from the SEC EDGAR database.

## How to run
1. Pull the docker image

```bash
    docker pull torbenconto/open-transaction-notifier:latest
```

2. Configure and run the docker container

```bash
    docker run -d
     -e NOTIFICATION_METHOD=your_notification_method
        -e NOTIFICATION_INFO=your_notification_info (eg: webhook url)
     open-transaction-notifier:latest
```

Notifying options:
- discord_webhook