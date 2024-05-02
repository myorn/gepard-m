# deposit server
Hi! Initially I thought that some client would be needed to test the server, but later it seemed like not too much of a good idea.
Anyways, sorry if it is sloppy in some places :)

**Important note:** this service was made with assumption that all sums are *__integers__* in which first two digits from the right side are cents

## How to run it:
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
