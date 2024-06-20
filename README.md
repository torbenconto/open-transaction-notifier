# Open Transaction Notifier
SEC form 4 transaction notifier for insiders trading in the US stock market.

Uses only data from the SEC EDGAR database.

## How to run
1. Clone the repository
```bash
    git clone git@github.com:torbenconto/open-transaction-notifier.git && cd open-transaction-notifier
```
2. Configure


3. Configure the notifier to your liking
Modify config.toml to your heart's content.

3. Run the docker-compose
```bash
    docker compose up -d
```

Notifying options:
- NOTIFICATION_METHOD=discord_webhook (NOTIFICATION_INFO=your_discord_webhook_url)