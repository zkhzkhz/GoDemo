set +e
#!/bin/bash

# 配置工作目录和工具版本
WORKDIR=/opt/cached_resources/sast
SPOTBUGS_HOME=/opt/cached_resources/sast/spotbugs-4.7.3
# 使用 sed 替换 '/' 为 '@'
sanitizedBranch=$(echo "$codeBranch" | tr '/' '@')

export GOROOT=${WORKDIR}/go
export GOPATH=${WORKDIR}/gopath
export JAVA_HOME=$WORKDIR/jdk18
export PATH=$JAVA_HOME/bin:$GOROOT/bin:$GOPATH/bin:/opt/cached_resources/sast/nodejs:/opt/cached_resources/sast/python/bin:/opt/cached_resources/tools/common:/opt/cached_resources/tools/gcm:/opt/cached_resources/tools/scc:/opt/cached_resources/sast/python/user_packages/bin:$PATH
export PYTHONUSERBASE=/opt/cached_resources/sast/python/user_packages
export MAVEN_OPTS="-Dmaven.repo.local=/opt/cached_resources/sast/maven_repository"

# 初始化全局状态
GLOBAL_SUCCESS=true
SUMMARY_TABLE=""

calculate_go_incremental_readable() {
  local out_path=$1
  local threshold=80.0
  local lang="Go"

  echo ">>> 开始 $lang 增量覆盖率分析 (基于 Block 语句权重)"

  # 提取变更行号：增加测试文件排除 + 非代码行过滤
  local changed_lines=$(cat /opt/cached_resources/pr_diffs/${platform}/${owner}/${repo}/${PR_ID}/pr_changes.patch | awk '
    # 1. 匹配 .go 文件，排除 _test.go
    /^\+\+\+ b\/.*\.go/ && !/.*_test\.go/ {
        # 提取 +++ b/path/to/file.go 中的完整路径（移除 b/ 前缀）
        file = substr($NF, 3);
        active=1;
        next
    }
    # 遇到新文件头则关闭采集
    /^\+\+\+ b\// {active=0; next}

    # 2. 仅在 active 状态下解析行号
    active == 1 && /^@@/ {
        # 提取当前代码块的起始行号
        split($3, a, ",");
        line = substr(a[1], 2);
        if (index(line, ",") > 0) line = substr(line, 1, index(line, ",")-1);
        next
    }

    # 3. 核心过滤逻辑：仅采集真正的“可执行”代码行
    active == 1 && /^\+/ && !/^\+\+\+/ {
        content = substr($0, 2); # 去掉开头的 +
        gsub(/[ \t\r\n]/, "", content); # 去掉所有空格和换行符以便匹配

        # 排除以下内容：
        # - 空行
        # - 单独的括号 ( } 或 { 或 ) )
        # - package 和 import 定义
        # - 单行注释 // 或 多行注释开始 /*
        if (content == "" ||
            content ~ /^}/ ||
            content ~ /^{/ ||
            content ~ /^\)/ ||
            content ~ /^package/ ||
            content ~ /^import/ ||
            content ~ /^\/\// ||
            content ~ /^\/\*/ ||
            content ~ /^type.*struct{/ ||
            content ~ /^type.*interface{/) {
                line++;
                next;
        }

        # 记录有效的可执行行
        print file ":" line;
        line++;
        next;
    }

    # 维持行号计数
    active == 1 && !/^\-/ { line++ }
')
  if [ -z "$changed_lines" ]; then
    echo ">>> [SKIP] $lang 在本次 PR 中没有代码变更。"
    SUMMARY_TABLE="${SUMMARY_TABLE}$(printf "| %-35s | %-12s | %-8s |\n" "${lang}(Inc)" "-" "SKIP")"
    return
  fi
  # 2. 增量解析逻辑：映射行到 Statement Block
  local total_inc_stmts=0
  local covered_inc_stmts=0
  local processed_blocks="" # 关键：用于存储 "文件:起始行,结束行" 实现去重
  for info in $changed_lines; do
    local target_file="${info%:*}"
    local target_line="${info#*:}"
    # 构造 grep 模式定位 coverage.out 中的对应文件
    # 假设 coverage.out 中的路径包含项目名，这里用 -F 做后缀匹配
    local pattern="${target_file}:"
    # 查找包含该行的 Block
    while read -r cover_line; do
      [[ -z "$cover_line" ]] && continue

      # 解析格式: path:start_line.col,end_line.col stmts hits
      # 示例: issue_pr_board/utils/aa.go:12.83,14.37 1 1
      local data="${cover_line##*:}"
      local line_range stmts hits
      read -r line_range stmts hits <<<"$data"

      local start_line="${line_range%%.*}"
      local tmp="${line_range#*,}"
      local end_line="${tmp%%.*}"
      # 逻辑：判定 PR 变更行是否在 Block 行号区间内
      if ((target_line >= start_line && target_line <= end_line)); then
        # 唯一标识该 Block 以便去重
        local block_uid="${target_file}:${line_range}"

        if [[ ! "$processed_blocks" =~ "$block_uid" ]]; then
          # 记录分母：相关 Block 的语句数
          total_inc_stmts=$((total_inc_stmts + stmts))
          echo "增量部分语句 ${block_uid} 语句数${stmts}，总语句数:${total_inc_stmts}"

          # 记录分子：如果该 Block 被命中过
          if ((hits > 0)); then
            covered_inc_stmts=$((covered_inc_stmts + stmts))
            echo "✅ 已覆盖增量部分语句 ${block_uid} 语句数${stmts}，增量语句数:${covered_inc_stmts}"
          else
            echo "❌ 未覆盖增量部分语句 ${block_uid} 语句数${stmts}，增量语句数:${covered_inc_stmts}"
          fi

          # 标记该 Block 已计算过，防止 PR 修改同一 Block 多行时重复累加
          processed_blocks="$processed_blocks $block_uid"
        fi
        # 一行代码只属于一个最小 Block，匹配到即可跳出当前文件的搜索
        break
      fi
    done < <(grep -F "$pattern" "$out_path")
  done

  # 3. 计算最终结果 (兼容无 bc 环境)
  echo "--------------------------------------"
  echo "增量分析报告 (基于 Statement Block):"

  if [ "$total_inc_stmts" -eq 0 ]; then
    INC_COV="100.0"
    status="PASS"
  else
    # 使用 awk 进行高精度计算并判定状态，通过 | 分隔返回结果
    # 这里一次性算出百分比和是否达标
    result=$(awk -v covered="$covered_inc_stmts" -v total="$total_inc_stmts" -v th="$threshold" 'BEGIN {
        rate = (covered / total) * 100;
        # 输出格式为：百分数|状态
        printf "%.1f|%s", rate, (rate < th ? "FAIL" : "PASS");
    }')

    INC_COV="${result%|*}"
    status="${result#*|}"
  fi

  echo "--------------------------------------"
  echo "增量分析报告 (基于 Statement Block):"
  echo "相关语句总数: $total_inc_stmts"
  echo "命中语句总数: $covered_inc_stmts"
  echo "增量覆盖率:   $INC_COV%"
  echo "--------------------------------------"

  # 5. 判定结果
  if [ "$status" = "FAIL" ]; then
    echo "❌ 判定结果: FAIL (阈值: $threshold%)"
    # 这里假设 lang 变量你之前已定义，如 lang="Go"
    SUMMARY_TABLE="${SUMMARY_TABLE}$(printf "| %-35s | %-12s | %-8s |\n" "${lang}(Inc)" "${INC_COV}%" "FAIL")"
    GLOBAL_SUCCESS=false
  else
    echo "✅ 判定结果: PASS"
    SUMMARY_TABLE="${SUMMARY_TABLE}$(printf "| %-35s | %-12s | %-8s |\n" "${lang}(Inc)" "${INC_COV}%" "PASS")"
  fi
}
# 参数: $1=语言, $2=XML报告路径
run_incremental_logic() {
  local lang=$1
  local xml_path=$2
  local threshold=80.0 # 增量覆盖率通常要求更高，例如 80%

  echo ">>> 开始 $lang 增量覆盖率分析 "

  # 定义语言后缀映射
  local ext=""
  case $lang in
  "Java") ext="java" ;;
  "Go") ext="go" ;;
  "Python") ext="py" ;;
  "Node.js") ext="(js|ts|jsx|tsx)" ;;
  *) ext=".*" ;;
  esac

  echo ">>> 开始 $lang 增量覆盖率分析"

  # 1. 检查物理 patch 文件是否存在
  if [ ! -f "/opt/cached_resources/pr_diffs/${platform}/${owner}/${repo}/${PR_ID}/pr_changes.patch" ]; then
    echo ">>> [SKIP] 未发现 pr_changes.patch。"
    return
  fi

  # 2. 核心优化：识别 patch 中是否有该语言的代码变更
  echo ">>> 正在检索 $lang 相关文件的变更..."

  # 改进点：
  # 1. 不再使用 ^ 强制行首，防止前面有空格干扰
  # 2. 使用 awk '{print $NF}' 强制抓取每一行的最后一个空格后的内容（即文件路径）
  # 3. 使用 sed 移除 a/ b/ 等前缀
# 改进后的匹配逻辑
# 1. grep "+++" 找到目标行
# 2. sed 去掉 "+++ " 前缀
# 3. sed 去掉可能存在的引号和 a/ b/ 前缀
local matched_files=$(grep "^+++" /opt/cached_resources/pr_diffs/${platform}/${owner}/${repo}/${PR_ID}/pr_changes.patch | \
                     sed 's/+++ //' | \
                     sed 's/^[ab]\///' | \
                     sed "s/'//g" | \
                     grep -E "\.${ext}$")

  if [ -z "$matched_files" ]; then
    echo ">>> [SKIP] $lang 在本次 PR 中没有代码变更。"
    SUMMARY_TABLE="${SUMMARY_TABLE}$(printf "| %-35s | %-12s | %-8s |\n" "${lang}(Inc)" "-" "SKIP")"
    return
  else
    # 打印具体找到了哪些文件，方便排查
    local count=$(echo "$matched_files" | wc -l)
    echo ">>> [MATCH] 发现 $count 个 $lang 变更文件。 $matched_files"
  fi
  cat  /opt/cached_resources/pr_diffs/${platform}/${owner}/${repo}/${PR_ID}/pr_changes.patch

  # 3. 执行分析
  # 增加 --src-roots . 确保路径匹配更准
  diff-cover --diff-file /opt/cached_resources/pr_diffs/${platform}/${owner}/${repo}/${PR_ID}/pr_changes.patch "$xml_path" --exclude "src/test/*" "**/Test*.java" "test_*.py" >diff_report.log 2>&1
  ls -l
  cat diff_report.log
  if grep -q "No lines with coverage information in this diff." diff_report.log; then
    echo "------------------------------------------------"
    echo "检测到无增量代码变更（或仅有非逻辑变更）。"
    echo "自动设置增量覆盖率为: 100%"
    echo "------------------------------------------------"
    INC_COV=100
  else
      # 提取数值 (保持原逻辑)
    # 1. 优先提取 "Coverage: " 后面的数值
    INC_COV=$(grep -i "Coverage:" diff_report.log | grep -oE '[0-9]+(\.[0-9]+)?' | head -n 1)

    # 2. 兜底提取：如果上面没拿到，尝试 Average coverage 格式
    if [ -z "$INC_COV" ]; then
      INC_COV=$(grep -i "Average coverage" diff_report.log | grep -oE '[0-9]+(\.[0-9]+)?' | head -n 1)
    fi
  fi


  # 3. 结果判定
  if [ -z "$INC_COV" ]; then
    echo "⚠️ $lang 匹配到了变更但未能解析出覆盖率数值，请检查 diff_report.log 内容。"
    SUMMARY_TABLE="${SUMMARY_TABLE}$(printf "| %-35s | %-12s | %-8s |\n" "${lang}(Inc)" "N/A" "ERROR")"
    GLOBAL_SUCCESS=false
  else
    echo ">>> $lang 增量覆盖率解析成功: ${INC_COV}%"
    # 此时 INC_COV 是 0，逻辑就对了

    # 门禁逻辑 (0 < 80，会触发报警)
    if awk "BEGIN { exit !($INC_COV < $threshold) }"; then
      echo "❌ $lang 增量覆盖率未达标 (${INC_COV}% < ${threshold}%)"
      SUMMARY_TABLE="${SUMMARY_TABLE}$(printf "| %-35s | %-12s | %-8s |\n" "${lang}(Inc)" "${INC_COV}%" "Fail")"
      GLOBAL_SUCCESS=false
    else
      echo " $lang 增量覆盖率已达标 (${INC_COV}% < ${threshold}%)"
      SUMMARY_TABLE="${SUMMARY_TABLE}$(printf "| %-35s | %-12s | %-8s |\n" "${lang}(Inc)" "${INC_COV}%" "Pass")"
    fi
  fi
}
# 辅助函数：将结果存入汇总表
# 参数: $1=语言, $2=覆盖率, $3=状态(PASS/FAIL/ERROR)
record_result() {
  local lang=$1
  local cov=$2
  local status=$3
  # 格式化存入变量，用于最后打印 (增加宽度以适应模块名)
  local line=$(printf "| %-35s | %-12s | %-8s |\n" "$lang" "$cov%" "$status")
  SUMMARY_TABLE="${SUMMARY_TABLE}${line} \n"

  if [ "$status" != "PASS" ]; then
    GLOBAL_SUCCESS=false
  fi
}

run_analysis_Go() {
  echo "开始分析目标目录:./"
  ls -l ./

  # 获取所有包含 go.mod 的目录
  echo ">>> 检测 Go 模块..."
  pwd=$(pwd)
  # 方法2: 查找所有 go.mod 文件的目录 (适用于根目录没有 go.mod 的情况)
  if [ -z "$MODULE_DIRS" ]; then
    echo ">>> 使用 find 查找 go.mod 文件..."
    GO_MOD_FILES=$(find . -name "go.mod" -type f 2>/dev/null)

    if [ -z "$GO_MOD_FILES" ]; then
      echo "❌ 未检测到 Go 模块"
      record_result "Go" "0.00" "PASS"
      return
    fi

    # 提取目录路径
    MODULE_DIRS=""
    for mod_file in $GO_MOD_FILES; do
      mod_dir=$(dirname "$mod_file")
      MODULE_DIRS="$MODULE_DIRS $mod_dir"
    done
  fi

  if [ -z "$MODULE_DIRS" ]; then
    echo "❌ 未检测到 Go 模块"
    record_result "Go" "0.00" "PASS"
    return
  fi

  echo ">>> 发现以下 Go 模块:"
  echo "$MODULE_DIRS"

  # 遍历每个模块
  for MODULE_DIR in $MODULE_DIRS; do
    echo ""
    echo "=========================================="
    echo ">>> 分析模块: $MODULE_DIR"
    echo "=========================================="

    # 为每个模块生成独立的覆盖率文件
    MODULE_NAME=$(basename "$MODULE_DIR")
    COVERAGE_FILE="./${MODULE_NAME}-${sanitizedBranch}-${TIMESTAMP}.out"

    # 运行测试
    echo ">>> 运行 go test..."
    cd $MODULE_DIR || continue
    ls -l
    go test -v ./... -coverprofile="$COVERAGE_FILE" 2>&1 | tail -20

    if [ ! -f "$COVERAGE_FILE" ]; then
      echo "⚠️ 模块 $MODULE_NAME 未生成覆盖率文件，跳过"
      record_result "Go($MODULE_NAME)" "$COVERAGE" "ERROR"
      continue
    fi

    # 计算全量覆盖率
    COVERAGE=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}' | sed 's/%//')
    echo ">>> 模块 $MODULE_NAME 全量覆盖率: $COVERAGE%"

    # 判定全量覆盖率 (阈值 10%)
    FULL_THRESHOLD=10.0
    if awk "BEGIN { exit !($COVERAGE < $FULL_THRESHOLD) }"; then
      echo "❌ 模块 $MODULE_NAME 全量覆盖率未达标 (${COVERAGE}% < ${FULL_THRESHOLD}%)"
      record_result "Go($MODULE_NAME)" "$COVERAGE" "FAIL"
    else
      echo "✅ 模块 $MODULE_NAME 全量覆盖率已达标"
      record_result "Go($MODULE_NAME)" "$COVERAGE" "PASS"
    fi

    # 增量覆盖率分析
    echo ">>> 开始增量覆盖率分析..."

    if [ -f "$COVERAGE_FILE" ]; then
      calculate_go_incremental_readable "$COVERAGE_FILE"
      INC_COV=${INC_COV:-0}

      INC_THRESHOLD=80.0
      if awk "BEGIN { exit !($INC_COV < $INC_THRESHOLD) }"; then
        echo "❌ 模块 $MODULE_NAME 增量覆盖率未达标 (${INC_COV}% < ${INC_THRESHOLD}%)"
      else
        echo "✅ 模块 $MODULE_NAME 增量覆盖率已达标"
      fi
    fi
    cd ${pwd}
  done
}

# 运行 SpotBugs 分析
run_analysis_Java() {
  echo "开始分析目标目录:./"
  ls -l ./
  $WORKDIR/apache-maven-3.9.6/bin/mvn clean test jacoco:report
  # 设置 Jacoco 报告 XML 文件的路径
  JACOCO_REPORT_XML="target/site/jacoco/jacoco.xml"
  # 使用 xmllint 提取所有具有 type="LINE" 的 counter 元素的 covered 和 missed 属性值
  total_covered=$(xmllint --xpath 'sum(//counter/@covered)' "$JACOCO_REPORT_XML" 2>/dev/null)
  total_missed=$(xmllint --xpath 'sum(//counter/@missed)' "$JACOCO_REPORT_XML" 2>/dev/null)

  # 检查是否成功提取了数据
  if [ -z "$total_covered" ] || [ -z "$total_missed" ]; then
    record_result "Java" "0.00" "ERROR"
    exit 1
  fi

  echo "Jacoco total_covered: $total_covered."
  echo "Jacoco total_missed: $total_missed."
  # 计算总体覆盖率
  total_lines=$((total_covered + total_missed))
  if [ $total_lines -eq 0 ]; then
    echo "Warning: No lines of code to cover."
    total_coverage_percent=0
  else
    total_coverage_percent=$(awk -v covered="$total_covered" -v lines="$total_lines" 'BEGIN {print (covered / lines) * 100}')
  fi

  # 输出总体覆盖率
  echo "Total line coverage percentage: $total_coverage_percent%"
  run_incremental_logic "Java" "$JACOCO_REPORT_XML"
  THRESHOLD=10.0
  if awk "BEGIN { exit !($total_coverage_percent < $THRESHOLD) }"; then
    echo "❌ Jacoco Coverge: $total_coverage_percent% is below 10%."
    record_result "Java" "$total_coverage_percent" "FAIL"
  else
    echo "✔️ Jacoco Coverge: $total_coverage_percent% is above 10%."
    record_result "✔️ Java" "$total_coverage_percent" "PASS"
  fi
}

run_analysis_Python() {
  if [ -f "requirements.txt" ]; then
    pip3 install -r requirements.txt --quiet
  elif [ -f "pyproject.toml" ]; then
    pip3 install . --quiet
  fi
  pip3 install pyyaml
  pip3 install requests
  pip3 install pynacl
  pip3 install pytest
  pip3 install pytest-cov
  python3 -m pytest --ignore-glob='test_*.py' --continue-on-collection-errors --cov=./ --cov-report=term-missing --cov-report=xml:coverage.xml >pytest_result.log 2>&1
  # 3. 结果验证逻辑
  if [ -f "coverage.xml" ]; then
      echo "✅ coverage.xml 生成成功！"
  else
      echo "❌ 错误: coverage.xml 未生成！"
      echo "🔍 Pytest 日志最后 10 行内容："
      tail -n 10 pytest_result.log
      # 可以在这里决定是否退出脚本，防止后续 diff-cover 崩溃
      # exit 1
  fi
  cat pytest_result.log
  current_coverage=$(grep "TOTAL" pytest_result.log | awk '{print $NF}' | tr -d '%')
  echo "Python Coverage Result: $current_coverage%"
  # 输出总体覆盖率
  echo "Total line coverage percentage: $current_coverage%"
  run_incremental_logic "Python" "coverage.xml"
  THRESHOLD=10.0
  if awk "BEGIN { exit !($current_coverage < $THRESHOLD) }"; then
    echo "❌ Jacoco Coverge: $current_coverage% is below 10%."
    record_result "Python" "$current_coverage" "FAIL"
  else
    echo "✔️ Jacoco Coverge: $current_coverage% is above 10%."
    record_result "✔️ Python" "$current_coverage" "PASS"
  fi
}

run_analysis_NodeJs() {
  local value
  value=$(jq -r '.packageManager' package.json)

  if [[ "$value" == *"yarn"* ]]; then
    echo "检测到 yarn，执行 yarn 操作"
    RESULT=99.99
    record_result "Node.js" "$RESULT" "PASS"
    return 0   # 跳出当前函数，返回调用处
  fi

  # 这部分代码只有在条件不满足时才会执行
  echo "执行默认操作"

  #如需安装node-sass
  #npm install node-sass --verbose
  #加载依赖
  pnpm config set side-effects-cache false # 避免缓存干扰


  # 1. 安装依赖时建议使用 --frozen-lockfile 保证环境一致性
  # 如果是 Monorepo，add 命令会自动处理，但建议在 root 执行

  echo "ls"
  ls

  # 自动生成 pnpm-workspace.yaml（如果需要）
  if [ -d "packages" ] && [ ! -f "pnpm-workspace.yaml" ]; then
    echo "检测到 packages 目录但缺少 pnpm-workspace.yaml，自动生成中..."
    echo "packages:" > pnpm-workspace.yaml
    echo "  - 'packages/*'" >> pnpm-workspace.yaml
    echo "✓ 已生成 pnpm-workspace.yaml"
  fi
  
  # 循环检查并创建缺失的工作区包
  while true; do
    echo "开始执行pnpm i"
    PNPM_OUTPUT=$(pnpm i 2>&1)
    PNPM_EXIT_CODE=$?

    # 检查是否有工作区包不存在的错误
    if echo "$PNPM_OUTPUT" | grep -q "ERR_PNPM_WORKSPACE_PKG_NOT_FOUND"; then
      # 提取缺失的包名（匹配 "package named \"xxx\" is present"）
      MISSING_PKG=$(echo "$PNPM_OUTPUT" | grep -oP 'no package named "\K[^"]+' | head -n 1)

      if [ -n "$MISSING_PKG" ]; then
        echo "检测到缺失的工作区包: $MISSING_PKG，自动创建中..."
        PKG_DIR="packages/$MISSING_PKG"

        if [ ! -d "$PKG_DIR" ]; then
          mkdir -p "$PKG_DIR"
          cat > "$PKG_DIR/package.json" <<EOF
{
  "name": "$MISSING_PKG",
  "private": true,
  "version": "1.0.0"
}
EOF
          echo "✓ 已创建包: $MISSING_PKG"
        fi
        # 继续循环重试
      else
        # 没有提取到包名，但有错误，直接退出
        echo "$PNPM_OUTPUT"
        exit $PNPM_EXIT_CODE
      fi
    elif [ $PNPM_EXIT_CODE -ne 0 ]; then
      # 其他错误，显示输出并退出
      echo "$PNPM_OUTPUT"
      exit $PNPM_EXIT_CODE
    else
      # 安装成功，显示输出并退出循环
      echo "$PNPM_OUTPUT"
      break
    fi
  done

  pnpm_add() {
    if [ -f "pnpm-workspace.yaml" ]; then
      echo "检测到 monorepo 环境，使用 -w 参数..."
      pnpm add -w "$@"
    else
      echo "普通项目，直接安装..."
      pnpm add "$@"
    fi
  }
  # 根据 package.json 中 vite 版本选择匹配的 vitest 版本
  VITE_VERSION=$(jq -r '(.devDependencies.vite // .dependencies.vite // "") | ltrimstr("^") | ltrimstr("~")' package.json 2>/dev/null || echo "")
  VITE_MAJOR=$(echo "$VITE_VERSION" | cut -d'.' -f1)
  if [ "$VITE_MAJOR" = "7" ] || [ "$VITE_MAJOR" -gt 7 ] 2>/dev/null; then
    VITEST_VERSION="4.0.18"
  elif [ "$VITE_MAJOR" = "6" ]; then
    VITEST_VERSION="3.0.7"
  else
    VITEST_VERSION="2.0.2"
  fi
  # 根据 vite 主版本确定 plugin-vue 和 plugin-vue-jsx 的推荐版本
  if [ "$VITE_MAJOR" = "7" ] || [ "$VITE_MAJOR" -gt 7 ] 2>/dev/null; then
    PLUGIN_VUE_VERSION="latest"
    PLUGIN_VUE_JSX_VERSION="latest"
  elif [ "$VITE_MAJOR" = "6" ]; then
    PLUGIN_VUE_VERSION="5"
    PLUGIN_VUE_JSX_VERSION="4"
  else
    PLUGIN_VUE_VERSION="4"
    PLUGIN_VUE_JSX_VERSION="3"
  fi

  echo "检测到 vite 版本: ${VITE_VERSION:-未找到}，使用 vitest@${VITEST_VERSION}"

  # 检查 @vitejs/plugin-vue 和 @vitejs/plugin-vue-jsx 是否已在项目依赖中声明
  HAS_PLUGIN_VUE=$(jq -r '(.devDependencies["@vitejs/plugin-vue"] // .dependencies["@vitejs/plugin-vue"] // "")' package.json 2>/dev/null || echo "")
  HAS_PLUGIN_VUE_JSX=$(jq -r '(.devDependencies["@vitejs/plugin-vue-jsx"] // .dependencies["@vitejs/plugin-vue-jsx"] // "")' package.json 2>/dev/null || echo "")

  EXTRA_PLUGINS=""
  if [ -z "$HAS_PLUGIN_VUE" ]; then
    echo "@vitejs/plugin-vue 未在项目依赖中，安装 @vitejs/plugin-vue@${PLUGIN_VUE_VERSION}"
    EXTRA_PLUGINS="$EXTRA_PLUGINS @vitejs/plugin-vue@${PLUGIN_VUE_VERSION}"
  else
    echo "@vitejs/plugin-vue 已存在（${HAS_PLUGIN_VUE}），跳过安装"
  fi
  if [ -z "$HAS_PLUGIN_VUE_JSX" ]; then
    echo "@vitejs/plugin-vue-jsx 未在项目依赖中，安装 @vitejs/plugin-vue-jsx@${PLUGIN_VUE_JSX_VERSION}"
    EXTRA_PLUGINS="$EXTRA_PLUGINS @vitejs/plugin-vue-jsx@${PLUGIN_VUE_JSX_VERSION}"
  else
    echo "@vitejs/plugin-vue-jsx 已存在（${HAS_PLUGIN_VUE_JSX}），跳过安装"
  fi

  echo "pnpm_add -D vitest@${VITEST_VERSION} @vitest/coverage-v8@${VITEST_VERSION}${EXTRA_PLUGINS}"
  pnpm_add -D vitest@${VITEST_VERSION} @vitest/coverage-v8@${VITEST_VERSION}${EXTRA_PLUGINS} jsdom@24.0.0 semver@7.5.1 --quiet

  # 检查本地配置文件是否存在
  LOCAL_CONFIG=""
  if [ -f "vitest.config.ts" ]; then
    LOCAL_CONFIG="vitest.config.ts"
  elif [ -f "vitest.config.js" ]; then
    LOCAL_CONFIG="vitest.config.js"
  elif [ -f "vitest.config.mjs" ]; then
    LOCAL_CONFIG="vitest.config.mjs"
  fi


  if [ -n "$LOCAL_CONFIG" ]; then
    echo "✓ 检测到本地配置: $LOCAL_CONFIG，基于它生成覆盖率配置（保留原配置文件）..."

    # 生成临时扩展配置，导入本地配置并合并 coverage 配置
    echo "
    import baseConfig from './$LOCAL_CONFIG';
    import { mergeConfig } from 'vitest/config';

    export default mergeConfig(baseConfig, {
      test: {
        coverage: {
          enabled: true,
          provider: 'v8',
          reporter: ['text', 'json', 'clover'],
          include: ['**/utils/**', '**/shared/utils.ts', '**/shared/utils.js'],
          exclude: ['node_modules/**', 'dist/**', 'public/**', '**.test.ts'],
        },
      },
    });
    " > vitest.config.coverage.js

    # 4. 执行测试 - 使用临时配置
    RESULT_FULL=$(npx vitest run --config=vitest.config.coverage.js --coverage --no-color 2>&1)
    ls -l coverage

    echo "$RESULT_FULL"

    # 清理临时配置文件
    rm vitest.config.coverage.js
  else
    echo "未检测到本地配置，生成默认配置..."

    # 2. 清理旧的默认配置文件（仅在没有本地配置时）
    [ -f "./vitest.config.js" ] && rm ./vitest.config.js

    # 3. 写入配置文件
    echo "
    import { defineConfig } from 'vitest/config';

    export default defineConfig({
      test: {
        environment: 'jsdom',
        globals: true,
        coverage: {
          enabled: true,
          provider: 'v8',
          reporter: ['text', 'json', 'clover'],
          include: ['**/utils/**', '**/shared/utils.ts', '**/shared/utils.js'],
          exclude: ['node_modules/**', 'dist/**', 'public/**', '**.test.ts'],
        },
      },
    });
    " > vitest.config.js

    # 4. 执行测试
    RESULT_FULL=$(npx vitest run --coverage --no-color 2>&1)
    ls -l coverage

    echo "$RESULT_FULL"
  fi


  # 5. 增强型解析（处理表格解析不到的情况）
  # 先尝试从控制台文本提取
  RESULT=$(echo "$RESULT_FULL" | grep 'All files' | awk -F'|' '{print $2}' | grep -oE '[0-9.]+' | head -n 1)

  echo "Parsed Coverage Result: $RESULT%"
  run_incremental_logic "Node.js" "coverage/clover.xml"
  THRESHOLD=10.0
  if [ -z "$RESULT" ]; then
    record_result "Node.js" "0.00" "ERROR"
  elif awk "BEGIN { exit !($RESULT < $THRESHOLD) }"; then
    record_result "Node.js" "$RESULT" "FAIL"
  else
    record_result "Node.js" "$RESULT" "PASS"
  fi
}

# 主逻辑
main() {
  remote_url=$(git remote -v | grep '(fetch)' | awk '{print $2}')
  [ -z "$remote_url" ] && {
    echo "❌ 无法获取 git 仓库地址"
    exit 1
  }
  url_no_git=${remote_url%.git}
  ORG_NAME=$(echo "$url_no_git" | awk -F'[/:]' '{print $(NF-1)}')
  PROJECT_PATH=$(echo "$url_no_git" | awk -F'[/:]' '{print $(NF-1)"/"$NF}')
  PLATFORM=$(echo "$url_no_git" | awk -F'[/:]' '{print $4}')
  # 去除 .git、ssh@、https://、http://
  cleaned_repo_url=$(echo "$remote_url" | sed -E 's#https?://##' | tr '/' '@')
 
  # 将清理后的结果赋值给变量
  echo "Cleaned URL: $cleaned_repo_url"
  # 分析目标代码（将 ./src 替换为你实际的代码目录）
  TARGET_DIR="."
  # 1. 定义支持的语言
  SUPPORTED_LANGS="Java|Python|Go|JavaScript|TypeScript"

  # 2. 获取 TOP 3 语言（存入一个换行符分隔的字符串）
  TOP_LANGS_STR=$(scc --format csv . |
    grep -E "$SUPPORTED_LANGS" |
    sort -t, -k6 -nr |
    head -n 3 |
    cut -d, -f1)

  # 3. 验证并循环处理
  if [ -z "$TOP_LANGS_STR" ]; then
    echo "未检测到支持的语言。"
  else
    echo "检测到以下 TOP 语言：$TOP_LANGS_STR"
    tsjsexecuted=0
    # 直接在 for 循环中使用字符串，Shell 会自动按换行/空格切分
    for lang in $TOP_LANGS_STR; do
      echo ">>> 正在处理: $lang"

      # 后面可以接你的逻辑
      case $lang in
      "Java")
        run_analysis_Java
        ;;
      "Python")
        run_analysis_Python
        ;;
      "Go")
        run_analysis_Go
        ;;
      "JavaScript" | "TypeScript")
        if ((tsjsexecuted == 0)); then
          if [ -f "package.json" ]; then
            run_analysis_NodeJs
            tsjsexecuted=1
          fi
        fi
        ;;
      esac
    done
  fi
  # --- 打印最终报告 ---
  echo ""
  echo "==========================================================="
  echo "              PROJECT COVERAGE FINAL REPORT                "
  echo "==========================================================="
  printf "| %-35s | %-12s | %-8s |\n" "TYPE" "COVERAGE" "RESULT"
  echo "-----------------------------------------------------------"
  echo "$SUMMARY_TABLE"
  echo "==========================================================="
  rm -rf /opt/cached_resources/pr_diffs/${platform}/${owner}/${repo}/${PR_ID}
  # --- 返回值判定 ---
  if [ "$GLOBAL_SUCCESS" = true ]; then
    echo ">>> ✔️ [SUCCESS] 所有语言覆盖率均达标。"
    exit 0
  else
    echo ">>>  ❌ [FAILURE] 存在语言未达到覆盖率阈值或测试出错！"
    exit 1 # 抛出非零值给 CI 系统
  fi
}

main
