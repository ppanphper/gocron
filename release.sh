#!/usr/bin/env bash

# command:
#   bash release.sh GITHUB_TOKEN
# 同时push到docker hub
#   bash release.sh GITHUB_TOKEN -d 1

GITHUB_TOKEN=$1

last_tag=$(git describe --tags "$(git rev-list --tags --max-count=1)")

printf "last_tag:%s\n" "${last_tag}"

CREATE_RESPONSE=$(curl \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  https://api.github.com/repos/ppanphper/gocron/releases \
  -d "{\"tag_name\":\"${last_tag}\",\"target_commitish\":\"master\",\"name\":\"${last_tag}\",\"body\":\"描述\",\"draft\":false,\"prerelease\":false,\"generate_release_notes\":false}")

echo "$CREATE_RESPONSE"

# python3
RELEASE_ID=$(echo "$CREATE_RESPONSE" | python -c "import sys, json; print(json.load(sys.stdin)['id'])")

# upload files
upload_assets() {
  ID=$1
  FILENAME=$2

  printf '\n upload %s\n' "${FILENAME}"

  GH_ASSET="https://uploads.github.com/repos/ppanphper/gocron/releases/${ID}/assets?name=$(basename "$FILENAME")"

  curl --data-binary @"${FILENAME}" \
    -H "Authorization: token ${GITHUB_TOKEN}" \
    -H "Content-Type: application/octet-stream" "${GH_ASSET}"
}

# shellcheck disable=SC2045
for f in $(ls gocron-package)
do
  upload_assets "${RELEASE_ID}" "gocron-package/${f}"
done

# shellcheck disable=SC2045
for f in $(ls gocron-node-package)
do
  upload_assets "${RELEASE_ID}" "gocron-node-package/${f}"
done


# Docker 相关变量
INPUT_PUSH_DOCKER=0 # 是否push到docker hub
DOCKER_USERNAME="" # 你的 Docker Hub 用户名 (或阿里云/其他仓库的命名空间)
DOCKER_REGISTRY="docker.io" # 如果是 Docker Hub 则留空，如果是阿里云则为 registry.cn-hangzhou.aliyuncs.com 等
DOCKER_IMAGE_NAME="gocron" # Docker 镜像名称，默认为 gocron
DOCKER_DOCKEFILE="Dockerfile" # Docker 镜像名称，默认为 gocron
DOCKER_ISLOGINED=0 # 是否已登录
 

# Docker 构建和推送
docker_build_and_push() {
    local IMAGE_FULL_NAME=""
    local DOCKER_IMAGE_TAG="${VERSION}" # Docker 镜像标签，默认为 VERSION 变量的值
    if [ -n "$DOCKER_REGISTRY" ]; then
        IMAGE_FULL_NAME="$DOCKER_REGISTRY/$DOCKER_USERNAME/$DOCKER_IMAGE_NAME:$DOCKER_IMAGE_TAG"
    else
        IMAGE_FULL_NAME="$DOCKER_USERNAME/$DOCKER_IMAGE_NAME:$DOCKER_IMAGE_TAG"
    fi

    
    if docker pull "$IMAGE_FULL_NAME" --dry-run > /dev/null 2>&1; then
      echo "Delete Docker image: $IMAGE_FULL_NAME"
      docker rmi "$IMAGE_FULL_NAME"
      if [[ $? -eq 0 ]]; then
        echo "Image $IMAGE_FULL_NAME deleted successfully."
      else
        echo "Failed to delete image $IMAGE_FULL_NAME."
      fi
    fi

    echo "Building Docker image: $IMAGE_FULL_NAME"
    
    docker build -t "$IMAGE_FULL_NAME" -f "$DOCKER_DOCKEFILE" .

    if [ "$DOCKER_ISLOGINED" -eq 0 ]; then
      echo "Logging in to Docker registry..."
      if [ -n "$DOCKER_REGISTRY" ]; then
          docker login --username="$DOCKER_USERNAME" "$DOCKER_REGISTRY"
      else
          docker login
      fi
      DOCKER_ISLOGINED=1
    fi

    echo "Pushing Docker image: $IMAGE_FULL_NAME"
    docker push "$IMAGE_FULL_NAME"

    echo "Docker build and push completed."
}

# 解析命令行参数
while getopts "d:" OPT; do
  case ${OPT} in
    d) INPUT_PUSH_DOCKER=$OPTARG ;; # 设置 INPUT_PUSH_DOCKER 的值
    *)
      ;;
  esac
done

docker_build_and_push_web() {
  DOCKER_IMAGE_NAME="gocron"
  DOCKER_DOCKEFILE="Dockerfile"
  docker_build_and_push
}

docker_build_and_push_node() {
  DOCKER_IMAGE_NAME="gocron-node"
  DOCKER_DOCKEFILE="Dockerfile-node"
  docker_build_and_push
}

if [ "$INPUT_PUSH_DOCKER" -eq 1 ]; then
    docker_build_and_push_web
    docker_build_and_push_node
fi
