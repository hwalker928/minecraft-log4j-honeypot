# Minecraft Log4j Honeypot

This honeypot runs a fake Minecraft server (1.7.2 - 1.19.4 without snapshots) waiting to be exploited. Payload classes are saved to `payloads/` directory.

## Requirements

- Golang 1.16+

## Running the Honeypot

### Natively

```
git clone https://github.com/hwalker928/minecraft-log4j-honeypot.git
cd minecraft-log4j-honeypot
go build .
./minecraft-log4j-honeypot
```

### Using Docker Compose

```
git clone https://github.com/hwalker928/minecraft-log4j-honeypot.git
cd minecraft-log4j-honeypot
cp config.example.json config.json
touch database.db
mkdir payloads
docker compose up -d
```

## Configuring AbuseIPDB

You can generate an API key from AbuseIPDB [here](https://www.abuseipdb.com/account/api). Once you have an API key, you can set it in the `config.json` file.

## Configuring Discord Webhook

You can create a Discord webhook [here](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks). Once you have a webhook, you can set it in the `config.json` file.
