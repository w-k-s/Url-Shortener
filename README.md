# Short-url

Visit the website: [small.ml](https://small.ml)

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
- [ ] Move HTTPS to nginx
- [x] Upload binary to docker instead of source
- [ ] Document this montrosity

Useful Resources

1. [Automating HTTPS with LetsEncrypt using autocert library](https://blog.kowalczyk.info/article/Jl3G/https-for-free-in-go-with-little-help-of-lets-encrypt.html)
2. [Gaussian Distribution](https://stackoverflow.com/questions/29325069/how-to-generate-random-numbers-biased-towards-one-value-in-a-range)
3. [Backup Mongodb to Amazon S3](https://gist.github.com/eladnava/96bd9771cd2e01fb4427230563991c8d)
4. [Using `net/http` as front facing server](https://blog.cloudflare.com/exposing-go-on-the-internet/)