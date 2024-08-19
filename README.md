# MVC
## SDSLabs MVC Assignment

### How to Run? (The Manual Way)
1. Install `go`.
2. Install `postgresql`.
3. Configure an `.env` according to `.env.sample`.
4. `go mod vendor`
5. `go mod tidy`
6. `go build -o mvc ./cmd/main.go`
7. `./mvc`
8. The website will be up at `localhost:5050`.

### How to Run? (Via Docker)
1. Configure an `.env` according to `.env.sample`.
2. `docker build -t mvc .`
3. `docker run -p <port>:5050 -d --rm mvc`
4. Site will be up at `localhost:<port>` since i cannot host it on `mvc.segfault.co` (because no money)