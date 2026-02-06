# --- 配置路径 ---
BASE_PATH="/opt/cached_resources"
export NVM_DIR="$BASE_PATH/nvm"
export PNPM_HOME="$BASE_PATH/pnpm"
mkdir -p "$NVM_DIR" "$PNPM_HOME" 

# 1. 创建目标目录
mkdir -p /opt/cached_resources/python

# 2. 下载 Miniconda (这是获取 Python+Pip 最便捷的独立包)
# 如果是内网环境，请提前下载好该 sh 文件并放入制品仓库
wget https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh -O miniconda.sh

# 3. 安装到指定目录 (-b 为静默安装, -p 为路径)
bash miniconda.sh -b -u -p /opt/cached_resources/python

# 4. 清理安装包
rm miniconda.sh

# 5. 配置软链接或环境变量
ln -sf /opt/cached_resources/python/bin/python3 /usr/local/bin/python3
ln -sf /opt/cached_resources/python/bin/pip3 /usr/local/bin/pip3

# 1. 创建一个专门存放库的目录
mkdir -p /opt/cached_resources/python/user_packages

# 2. 设置环境变量（建议写入全局配置文件如 /etc/profile）
export PYTHONUSERBASE=/opt/cached_resources/python/user_packages

# 3. 使用 --user 安装
/opt/cached_resources/python/bin/pip3 install --user \
    pytest \
    pytest-cov \
    coverage \
    diff-cover 

# 1. 安装 gocov 和 gocov-xml
export GOPATH=/opt/cached_resources/sast/gopath
export PATH=$GOPATH/bin:$PATH

go install github.com/axw/gocov/gocov@latest
go install github.com/AlekSi/gocov-xml@latest

cd /tmp
# 下载 deb 包及其核心依赖
apt-get update
apt-get download libxml2-utils libxml2

# 解压 deb 包内容
find . -name "*.deb" -exec dpkg -x {} . \;

# 拷贝二进制文件
cp -f usr/bin/xmllint $BASE_PATH/tools/common
# 拷贝必要的动态链接库 (xmllint 运行需要 libxml2.so.2)
cp -f usr/lib/x86_64-linux-gnu/libxml2.so* $LIB_DIR/ 2>/dev/null || cp -f usr/lib/libxml2.so* $LIB_DIR/

chmod +x $BASE_PATH/tools/common/xmllint

# --- 2. 下载并安装 nvm ---
# 显式指定 NVM_DIR，nvm 会将其内部组件安装到该目录
mkdir -p "$BASE_PATH/nodejs"
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | NVM_DIR="$NVM_DIR" bash

# 加载 nvm 环境（当前进程生效）
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

# --- 3. 安装 Node.js 24 ---
echo ">>> 正在安装 Node.js 24..."
nvm install 24

# 验证安装
node_path=$(nvm which 24)
echo "Node 实际路径: $node_path"

# --- 4. 激活 pnpm (参照官方 Corepack 方式) ---
echo ">>> 正在通过 Corepack 激活 pnpm..."
corepack enable pnpm

# 关键：配置 pnpm 的持久化存储和全局目录
# 确保 pnpm bin 目录也在持久化路径下
export PATH="$PNPM_HOME:$PATH"
pnpm config set store-dir "$BASE_PATH/pnpm_store" --global
pnpm config set global-bin-dir "$PNPM_HOME" --global

# --- 5. 建立全局软链接 (方便 ut_scan.sh 直接调用) ---
# 这样你的扫描脚本只需把 /opt/cached_resources/bin 加入 PATH 即可
ln -sf "$node_path" "$BASE_PATH/nodejs"
ln -sf "$(dirname "$node_path")/npm" "$BASE_PATH/nodejs/npm"
ln -sf "$(dirname "$node_path")/npx" "$BASE_PATH/nodejs/npx"

# 找到 corepack 激活后的 pnpm 真实路径并链接
PNPM_REAL_PATH=$(which pnpm)
ln -sf "$PNPM_REAL_PATH" "$BASE_PATH/nodejs/pnpm"

# --- 验证结果 ---
echo "--------------------------------------"
echo "验证持久化工具链："
"$BASE_PATH/nodejs/node" -v
"$BASE_PATH/nodejs/pnpm" -v
echo "所有工具已链接至: $BASE_PATH/nodejs"
echo "--------------------------------------"

echo ">>> 检查 jq 版本并尝试更新 (架构: amd64)..."

# 1. 获取最新 Release 的 Tag Name (例如 "jq-1.7.1")
# 使用 jqlang 组织下的新仓库地址
LATEST_TAG=$(curl -s https://api.github.com/repos/jqlang/jq/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo ">>> [警告] 无法获取最新版本号，尝试使用静态备份地址..."
    DOWNLOAD_URL="https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64"
else
    # 2. 拼接 amd64 专用的下载 URL
    # 注意：新版本文件名通常为 jq-linux-amd64 或 jq-linux64
    DOWNLOAD_URL="https://github.com/jqlang/jq/releases/download/${LATEST_TAG}/jq-linux-amd64"
    echo ">>> 发现最新版本: $LATEST_TAG"
fi

# 3. 下载并覆盖旧版本
# -N 选项可以检查服务器文件是否比本地新，节省带宽
curl -L "$DOWNLOAD_URL" -o "$BASE_PATH/tools/common/jq"

if [ $? -eq 0 ]; then
    chmod +x "$BASE_PATH/tools/common/jq"
    echo ">>> jq amd64 预置成功！版本信息："
    "$BASE_PATH/tools/common/jq" --version
else
    echo ">>> [错误] 下载失败，请检查网络连接。"
    exit 1
fi
