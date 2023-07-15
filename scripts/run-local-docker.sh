docker build -t template-api .

docker network create template-network || echo "template-network already created"
docker container stop template-api || echo "template-api's container already stopped"
docker container rm template-api || echo "template-api's container already removed"

docker run -d \
  --network host \
  --name template-api \
  template-api