# --- 配置路径 ---
BASE_PATH="/opt/cached_resources"
export NVM_DIR="$BASE_PATH/sast/nodejs/nvm"
export PNPM_HOME="$BASE_PATH/sast/nodejs"
mkdir -p "$NVM_DIR" "$PNPM_HOME"

# 1. 创建目标目录
mkdir -p /opt/cached_resources/sast/python

# 2. 下载 Miniconda (这是获取 Python+Pip 最便捷的独立包)
# 如果是内网环境，请提前下载好该 sh 文件并放入制品仓库
wget https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh -O miniconda.sh

# 3. 安装到指定目录 (-b 为静默安装, -p 为路径)
bash miniconda.sh -b -u -p /opt/cached_resources/sast/python

# 4. 清理安装包
rm miniconda.sh

# 5. 配置软链接或环境变量
ln -sf /opt/cached_resources/sast/python/bin/python3 /usr/local/bin/python3
ln -sf /opt/cached_resources/sast/python/bin/pip3 /usr/local/bin/pip3

# 1. 创建一个专门存放库的目录
mkdir -p /opt/cached_resources/sast/python/user_packages

# 2. 设置环境变量（建议写入全局配置文件如 /etc/profile）
export PYTHONUSERBASE=/opt/cached_resources/sast/python/user_packages

# 3. 使用 --user 安装
/opt/cached_resources/sast/python/bin/pip3 install --user \
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
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | NVM_DIR="$NVM_DIR" bash

# 加载 nvm 环境（当前进程生效）
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

# --- 3. 安装 Node.js 24 ---
echo ">>> 正在安装 Node.js 24..."
nvm install 24

# 验证安装
node_path=$(nvm which 24)
echo "Node 实际路径: $node_path"
echo ">>> 清理旧的 pnpm 配置..."
rm -rm ~/.config/pnpm
export PATH="$PNPM_HOME:$PATH"

# --- 4. 激活 pnpm (参照官方 Corepack 方式) ---
echo ">>> 正在通过 Corepack 激活 pnpm..."
corepack enable pnpm
which pnpm
# 关键：配置 pnpm 的持久化存储和全局目录
# 确保 pnpm bin 目录也在持久化路径下
pnpm config set store-dir "$BASE_PATH/sast/nodejs/pnpm_store" --global
pnpm config set global-bin-dir "$PNPM_HOME" --global

# --- 5. 建立全局软链接 (方便 ut_scan.sh 直接调用) ---
# 这样你的扫描脚本只需把 /opt/cached_resources/bin 加入 PATH 即可
ln -sf "$node_path" "$BASE_PATH//sast/nodejs"
ln -sf "$(dirname "$node_path")/npm" "$BASE_PATH/sast/nodejs/npm"
ln -sf "$(dirname "$node_path")/npx" "$BASE_PATH/sast/nodejs/npx"

# 找到 corepack 激活后的 pnpm 真实路径并链接
PNPM_REAL_PATH=$(which pnpm)
ln -sf "$PNPM_REAL_PATH" "$BASE_PATH/sast/nodejs/pnpm"

# --- 验证结果 ---
echo "--------------------------------------"
echo "验证持久化工具链："
"$BASE_PATH/sast/nodejs/node" -v
"$BASE_PATH/sast/nodejs/pnpm" -v
echo "所有工具已链接至: $BASE_PATH/sast/nodejs"
echo "--------------------------------------"

mkdir -p $BASE_PATH/tools/common

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

echo ">>> 正在检查 Gitleaks 版本并尝试更新..."

# --- 2. 获取最新版本号 (使用 grep + sed 替代 jq) ---
# 访问 gitleaks 官方仓库获取最新 Release 标签
LATEST_TAG=$(curl -s https://api.github.com/repos/gitleaks/gitleaks/releases/latest | grep '"tag_name":' | sed -E 's/.*"v?([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo ">>> [警告] 无法通过 API 获取版本号，尝试下载静态版本 v8.23.1..."
    G_VER="8.23.1"
    G_TAG="v8.23.1"
else
    G_VER="$LATEST_TAG"
    G_TAG="v$LATEST_TAG"
    echo ">>> 发现 Gitleaks 最新版本: $G_TAG"
fi

# --- 3. 识别系统架构 ---
ARCH=$(uname -m)
case $ARCH in
    x86_64)  G_ARCH="x64" ;;
    aarch64) G_ARCH="arm64" ;;
    *)       echo ">>> [错误] 不支持的架构: $ARCH"; exit 1 ;;
esac

# --- 4. 构造下载链接并执行下载 ---
# Gitleaks 格式示例: gitleaks_8.23.1_linux_x64.tar.gz
DOWNLOAD_URL="https://github.com/gitleaks/gitleaks/releases/download/${G_TAG}/gitleaks_${G_VER}_linux_${G_ARCH}.tar.gz"

echo ">>> 正在下载: $DOWNLOAD_URL"
TMP_DIR=$(mktemp -d)
curl -L "$DOWNLOAD_URL" -o "$TMP_DIR/gitleaks.tar.gz"

if [ $? -eq 0 ]; then
    # 解压并将二进制文件移动到你的 BIN_DIR
    tar -xzf "$TMP_DIR/gitleaks.tar.gz" -C "$TMP_DIR"
    mkdir -p "$BASE_PATH/tools/gitleaks"
    mv -f "$TMP_DIR/gitleaks" "$BASE_PATH/tools/gitleaks/gitleaks"
    chmod +x "$BASE_PATH/tools/gitleaks/gitleaks"
    
    echo ">>> Gitleaks $G_ARCH 预置成功！版本信息："
    "$BASE_PATH/tools/gitleaks/gitleaks" version
else
    echo ">>> [错误] 下载失败，请检查网络。"
    rm -rf "$TMP_DIR"
    exit 1
fi

# 清理临时目录
rm -rf "$TMP_DIR"

echo ">>> 正在检查 Git Credential Manager (GCM) 版本并尝试更新..."

# --- 2. 获取最新版本号 (使用 grep + sed 替代 jq) ---
# 目标：从 tag_name 中提取 2.7.1
LATEST_TAG=$(curl -s https://api.github.com/repos/git-ecosystem/git-credential-manager/releases/latest | grep '"tag_name":' | sed -E 's/.*"v?([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo ">>> [警告] 无法通过 API 获取版本号，使用静态版本 2.7.1..."
    GCM_VER="2.7.1"
else
    GCM_VER="$LATEST_TAG"
    echo ">>> 发现 GCM 最新版本: v$GCM_VER"
fi

# --- 3. 识别系统架构 ---
ARCH=$(uname -m)
case $ARCH in
    x86_64)  GCM_ARCH="x64" ;;
    aarch64) GCM_ARCH="arm64" ;;
    *)       echo ">>> [错误] 不支持的架构: $ARCH"; exit 1 ;;
esac

# --- 4. 构造下载链接 ---
# 适配最新格式: gcm-linux-<arch>-<version>.tar.gz
# 示例: https://github.com/git-ecosystem/git-credential-manager/releases/download/v2.7.1/gcm-linux-x64-2.7.1.tar.gz
DOWNLOAD_URL="https://github.com/git-ecosystem/git-credential-manager/releases/download/v${GCM_VER}/gcm-linux-${GCM_ARCH}-${GCM_VER}.tar.gz"

echo ">>> 正在下载: $DOWNLOAD_URL"
TMP_DIR=$(mktemp -d)
curl -L "$DOWNLOAD_URL" -o "$TMP_DIR/gcm.tar.gz"

if [ $? -eq 0 ]; then
    # --- 5. 解压并安装 ---
    # GCM 非单文件程序，需独立目录存放依赖库
    GCM_HOME="$BASE_PATH/tools/gcm"
    mkdir -p "$GCM_HOME"
    
    # 清理旧版本并解压到专用目录
    rm -rf "$GCM_HOME"/*
    tar -xzf "$TMP_DIR/gcm.tar.gz" -C "$GCM_HOME"
    
    # 建立主程序软链接到 BIN_DIR，方便全局调用
    chmod +x "$GCM_HOME/git-credential-manager"
    
    echo ">>> GCM $GCM_ARCH 预置成功！版本信息："
    "$GCM_HOME/git-credential-manager" --version
else
    echo ">>> [错误] 下载失败，请确认 URL 是否正确。"
    rm -rf "$TMP_DIR"
    exit 1
fi

# 清理临时文件
rm -rf "$TMP_DIR"

echo ">>> 正在检查 Kustomize 版本并尝试更新..."

# --- 2. 获取最新版本号 (grep + sed) ---
# 官方仓库: kubernetes-sigs/kustomize
# 注意：kustomize 的 tag 格式通常是 kustomize/v5.x.x
LATEST_TAG_FULL=$(curl -s https://api.github.com/repos/kubernetes-sigs/kustomize/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# 提取纯版本号部分 (例如从 kustomize/v5.4.1 提取 v5.4.1)
LATEST_TAG=$(echo $LATEST_TAG_FULL | sed -E 's/.*(v[0-9]+\.[0-9]+\.[0-9]+).*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo ">>> [警告] 无法获取版本号，使用静态备份版本 v5.4.1..."
    K_TAG="v5.4.1"
else
    K_TAG="$LATEST_TAG"
    echo ">>> 发现 Kustomize 最新版本: $K_TAG"
fi

# --- 3. 识别系统架构 ---
ARCH=$(uname -m)
case $ARCH in
    x86_64)  K_ARCH="amd64" ;;
    aarch64) K_ARCH="arm64" ;;
    *)       echo ">>> [错误] 不支持的架构: $ARCH"; exit 1 ;;
esac

# --- 4. 构造下载链接 ---
# 格式示例: https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv5.4.1/kustomize_v5.4.1_linux_amd64.tar.gz
# 注意 URL 中的路径需要对 / 进行转义，或者直接使用拼接好的 tag
DOWNLOAD_URL="https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2F${K_TAG}/kustomize_${K_TAG}_linux_${K_ARCH}.tar.gz"

echo ">>> 正在下载: $DOWNLOAD_URL"
TMP_DIR=$(mktemp -d)
curl -L "$DOWNLOAD_URL" -o "$TMP_DIR/kustomize.tar.gz"

if [ $? -eq 0 ]; then
    # 5. 解压并安装
    tar -xzf "$TMP_DIR/kustomize.tar.gz" -C "$BASE_PATH/tools/common"
    chmod +x "$BASE_PATH/tools/common/kustomize"
    
    echo ">>> Kustomize $K_ARCH 预置成功！版本信息："
    "$BASE_PATH/tools/common/kustomize" version
else
    echo ">>> [错误] 下载失败，请检查网络。"
    rm -rf "$TMP_DIR"
    exit 1
fi

# 清理
rm -rf "$TMP_DIR"

echo ">>> 正在检查 SCC 版本并尝试更新..."
mkdir -p $BASE_PATH/tools/scc
# --- 2. 获取最新版本号 (grep + sed) ---
# 提取纯数字版本号，例如 3.6.0
LATEST_TAG=$(curl -s https://api.github.com/repos/boyter/scc/releases/latest | grep '"tag_name":' | sed -E 's/.*"v?([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo ">>> [警告] 无法获取版本号，使用静态版本 3.6.0..."
    S_VER="3.6.0"
else
    S_VER="$LATEST_TAG"
    echo ">>> 发现 SCC 最新版本: $S_VER"
fi

# --- 3. 识别系统架构 ---
# 根据你提供的链接，后缀固定为 Linux_x86_64 或 Linux_arm64
ARCH=$(uname -m)
case $ARCH in
    x86_64)  S_ARCH="Linux_x86_64" ;;
    aarch64) S_ARCH="Linux_arm64" ;;
    *)       echo ">>> [错误] 不支持的架构: $ARCH"; exit 1 ;;
esac

# --- 4. 构造下载链接 ---
# 适配格式: scc_Linux_x86_64.tar.gz
# 注意：SCC 的这个包名在 3.6.0 版本中连版本号都从文件名里去掉了
DOWNLOAD_URL="https://github.com/boyter/scc/releases/download/v${S_VER}/scc_Linux_x86_64.tar.gz"

# 如果你希望更健壮，支持多架构拼接：
# DOWNLOAD_URL="https://github.com/boyter/scc/releases/download/v${S_VER}/scc_${S_ARCH}.tar.gz"

echo ">>> 正在下载: $DOWNLOAD_URL"
TMP_DIR=$(mktemp -d)
curl -L "$DOWNLOAD_URL" -o "$TMP_DIR/scc.tar.gz"

if [ $? -eq 0 ]; then
    # --- 5. 解压并安装 ---
    tar -xzf "$TMP_DIR/scc.tar.gz" -C "$TMP_DIR"
    # 找到解压后的二进制文件并移动
    mv -f "$TMP_DIR/scc" "$BASE_PATH/tools/scc/scc"
    chmod +x "$BASE_PATH/tools/scc/scc"
    
    echo ">>> SCC $S_ARCH 预置成功！版本信息："
    "$BASE_PATH/tools/scc/scc" --version
else
    echo "下载失败"
    rm -rf "$TMP_DIR"
    exit 1
fi

# 清理
rm -rf "$TMP_DIR"
