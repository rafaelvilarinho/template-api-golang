APP_NAME="template-api"

docker build -t $APP_NAME .

docker container stop $APP_NAME || echo "$APP_NAME's container already stopped"
docker container rm $APP_NAME || echo "$APP_NAME's container already removed"

# host network don't work on Windows or Mac, only Linux
docker run -d \
  --network host \
  --name $APP_NAME \
  $APP_NAME