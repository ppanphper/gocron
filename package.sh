#!/usr/bin/env bash
 
# 生成压缩包 xx.tar.gz或xx.zip
# 使用 ./package.sh -a amd664 -p linux -v v2.0.0

# 任何命令返回非0值退出
set -o errexit
# 使用未定义的变量退出
set -o nounset
# 管道中任一命令执行失败退出
set -o pipefail

eval $(go env)
 
# 二进制文件名
BINARY_NAME=''
# main函数所在文件
MAIN_FILE=""
 
# 提取git最新tag作为应用版本
VERSION=''
# 最新git commit id
GIT_COMMIT_ID=''
 
# 外部输入的系统
INPUT_OS=()
# 外部输入的架构
INPUT_ARCH=()
# 未指定OS，默认值
DEFAULT_OS=${GOHOSTOS}
# 未指定ARCH,默认值
DEFAULT_ARCH=${GOHOSTARCH}
# 支持的系统
SUPPORT_OS=(linux darwin windows)
# 支持的架构
SUPPORT_ARCH=(386 amd64 arm64)

 
# 编译参数
LDFLAGS=''
# 需要打包的文件
INCLUDE_FILE=()
# 打包文件生成目录
PACKAGE_DIR=''
# 编译文件生成目录
BUILD_DIR=''

# 获取git 最新tag name
git_latest_tag() {
    local COMMIT_ID=""
    local TAG_NAME=""
    COMMIT_ID=`git rev-list --tags --max-count=1`
    TAG_NAME=`git describe --tags "${COMMIT_ID}"`
 
    echo ${TAG_NAME}
}
 
# 获取git 最新commit id
git_latest_commit() {
    echo "$(git rev-parse --short HEAD)"
}
 
# 打印信息
print_message() {
    echo "$1"
}
 
# 打印信息后推出
print_message_and_exit() {
    if [[ -n $1 ]]; then
        print_message "$1"
    fi
    exit 1
}
 
# 设置系统、CPU架构
set_os_arch() {
    if [[ ${#INPUT_OS[@]} = 0 ]];then
        INPUT_OS=("${DEFAULT_OS}")
    fi
 
    if [[ ${#INPUT_ARCH[@]} = 0 ]];then
        INPUT_ARCH=("${DEFAULT_ARCH}")
    fi
 
    for OS in "${INPUT_OS[@]}"; do
        if [[  ! "${SUPPORT_OS[*]}" =~ ${OS} ]]; then
            print_message_and_exit "不支持的系统${OS}"
        fi
    done
 
    for ARCH in "${INPUT_ARCH[@]}"; do
        if [[ ! "${SUPPORT_ARCH[*]}" =~ ${ARCH} ]]; then
            print_message_and_exit "不支持的CPU架构${ARCH}"
        fi
    done
}
 
# 初始化
init() {
    set_os_arch
 
    if [[ -z "${VERSION}" ]];then
        VERSION=`git_latest_tag`
    fi
    GIT_COMMIT_ID=`git_latest_commit`
    LDFLAGS="-w -X 'main.AppVersion=${VERSION}' -X 'main.BuildDate=`date '+%Y-%m-%d %H:%M:%S'`' -X 'main.GitCommit=${GIT_COMMIT_ID}'"
 
    PACKAGE_DIR=${BINARY_NAME}-package
    BUILD_DIR=${BINARY_NAME}-build
 
    if [[ -d ${BUILD_DIR} ]];then
        rm -rf ${BUILD_DIR}
    fi
    if [[ -d ${PACKAGE_DIR} ]];then
        rm -rf ${PACKAGE_DIR}
    fi
 
    mkdir -p ${BUILD_DIR}
    mkdir -p ${PACKAGE_DIR}
}
 
# 编译
build() {
    local FILENAME=''
    for OS in "${INPUT_OS[@]}";do
        for ARCH in "${INPUT_ARCH[@]}";do
            if [[ "${OS}" = "windows"  ]];then
                FILENAME=${BINARY_NAME}.exe
            else
                FILENAME=${BINARY_NAME}
            fi
            env CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${BINARY_NAME}-${OS}-${ARCH}/${FILENAME} ${MAIN_FILE}
        done
    done
}
 
# 打包
package_binary() {
    cd ${BUILD_DIR}
 
    for OS in "${INPUT_OS[@]}";do
        for ARCH in "${INPUT_ARCH[@]}";do
        package_file ${BINARY_NAME}-${OS}-${ARCH}
        if [[ "${OS}" = "windows" ]];then
            zip -rq ../${PACKAGE_DIR}/${BINARY_NAME}-${VERSION}-${OS}-${ARCH}.zip ${BINARY_NAME}-${OS}-${ARCH}
        else
            tar czf ../${PACKAGE_DIR}/${BINARY_NAME}-${VERSION}-${OS}-${ARCH}.tar.gz ${BINARY_NAME}-${OS}-${ARCH}
        fi
        done
    done
 
    cd ${OLDPWD}
}
 
# 打包文件
package_file() {
    if [[ "${#INCLUDE_FILE[@]}" = "0" ]];then
        return
    fi
    for item in "${INCLUDE_FILE[@]}"; do
            cp -r ../${item} $1
    done
}
 
# 清理
clean() {
    if [[ -d ${BUILD_DIR} ]];then
        rm -rf ${BUILD_DIR}
    fi
}
 
# 运行
run() {
    init
    build
    package_binary
    clean
}

package_gocron() {
    BINARY_NAME='gocron'
    MAIN_FILE="./cmd/gocron/gocron.go"
    INCLUDE_FILE=()


    run
}

package_gocron_node() {
    BINARY_NAME='gocron-node'
    MAIN_FILE="./cmd/node/node.go"
    INCLUDE_FILE=()

    run
}


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
 
# p 平台 linux darwin windows
# a 架构 386 amd64
# v 版本号  默认取git最新tag
while getopts "p:a:v:d:" OPT;
do
    case ${OPT} in
    p) IFS=',' read -r -a INPUT_OS <<< "$OPTARG" ;;
    a) IFS=',' read -r -a INPUT_ARCH <<< "$OPTARG" ;;
    d) INPUT_PUSH_DOCKER=$OPTARG ;; # 设置 INPUT_PUSH_DOCKER 的值
    v) VERSION=$OPTARG ;;
    *) ;;
    esac
done
 
package_gocron
package_gocron_node


if [ "$INPUT_PUSH_DOCKER" -eq 1 ]; then
    echo "push docker"
    docker_build_and_push_web
    docker_build_and_push_node
fi
