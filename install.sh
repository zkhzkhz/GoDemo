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


# 先用 yum 下载，然后提取二进制，确保跨容器可用
cd /tmp
yum install -y yum-utils
yumdownloader libxml2
rpm2cpio libxml2-*.rpm | cpio -idmv
cp -f usr/bin/xmllint $BIN_DIR/
# 赋予执行权限
chmod +x $BIN_DIR/xmllint
