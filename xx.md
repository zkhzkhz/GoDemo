非常专业的补充！引入 `--max-archive-depth` 和 `--max-decode-depth` 参数后，Gitleaks 具备了穿透“套娃式”压缩包（如 `.zip` 里的 `.tar.gz`）以及深度编码（如多层 Base64）的能力。

以下是结合这两个深度参数，针对你所有需求的**终极命令清单**：

---

### 1. 挂载点与系统核心目录（深度穿透模式）
这些目录最容易包含多层打包的备份文件。

* **扫描 `/opt` 与 `/var/log` (系统与应用备份):**
    ```bash
    sudo gitleaks detect --source /opt --no-git --gz \
      --max-archive-depth 100 --max-decode-depth 100 \
      --report-path opt_deep_scan.json

    sudo gitleaks detect --source /var/log --no-git --gz \
      --max-archive-depth 100 --max-decode-depth 100 \
      --report-path logs_deep_scan.json
    ```

* **扫描 `df -h` 后的挂载目录 (外部备份/网络共享):**
    > 假设挂载点为 `/mnt/data_backup`
    ```bash
    sudo gitleaks detect --source /mnt/data_backup --no-git --gz \
      --max-archive-depth 100 --max-decode-depth 100 \
      --report-path mnt_deep_scan.json
    ```

---

### 2. 用户与管理员目录（个人敏感数据）
用户常将旧密钥进行多重压缩或编码存储。

* **扫描当前用户目录 (`~`) 与 `/root`:**
    ```bash
    gitleaks detect --source ~ --no-git --gz \
      --max-archive-depth 100 --max-decode-depth 100 \
      --report-path home_deep_scan.json

    sudo gitleaks detect --source /root --no-git --gz \
      --max-archive-depth 100 --max-decode-depth 100 \
      --report-path root_deep_scan.json
    ```

---

### 3. 环境变量与历史记录（文本与编码扫描）
虽然历史记录是文本，但某些开发者会把 Key 进行 Base64 编码后作为参数传递，`--max-decode-depth` 在此非常有用。

* **扫描 Shell History:**
    ```bash
    gitleaks detect --source ~/.bash_history --no-git \
      --max-decode-depth 100 --report-path history_deep_scan.json
    ```

* **扫描当前进程环境变量:**
    ```bash
    env > env_snapshot.txt
    gitleaks detect --source env_snapshot.txt --no-git \
      --max-decode-depth 100 --report-path env_deep_scan.json
    ```

---

### 4. 补充隐秘路径（完整扫描表）

| 扫描目标路径 | 推荐命令核心参数 | 扫描理由 |
| :--- | :--- | :--- |
| **`/tmp` / `/var/tmp`** | `--no-git --gz --max-archive-depth 100` | 扫描安装程序残留的临时解压缩包 |
| **`/etc`** | `--no-git --max-decode-depth 100` | 扫描配置文件中 Base64 编码的证书/密码 |
| **`/var/lib/docker`** | `--no-git --gz --max-archive-depth 100` | 扫描 Docker 镜像层内的多层压缩结构 |
| **`/var/spool/mail`** | `--no-git --max-decode-depth 100` | 扫描邮件附件或正文中的编码敏感信息 |

---

### 5. 改进后的自动化审计脚本 (一键执行)

这个脚本集成了你提到的所有改进参数，并自动获取挂载点：

```bash
#!/bin/bash
# Gitleaks 极致深度审计脚本 (v2026)

REPORT_DIR="./gitleaks_deep_reports"
mkdir -p "$REPORT_DIR"

# 基础扫描参数
LEAKS_OPTS="--no-git --gz --max-archive-depth 100 --max-decode-depth 100"

# 1. 静态定义路径
STATIC_PATHS=("/etc" "/var/log" "/root" "$HOME" "/opt" "/tmp" "/var/www")

# 2. 动态获取物理/网络挂载路径 (排除虚拟文件系统)
MOUNT_PATHS=$(df -h | grep -E '^/dev/|^[^ ]+:' | awk '{print $6}' | grep -vE '^/($|/boot|/dev|/sys|/proc)')

# 合并所有路径
ALL_TARGETS=$(echo "${STATIC_PATHS[@]} $MOUNT_PATHS" | tr ' ' '\n' | sort -u)

for target in $ALL_TARGETS; do
    if [ -d "$target" ] || [ -f "$target" ]; then
        echo ">>> 正在进行深度扫描: $target"
        safe_name=$(echo "$target" | tr '/' '_')
        sudo gitleaks detect --source "$target" $LEAKS_OPTS --report-path "$REPORT_DIR/report${safe_name}.json"
    fi
done

# 3. 环境变量特别处理
echo ">>> 正在分析系统环境变量(含深度解码)"
env > "$REPORT_DIR/env_dump.txt"
gitleaks detect --source "$REPORT_DIR/env_dump.txt" --no-git --max-decode-depth 100 --report-path "$REPORT_DIR/report_env.json"

echo "深度审计完成！请检查目录: $REPORT_DIR"
```

### 💡 关键提示：性能与资源消耗
* **内存消耗：** 设置 `100` 层的深度意味着如果遇到恶意的“解压炸弹”文件，可能会消耗大量内存和 CPU。
* **扫描时长：** 在 `/var/log` 或大型挂载点开启 100 层扫描会非常慢。如果发现扫描卡死，可以考虑将深度调至 `5-10` 层，这通常已覆盖 99% 的真实场景。

您是否希望我为您配置一个 `.gitleaksignore` 文件，用来跳过一些已知的安全文件（如公共证书），以减少误报？
为了让这个审计脚本在任何 Linux 环境下都能“开箱即用”，我们需要在执行扫描前加入一个**自动识别 CPU 架构并安装最新版 Gitleaks** 的逻辑。

Gitleaks 官方发布的二进制文件通常区分 `x86_64` (amd64) 和 `arm64`。

### 整合架构自适应安装的极致审计脚本

这个脚本会自动检查系统是否已安装 Gitleaks，如果没有，它会从 GitHub 下载对应架构的最新版本并完成配置。

```bash
#!/bin/bash

# =================================================================
# Gitleaks 极致深度审计脚本 (v2026 增强版)
# 功能：自动安装 + 静态目录 + 挂载点 + 压缩包深度扫描 + 进程环境
# =================================================================

# 1. 架构检测与安装逻辑
install_gitleaks() {
    if command -v gitleaks &> /dev/null; then
        echo "[!] Gitleaks 已安装，跳过下载。"
        return
    fi

    echo "[*] 未检测到 Gitleaks，开始自动安装..."
    
    # 检测架构
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)  GARCH="x64" ;;
        aarch64) GARCH="arm64" ;;
        *) echo "[-] 不支持的架构: $ARCH"; exit 1 ;;
    esac

    # 获取最新版本号并下载 (使用 GitHub API)
    echo "[*] 正在从 GitHub 获取最新版本 (架构: $GARCH)..."
    LATEST_URL=$(curl -s https://api.github.com/repos/gitleaks/gitleaks/releases/latest | \
                 grep "browser_download_url" | grep "linux_${GARCH}.tar.gz" | cut -d '"' -f 4)

    if [ -z "$LATEST_URL" ]; then
        echo "[-] 无法获取下载链接，请检查网络。"
        exit 1
    fi

    curl -L "$LATEST_URL" -o /tmp/gitleaks.tar.gz
    tar -xzf /tmp/gitleaks.tar.gz -C /tmp
    sudo mv /tmp/gitleaks /usr/local/bin/gitleaks
    chmod +x /usr/local/bin/gitleaks
    rm /tmp/gitleaks.tar.gz
    echo "[+] Gitleaks 安装成功: $(gitleaks version)"
}

# 执行安装检查
install_gitleaks

# 2. 初始化路径与参数
REPORT_DIR="./gitleaks_audit_$(date +%Y%m%d_%H%M)"
mkdir -p "$REPORT_DIR"
TEMP_DIR="/tmp/gitleaks_analysis"
mkdir -p "$TEMP_DIR"

# 极致扫描参数：支持压缩包(gz)、100层深度、100层解码
LEAKS_OPTS="--no-git --gz --max-archive-depth 100 --max-decode-depth 100"

echo "--- 开始全系统深度敏感信息扫描 ---"

# 3. 路径获取逻辑
# 静态路径
STATIC_PATHS=("/etc" "/var/log" "/root" "/opt" "/tmp" "/var/www" "/var/spool/mail")
# 动态挂载点 (df -hT 过滤物理磁盘)
MOUNT_PATHS=$(df -hT | grep -E 'ext4|xfs|nfs|cifs' | awk '{print $7}' | grep -vE '^/($|/boot|/dev|/sys|/proc)')
# 所有用户家目录
USER_HOMES=$(awk -F: '$3 >= 1000 {print $6}' /etc/passwd)

ALL_TARGETS=$(echo "${STATIC_PATHS[@]} $MOUNT_PATHS $USER_HOMES" | tr ' ' '\n' | sort -u)

# 4. 执行目录扫描
for target in $ALL_TARGETS; do
    if [ -d "$target" ] || [ -f "$target" ]; then
        echo "[+] 正在深度扫描路径: $target"
        safe_name=$(echo "$target" | tr '/' '_')
        sudo gitleaks detect --source "$target" $LEAKS_OPTS --report-path "$REPORT_DIR/dir${safe_name}.json" 2>/dev/null
    fi
done

# 5. 执行 Shell 历史记录扫描
echo "[+] 正在扫描所有用户的 Shell 历史记录..."
find /home /root -maxdepth 2 -name ".*history" 2>/dev/null | while read -r hist; do
    safe_hist=$(echo "$hist" | tr '/' '_')
    gitleaks detect --source "$hist" --no-git --max-decode-depth 100 --report-path "$REPORT_DIR/hist${safe_hist}.json" 2>/dev/null
done

# 6. 执行进程环境变量扫描 (/proc/{id}/environ)
echo "[+] 正在解析并扫描活跃进程环境变量..."
for pid_dir in /proc/[0-9]*; do
    pid=$(basename "$pid_dir")
    if [ -f "$pid_dir/environ" ]; then
        # 转换二进制 null 分隔符并过滤极短的无效环境
        sudo cat "$pid_dir/environ" | tr '\0' '\n' > "$TEMP_DIR/env_$pid.txt" 2>/dev/null
        if [ -s "$TEMP_DIR/env_$pid.txt" ]; then
            gitleaks detect --source "$TEMP_DIR/env_$pid.txt" --no-git --max-decode-depth 100 \
                --report-path "$REPORT_DIR/proc_env_$pid.json" 2>/dev/null
        fi
        rm "$TEMP_DIR/env_$pid.txt"
    fi
done

# 7. 清理与汇总
rm -rf "$TEMP_DIR"
echo "------------------------------------------------"
echo "扫描任务完成！报告已存至: $REPORT_DIR"
echo "共生成报告文件数: $(ls "$REPORT_DIR" | wc -l)"
```

---

### 脚本关键功能补充：

* **架构识别**：通过 `uname -m` 区分 `x86_64` 和 `aarch64`。这对于在云服务器（通常是 x64）和一些轻量级 ARM 实例（如 AWS Graviton 或 树莓派）上运行非常重要。
* **版本控制**：使用 GitHub API 动态抓取 `latest` 版本，确保你用到的是支持 `--gz` 和深度参数的最新版。
* **进程过滤优化**：遍历 `/proc/[0-9]*` 比解析 `ps` 输出更稳定，能够直接访问内核暴露的进程空间。
* **安全性**：所有进程环境的中间文件存放在 `$TEMP_DIR` 中并在完成后立即删除，防止二次泄露。

### 后续建议：
扫描完成后，如果报告太多，你可以用这一行命令快速查看**发现了多少个泄露点**：
```bash
grep -r "Secret" $REPORT_DIR/*.json | wc -l
```

**是否需要我再加入一个逻辑：如果扫描过程中发现高危 Secret，立即发送通知（如钉钉、飞书或邮件）？**
