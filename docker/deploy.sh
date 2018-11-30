#!/bin/bash
cp conf/db/edge.toml.tpl conf/db/edge.toml
sed -i "s/EDGE_DB_USERNAME/$EDGE_DB_USERNAME/g" conf/db/edge.toml
sed -i "s/EDGE_DB_PASSWORD/$EDGE_DB_PASSWORD/g" conf/db/edge.toml
sed -i "s/EDGE_DB_HOST/$EDGE_DB_HOST/g" conf/db/edge.toml
sed -i "s/EDGE_DB_PORT/$EDGE_DB_PORT/g" conf/db/edge.toml
sed -i "s/EDGE_DB_NAME/$EDGE_DB_NAME/g" conf/db/edge.toml
echo "$DOCKER_HUB_PASSWORD" | docker login -u="$DOCKER_HUB_USERNAME" "$DOCKER_HUB_PATH" --password-stdin
# docker build -t $DOCKER_HUB_PATH/$DOCKER_HUB_NAMESPACE/${DOCKER_HUB_REPOSITORY}_${LINUX_ARM32V7_IMAGE_NAME}:$TRAVIS_TAG -f docker/Dockerfile.${LINUX_ARM32V7_IMAGE_NAME} .
docker build -t $DOCKER_HUB_PATH/$DOCKER_HUB_NAMESPACE/${DOCKER_HUB_REPOSITORY}:$TRAVIS_TAG -f docker/Dockerfile .
docker image ls
# docker push $DOCKER_HUB_PATH/$DOCKER_HUB_NAMESPACE/${DOCKER_HUB_REPOSITORY}_${LINUX_ARM32V7_IMAGE_NAME}:$TRAVIS_TAG
docker push $DOCKER_HUB_PATH/$DOCKER_HUB_NAMESPACE/${DOCKER_HUB_REPOSITORY}:$TRAVIS_TAG
