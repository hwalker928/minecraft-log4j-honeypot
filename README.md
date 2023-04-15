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

### Using Docker

```
git clone https://github.com/hwalker928/minecraft-log4j-honeypot.git
cd minecraft-log4j-honeypot
docker build -t minecraft-log4j-honeypot .
mkdir payloads
docker run --rm -it --mount type=bind,source="${PWD}/payloads",target=/payloads --user=`id -u` -p 25565:25565 minecraft-log4j-honeypot
```

## Configuring AbuseIPDB
