# --- 配置路径 ---
BASE_PATH="/opt/cached_resources"
BIN_DIR="$BASE_PATH/bin"
export NVM_DIR="$BASE_PATH/nvm"
export PNPM_HOME="$BASE_PATH/pnpm"
BIN_DIR="$BASE_PATH/bin"
mkdir -p "$NVM_DIR" "$PNPM_HOME" "$BIN_DIR"
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
# --- 2. 下载并安装 nvm ---
# 显式指定 NVM_DIR，nvm 会将其内部组件安装到该目录
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
ln -sf "$node_path" "$BIN_DIR/node"
ln -sf "$(dirname "$node_path")/npm" "$BIN_DIR/npm"
ln -sf "$(dirname "$node_path")/npx" "$BIN_DIR/npx"

# 找到 corepack 激活后的 pnpm 真实路径并链接
PNPM_REAL_PATH=$(which pnpm)
ln -sf "$PNPM_REAL_PATH" "$BIN_DIR/pnpm"

# --- 验证结果 ---
echo "--------------------------------------"
echo "验证持久化工具链："
"$BIN_DIR/node" -v
"$BIN_DIR/pnpm" -v
echo "所有工具已链接至: $BIN_DIR"
echo "--------------------------------------"


echo ">>> 正在检查 jq 版本更新..."

# 1. 获取 GitHub 上的最新版本号 (例如 "jq-1.7.1")
LATEST_TAG=$(curl -s https://api.github.com/repos/stedolan/jq/releases/latest | grep '"tag_name":' | cut -d '"' -f 4)

# 2. 获取本地已安装的版本 (如果不存在则设为空)
CURRENT_VER=""
if [ -f "$BIN_DIR/jq" ]; then
    # jq --version 通常输出 "jq-1.7.1"
    CURRENT_VER=$("$BIN_DIR/jq" --version 2>/dev/null)
fi

# 3. 对比版本并决定是否下载
if [ "$LATEST_TAG" != "$CURRENT_VER" ]; then
    echo ">>> 检测到新版本: [本地: ${CURRENT_VER:-无}] -> [远程: $LATEST_TAG]"
    
    # 获取下载地址 (匹配 linux64)
    DOWNLOAD_URL=$(curl -s https://api.github.com/repos/stedolan/jq/releases/latest | \
                  grep "browser_download_url.*linux64" | \
                  cut -d '"' -f 4)

    if [ -n "$DOWNLOAD_URL" ]; then
        echo ">>> 正在更新 jq..."
        curl -L -s -S "$DOWNLOAD_URL" -o "$BIN_DIR/jq.tmp"
        
        # 验证下载的文件是否有效
        chmod +x "$BIN_DIR/jq.tmp"
        if "$BIN_DIR/jq.tmp" --version > /dev/null 2>&1; then
            mv "$BIN_DIR/jq.tmp" "$BIN_DIR/jq"
            echo "✅ jq 已成功更新至: $LATEST_TAG"
        else
            echo "❌ 下载的文件无效，取消更新。"
            rm -f "$BIN_DIR/jq.tmp"
        fi
    else
        echo "❌ 无法获取下载链接，跳过更新。"
    fi
else
    echo ">>> 当前 jq 已是最新版本 ($CURRENT_VER)，无需更新。"
fi
