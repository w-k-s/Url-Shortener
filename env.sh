touch .env
echo "TAG=$CIRCLE_BUILD_NUM" > .env
echo "MONGO_ADDRESS=$MONGO_ADDRESS" >> .env