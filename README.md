# deposit server

## Run

Option one:
```bash
sudo docker run --name postgres2 --rm -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=postgres postgres
sudo docker build -t depositsrv:latest -f dockerfiles/Dockerfile_server .
sudo docker run --rm -it --name depositsrv --link postgres2 -p 8888:8888 depositsrv
```
Option two:
```bash
sudo docker-compose -f dockerfiles/docker-compose.yaml up -d
```
