# Open Transaction Notifier
SEC form 4 transaction notifier for insiders trading in the US stock market.

## How to run
1. Clone the repository or pull the docker image (wip)

```bash
git clone git@github.com:torbenconto/open-transaction-notifier.git && cd open-transaction-notifier
```

2. Run the docker image

```bash
  docker-compose up
```

Notifying options:
- discord_webhook 

Need to split documents harvested into transactions before notifying. Each doc can
have multiple transactions.