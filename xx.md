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
