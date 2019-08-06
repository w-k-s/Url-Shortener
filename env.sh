touch .env
echo "TAG=$CIRCLE_BUILD_NUM" > .env
echo "DB_CONN_STRING=$DB_CONN_STRING" >> .env