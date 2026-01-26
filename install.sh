# --- 配置路径 ---
BASE_PATH="/opt/cached_resources"
BIN_DIR="$BASE_PATH/bin"
PNPM_HOME="$BASE_PATH/pnpm"
PNPM_STORE="$BASE_PATH/pnpm_store"
export NVM_DIR="/opt/cached_resources/nvm"
mkdir -p "$NVM_DIR" "$PNPM_HOME" "$PNPM_STORE"
mkdir -p  $BIN_DIR 

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
    coverage


cd /tmp
# 下载 deb 包及其核心依赖
apt-get update
apt-get download libxml2-utils libxml2

# 解压 deb 包内容
find . -name "*.deb" -exec dpkg -x {} . \;

# 拷贝二进制文件
cp -f usr/bin/xmllint $BIN_DIR/
# 拷贝必要的动态链接库 (xmllint 运行需要 libxml2.so.2)
cp -f usr/lib/x86_64-linux-gnu/libxml2.so* $LIB_DIR/ 2>/dev/null || cp -f usr/lib/libxml2.so* $LIB_DIR/

chmod +x $BIN_DIR/xmllint

# 显式指定 NVM_DIR 环境变量，安装脚本会自动识别并装到这里
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | NVM_DIR="$NVM_DIR" bash

# --- 3. 加载 nvm 并安装 Node ---
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

nvm install 24

# --- 4. 关键：将 node/npm 软链接到系统能找到的地方 ---
# 这样你在扫描脚本里就不需要每次都加载 nvm 了
BIN_DIR="/opt/cached_resources/bin"
mkdir -p "$BIN_DIR"

# 获取 nvm 安装的实际 node 路径
CURRENT_NODE_BIN=$(nvm which 24)
CURRENT_NODE_DIR=$(dirname "$CURRENT_NODE_BIN")

ln -sf "$CURRENT_NODE_BIN" "$BIN_DIR/node"
ln -sf "$CURRENT_NODE_DIR/npm" "$BIN_DIR/npm"
ln -sf "$CURRENT_NODE_DIR/npx" "$BIN_DIR/npx"

# --- 5. 设置全局缓存到挂载点 ---
echo ">>> 正在持久化安装 pnpm..."
# 使用 Corepack 安装 pnpm (Node 24 自带)
export PNPM_HOME="$PNPM_HOME"
export PATH="$PNPM_HOME:$PATH"
corepack enable
corepack prepare pnpm@latest --activate

# 建立 pnpm 软链接到全局 BIN
ln -sf "$PNPM_HOME/pnpm" "$BIN_DIR/pnpm"
ln -sf "$PNPM_HOME/pnpx" "$BIN_DIR/pnpx"

# 关键：设置 pnpm 存储路径到挂载目录，实现真正的持久化缓存
"$BIN_DIR/pnpm" config set store-dir "$PNPM_STORE" --global
"$BIN_DIR/npm" config set cache "$BASE_PATH/npm_cache" --global

echo ">>> 环境预装检查:"
node -v
pnpm -v
xmllint --version | head -n 1
