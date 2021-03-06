# Maxmind GeoIP http server [![Build Status](https://travis-ci.org/OndrejIT/glok.svg?branch=master)](https://travis-ci.org/OndrejIT/glok)

### Install
  - Install [glide](https://github.com/Masterminds/glide) (Package Management for Golang)
  - glide install
  
  - or: go get ./... (not recommended)

### Test
  - go test $(glide novendor)

### Run
 - go run ./main.go (-h help, -d debug, -p Set server port, -f Set flag redirect, -m Set map redirect, -c config path)

### Usage
  - http://127.0.0.1:8888/v1/lookup - Return host IP info
    - {"country":"CZ","city":"Prague","latitude":xx.xxxx,"longitude":yy.yyyy}
  - http://127.0.0.1:8888/v1/lookup?ip=xxx.xxx.xxx.xxx - Return IP info
  - http://127.0.0.1:8888/v1/map?ip=xxx.xxx.xxx.xxx - Redirect to maps (Default google maps.)
  - http://127.0.0.1:8888/v1/flag?ip=xxx.xxx.xxx.xxx - Redirect to flag country (Default theffodora.)
