# practical-crawler

pratical-crawler is a large-scale website crawler to collect all identifiable URL during processing. The original websites to start crawling are the URL seeds inside the config.go file, and the crawler will keep requesting the URLs it found recursively to implement a BFS-like policy.

# Prerequisites
- go 1.15
- docker 19.03.12
- docker-compose 1.26.2

# Getting Started

## Run Database
```
docker-compose up -d
```

## Run Crawler
```
go run cmd/main.go
```

## Configurations
Please refer to `config/config.go` for more advanced details.

# Benchmarks

- Network
    - Wifi 2.4Ghz
    - 72 Mbps

- Device
    - [Macbool pro 16 inches (basic build)](https://www.apple.com/tw/shop/buy-mac/macbook-pro/16-%E5%90%8B-%E5%A4%AA%E7%A9%BA%E7%81%B0%E8%89%B2-2.6ghz-6-%E6%A0%B8%E5%BF%83%E8%99%95%E7%90%86%E5%99%A8-512gb#)

- Performance
```
Benchmark 60 seconds, Left 4096 jobs
Encounter 806 errors, Recieve 1153 responses
Total 37187 records 
```

# Contributors
[RainrainWu](https://github.com/RainrainWu)