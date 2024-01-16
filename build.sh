docker login repo.name
docker build -t davidacarter/cert-scanner .
docker tag davidacarter/cert-scanner repo.name/davidacarter/cert-scanner:latest
docker push repo.name/davidacarter/cert-scanner:latest