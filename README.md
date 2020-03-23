# Short-url

[![CircleCI](https://circleci.com/gh/w-k-s/UrlShortener.svg?style=svg)](https://circleci.com/gh/w-k-s/UrlShortener)
[![Go Report Card](https://goreportcard.com/badge/github.com/w-k-s/UrlShortener)](https://goreportcard.com/report/github.com/w-k-s/UrlShortener)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Visit the website: [shortest.ml](https://shortest.ml)

## Check List

1. Implement Backend API

- [x] Return Short Url
- [x] Return Long Url
- [x] Redirect using Short Url
- [x] Document usage on home page

2. Host locally on docker

- [x] Run on docker
- [x] Use Environmental variables to supply port and database connection
- [x] Create Docker-compose file

3. Publish on AWS

- [x] Set up remote mongodb
- [x] Create seperate local and prod docker-compose files
- [x] Configure app to run on HTTPS
- [x] Host on AWS

4. Refactor 
- [x] Move database logic to repository
- [x] Improve error responses. Go's error messages are short and useless.
- [x] Unit Tests using testify test suite
- [x] Set up CircleCI
- [x] Use swagger.io on the home screen
- [x] Consistent Naming: Url -> URL. Db -> DB. FullURL -> OriginalURL. ShortId -> ShortID
- [x] Remove swagger (slow and ugly)

5. Enhance
- [x] Send Cache headers when converting
- [x] Log all req/resp via middleware to db
- [x] Save db backups to s3, set up midnight task
- [x] Minify JS and CSS and serve via nginx
- [x] Move HTTPS to nginx
- [x] Upload binary to docker instead of source
- [ ] Document this montrosity

Useful Resources

1. [Automating HTTPS with LetsEncrypt using autocert library](https://blog.kowalczyk.info/article/Jl3G/https-for-free-in-go-with-little-help-of-lets-encrypt.html)
2. [Gaussian Distribution](https://stackoverflow.com/questions/29325069/how-to-generate-random-numbers-biased-towards-one-value-in-a-range)
3. [Backup Mongodb to Amazon S3](https://gist.github.com/eladnava/96bd9771cd2e01fb4427230563991c8d)
4. [Using `net/http` as front facing server](https://blog.cloudflare.com/exposing-go-on-the-internet/)
5. [Setting up https on ec2 ubuntu](https://blog.cloudboost.io/setting-up-an-https-sever-with-node-amazon-ec2-nginx-and-lets-encrypt-46f869159469)
5. [Setting up https on ami](https://coderwall.com/p/e7gzbq/https-with-certbot-for-nginx-on-amazon-linux)
6. [Building Go for Alpine Linux](https://www.blang.io/posts/2015-04_golang-alpine-build-golang-binaries-for-alpine-linux/)
7. [Migrating to go modules](https://blog.callr.tech/migrating-from-dep-to-go-1.11-modules/)
8. [Reduce Docker image size with multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build/)
9. [Only allow CD user to restart nginx](https://serverfault.com/questions/841099/systemd-grant-an-unprivileged-user-permission-to-alter-one-specific-service?newreg=c9a1362791da43f695cd5eb08b6e01c6)