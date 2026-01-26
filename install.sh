# --- 配置路径 ---
BASE_PATH="/opt/cached_resources"
BIN_DIR="$BASE_PATH/bin"

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
