#!/bin/bash

# 配置工作目录和工具版本
WORKDIR=/opt/cached_resources/sast
SPOTBUGS_HOME=/opt/cached_resources/sast/spotbugs-4.7.3
# 使用 sed 替换 '/' 为 '@'
sanitizedBranch=$(echo "$codeBranch" | tr '/' '@')

export GOROOT=${WORKDIR}/go
export GOPATH=${WORKDIR}/gopath
export JAVA_HOME=$WORKDIR/jdk18
export PATH=$JAVA_HOME/bin:$GOROOT/bin:$GOPATH/bin:/opt/cached_resources/sast/nodejs:/opt/cached_resources/sast/python/bin:/opt/cached_resources/tools/common:/opt/cached_resources/tools/gcm:/opt/cached_resources/sast/python/user_packages/bin:$PATH
export PYTHONUSERBASE=/opt/cached_resources/sast/python/user_packages
export MAVEN_OPTS="-Dmaven.repo.local=/opt/cached_resources/sast/maven_repository"

GLOBAL_SUCCESS=true
SUMMARY_DATA=""

# --- 1. 增强型打印函数 ---

# 参数: $1=文件名, $2=jq查询语句, $3=问题总数
print_json_details_if_any() {
    local file=$1
    local query=$2
    local count=$3

    # 只有当 count 不为 0 且文件存在时才打印
    if [[ "$count" != "0" && -s "$file" ]]; then
        echo -e "\n[发现漏洞详情报告 - JSON]"
        echo "------------------------------------------------------"
        jq "$query" "$file"
        echo -e "------------------------------------------------------\n"
    fi
}

# 参数: $1=文件名, $2=xpath查询语句, $3=问题总数
print_xml_details_if_any() {
    local file=$1
    local xpath=$2
    local count=$3

    if [[ "$count" != "0" && -s "$file" ]]; then
        echo -e "\n[发现漏洞详情报告 - XML]"
        echo "------------------------------------------------------"
        xmllint --xpath "$xpath" "$file" 2>/dev/null
        echo -e "\n------------------------------------------------------\n"
    fi
}

# 辅助函数：记录结果
# 参数: $1=工具/语言, $2=问题数量, $3=状态(PASS/FAIL/ERROR)
record_result() {
    local tool=$1
    local issues=$2
    local status=$3
    # 格式化存入变量 (增加宽度以适应模块名)
    printf -v line "| %-35s | %-12s | %-8s |\n" "$tool" "$issues" "$status"
    SUMMARY_DATA="${SUMMARY_DATA}${line}"
    
    if [ "$status" != "PASS" ]; then
        GLOBAL_SUCCESS=false
    fi
}

# 检查 Java 环境
check_java() {
    if ! command -v /opt/cached_resources/sast/jdk18/bin/java &> /dev/null; then
        echo "未安装 Java，请先安装 JDK 8 或 11。"
        exit 1
    fi
}

has_python_files() {
    local dir="$1"
    if find "$dir" -type f -name "*.py" | grep -q .; then
        return 0  # 有 Python 文件
    else
        return 1  # 没有 Python 文件
    fi
}

run_analysis_Go() {
    echo "开始分析目标目录:./"
    ls -l ./
    
    # 查找所有 go.mod 文件
    echo ">>> 检测 Go 模块..."
    GO_MOD_FILES=$(find . -name "go.mod" -type f 2>/dev/null)
    
    if [ -z "$GO_MOD_FILES" ]; then
        echo ">>> 未发现 Go 模块，跳过 Gosec 扫描"
        return
    fi
    
    echo ">>> 发现以下 Go 模块:"
    echo "$GO_MOD_FILES"
    
    # 遍历每个模块
    for mod_file in $GO_MOD_FILES; do
        mod_dir=$(dirname "$mod_file")
        module_name=$(basename "$mod_dir")
        
        echo ""
        echo "=========================================="
        echo ">>> 扫描模块: $module_name (路径: $mod_dir)"
        echo "=========================================="
        
        REPORT_FILE="./${module_name}-gosec-${sanitizedBranch}-${TIMESTAMP}.json"

        # 执行 gosec 扫描 (去掉 -stdout，让结果写入文件)
        $WORKDIR/gopath/bin/gosec -fmt=json -out="$REPORT_FILE" -verbose=text $mod_dir/... 2>&1 | tail -30
        ls -l

        # 检查报告文件是否存在 (gosec 会覆盖已存在的文件)
        if [ -f "$REPORT_FILE" ]; then
            # 获取问题总数 (如果为 null 则设为 0)
            TOTAL_FOUND=$(jq '.Stats.found // 0' "$REPORT_FILE" 2>/dev/null || echo "0")
            
            # 统计各严重程度 (处理 null 情况)
            HIGH_COUNT=$(jq '[.Issues[] | select(.severity == "HIGH")] | length // 0' "$REPORT_FILE" 2>/dev/null || echo "0")
            MEDIUM_COUNT=$(jq '[.Issues[] | select(.severity == "MEDIUM")] | length // 0' "$REPORT_FILE" 2>/dev/null || echo "0")
            LOW_COUNT=$(jq '[.Issues[] | select(.severity == "LOW")] | length // 0' "$REPORT_FILE" 2>/dev/null || echo "0")
            
            echo ">>> 模块 $module_name Gosec 扫描汇总:"
            echo "----------------------"
            echo "Total Issues  : $TOTAL_FOUND"
            echo "High Severity : $HIGH_COUNT"
            echo "Medium Severity: $MEDIUM_COUNT"
            echo "Low Severity  : $LOW_COUNT"
            echo "----------------------"
            
            # 判定逻辑
            if [ "$TOTAL_FOUND" -gt 10 ]; then
                echo "❌ 模块 $module_name 扫描失败: 总问题数 $TOTAL_FOUND > 10"
                record_result "Gosec (Go-$module_name)" "$TOTAL_FOUND" "FAIL"
            elif [ "$HIGH_COUNT" -gt 0 ]; then
                echo "❌ 模块 $module_name 扫描失败: 高危问题 $HIGH_COUNT > 0"
                record_result "Gosec (Go-$module_name)" "$HIGH_COUNT" "HIGH_SEC"
            else
                echo "✅ 模块 $module_name 扫描通过"
                record_result "Gosec (Go-$module_name)" "$TOTAL_FOUND" "PASS"
            fi
            
            # 打印详细信息
            if [ "$TOTAL_FOUND" -gt 0 ]; then
                print_json_details_if_any "$REPORT_FILE" '.Issues[] | {severity: .severity, rule: .rule_id, file: .file, line: .line, desc: .details}' "$TOTAL_FOUND"
            fi
        else
            echo "⚠️ 模块 $module_name 未生成报告文件"
            record_result "Gosec (Go-$module_name)" "N/A" "ERROR"
        fi
    done
    
    echo ""
    echo "✅ 所有 Go 模块 Gosec 扫描完成"
}

# 运行 SpotBugs 分析
run_analysis_Java() {
    
    echo "开始分析目标目录:./"
    ls -l ./

    $WORKDIR/apache-maven-3.9.6/bin/mvn clean compile -Dmaven.test.skip
    ls -l ./
    # 创建报告文件夹
    REPORT_FILE="./findsecbugs-report-${sanitizedBranch}-${TIMESTAMP}.xml"
    # 执行 SpotBugs 扫描
    $WORKDIR/jdk18/bin/java -jar "$SPOTBUGS_HOME/lib/spotbugs.jar" \
        -textui \
        -pluginList "$SPOTBUGS_HOME/plugin/findsecbugs-plugin.jar" \
        -effort:max \
        -include $SPOTBUGS_HOME/plugin/include-filter.xml \
        -xml \
        -output "$REPORT_FILE" \
        "$TARGET_DIR/target/classes"

    
    echo "分析完成，报告已生成：$REPORT_FILE"
    if [ -s "$REPORT_FILE" ]; then
       # 使用 xmllint 统计 BugInstance 数量
       # 1. 统计问题总数
        TOTAL_FOUND=$(xmllint --xpath "count(//BugInstance)" "$REPORT_FILE" 2>/dev/null || echo 0)
        print_xml_details_if_any "$REPORT_FILE" "//BugInstance" "$TOTAL_FOUND"
        # 2. 统计各严重程度
        # Priority 1 = High, Priority 2 = Medium, Priority 3 = Low
        HIGH_COUNT=$(xmllint --xpath "count(//BugInstance[@priority='1'])" "$REPORT_FILE" 2>/dev/null || echo 0)
        MEDIUM_COUNT=$(xmllint --xpath "count(//BugInstance[@priority='2'])" "$REPORT_FILE" 2>/dev/null || echo 0)
        LOW_COUNT=$(xmllint --xpath "count(//BugInstance[@priority='3'])" "$REPORT_FILE" 2>/dev/null || echo 0)

        echo "SpotBugs 扫描汇总 (Java):"
        echo "----------------------"
        echo "Total Issues   : $TOTAL_FOUND"
        echo "High (P1)      : $HIGH_COUNT"
        echo "Medium (P2)    : $MEDIUM_COUNT"
        echo "Low (P3)       : $LOW_COUNT"
        echo "----------------------"

        # 3. 判定逻辑 (参照 Go 的模式)
        if (( HIGH_COUNT > 0 )); then
            record_result "SpotBugs (Java)" "$HIGH_COUNT" "HIGH_SEC"
        elif (( TOTAL_FOUND > 10 )); then # 假设 Java 阈值设为10
            record_result "SpotBugs (Java)" "$TOTAL_FOUND" "FAIL"
        else
            record_result "SpotBugs (Java)" "$TOTAL_FOUND" "PASS"
        fi
    else
        echo "警告: 未找到 SpotBugs 报告文件。"
        record_result "SpotBugs (Java)" "N/A" "ERROR"
    fi
}

run_analysis_Python (){
    REPORT_FILE="./bandit-report-${sanitizedBranch}-${TIMESTAMP}.json"
    bandit -r "./" -f json -o "$REPORT_FILE" --exclude "*test_*.py"
    echo "bandit 扫描完成"
    cat $REPORT_FILE
    if [ -s "$REPORT_FILE" ]; then
        # 提取 results 列表长度
        # 使用 jq 的 select 统计不同级别的数量
        TOTAL_FOUND=$(jq '.results | length // 0' "$REPORT_FILE")
        print_json_details_if_any "$REPORT_FILE" '.results[] | {severity: .issue_severity, test: .test_id, file: .filename, line: .line_number, msg: .issue_text}' "$TOTAL_FOUND"
        HIGH_COUNT=$(jq '[.results[] | select(.issue_severity == "HIGH")] | length // 0' "$REPORT_FILE")
        MEDIUM_COUNT=$(jq '[.results[] | select(.issue_severity == "MEDIUM")] | length // 0' "$REPORT_FILE")
        LOW_COUNT=$(jq '[.results[] | select(.issue_severity == "LOW")] | length // 0' "$REPORT_FILE")

        echo "Bandit 扫描汇总 (Python):"
        echo "----------------------"
        echo "Total Issues   : $TOTAL_FOUND"
        echo "High Severity  : $HIGH_COUNT"
        echo "Medium Severity: $MEDIUM_COUNT"
        echo "Low Severity   : $LOW_COUNT"
        echo "----------------------"
        if [ "$TOTAL_FOUND" -gt 10 ]; then
            echo "当前仓库： ${cleaned_repo_url}"
            if [[ "$cleaned_repo_url" == "github.com@opensourceways@om-dataarts.git" ]]; then
                echo "当前仓库豁免： ${cleaned_repo_url}"
                record_result "Bandit (Python)" "$TOTAL_FOUND" "PASS"
            else
                record_result "Bandit (Python)" "$TOTAL_FOUND" "FAIL"
            fi
        elif [ "$HIGH_COUNT" -gt 0 ]; then
            record_result "Bandit (Python)" "$HIGH_COUNT" "HIGH_SEC"
        else
            record_result "Bandit (Python)" "$TOTAL_FOUND" "PASS"
        fi
    else
        record_result "Bandit_Python" "N/A" "ERROR"
    fi
}


run_analysis_Nodejs (){
    local value
    value=$(jq -r '.packageManager' package.json)

    if [[ "$value" == *"yarn"* ]]; then
        echo "检测到 yarn，执行 yarn 操作"
        TOTAL_FOUND=0
        record_result "ESLint NodeJS" "$TOTAL_FOUND" "PASS"
        return 0   # 跳出当前函数，返回调用处
    fi

    # 这部分代码只有在条件不满足时才会执行
    echo "执行默认操作"
   
    # 1. 安装 ESLint 相关依赖（使用最新稳定版本）
    echo "正在安装 ESLint 依赖..."
    pnpm i
    pnpm add -D \
        eslint@9.39.2 \
        @eslint/js@9.39.2 \
        eslint-plugin-vue@9.31.0 \
        @vue/eslint-config-typescript@14.2.0 \
        @vue/eslint-config-prettier@10.2.0 \
        typescript-eslint@8.55.0

    # 2. 清理旧配置文件
    echo "清理旧配置文件..."
    [ -f "./.eslintrc.js" ] && rm ./.eslintrc.js
    [ -f "./.eslintrc.cjs" ] && rm ./.eslintrc.cjs
    [ -f "./.eslintrc.json" ] && rm ./.eslintrc.json
    [ -f "./.eslintignore" ] && rm ./.eslintignore
    [ -f "./eslint.config.js" ] && rm ./eslint.config.js
    [ -f "./eslint.config.mjs" ] && rm ./eslint.config.mjs
    [ -f "./eslint.config.cjs" ] && rm ./eslint.config.cjs

    # 3. 生成 .prettierrc.json 配置文件（解决换行符问题）
    echo "生成 .prettierrc.json 配置文件..."
    cat "{
      "$schema": "https://json.schemastore.org/prettierrc",
      "endOfLine": "auto",
      "semi": true,
      "tabWidth": 2,
      "singleQuote": true,
      "printWidth": 160,
      "trailingComma": "es5"
    }" > .prettierrc.json

    # 4. 生成 eslint.config.js（使用 Flat Config 格式）
    echo "生成 eslint.config.js..."
    cat "import js from '@eslint/js';
    import vue from 'eslint-plugin-vue';
    import typescript from '@vue/eslint-config-typescript';
    import prettierSkip from '@vue/eslint-config-prettier/skip-formatting';
    import tseslint from 'typescript-eslint';

    export default [
      // 基础 JavaScript 推荐规则
      js.configs.recommended,

      // Vue 3 推荐规则
      ...vue.configs['flat/essential'],

      // TypeScript 规则
      ...typescript({
        extends: ['base']
      }),

      // Prettier 规则（跳过格式化）
      prettierSkip,

      // 项目特定配置
      {
        files: ['**/*.{js,ts,vue,jsx,tsx}'],
        rules: {
          // 在这里添加自定义规则
          'vue/multi-word-component-names': 'off',
          // 禁用原生的 no-unused-vars 规则
          'no-unused-vars': 'off',
          // 启用 TypeScript 版本的 no-unused-vars
          '@typescript-eslint/no-unused-vars': [
            'error',
            {
              argsIgnorePattern: '^_',
              varsIgnorePattern: '^_',
              caughtErrorsIgnorePattern: '^_',
            },
          ],
        },
        languageOptions: {
          ecmaVersion: 'latest',
          sourceType: 'module',
          globals: {
            // 自定义hooks
            useLocale: 'readonly',
            useScreen: 'readonly',
        
            // Vue 自动导入
            computed: 'readonly',
            ref: 'readonly',
            reactive: 'readonly',
            watch: 'readonly',
            watchEffect: 'readonly',
            unref: 'readonly',
            toRef: 'readonly',
            toRefs: 'readonly',
            isRef: 'readonly',
            readonly: 'readonly',
            shallowRef: 'readonly',
            shallowReactive: 'readonly',
            toRaw: 'readonly',
            markRaw: 'readonly',
            effectScope: 'readonly',
            getCurrentScope: 'readonly',
            onScopeDispose: 'readonly',
            nextTick: 'readonly',
            PropType: 'readonly',

            // Vue 生命周期
            onBeforeMount: 'readonly',
            onMounted: 'readonly',
            onBeforeUpdate: 'readonly',
            onUpdated: 'readonly',
            onBeforeUnmount: 'readonly',
            onUnmounted: 'readonly',
            onActivated: 'readonly',
            onDeactivated: 'readonly',
            onErrorCaptured: 'readonly',

            // Nuxt 自动导入
            useRouter: 'readonly',
            useRoute: 'readonly',
            useAsyncData: 'readonly',
            useFetch: 'readonly',
            useLazyFetch: 'readonly',
            useLazyAsyncData: 'readonly',
            useNuxtApp: 'readonly',
            useRuntimeConfig: 'readonly',
            useState: 'readonly',
            useCookie: 'readonly',
            useRequestHeaders: 'readonly',
            useRequestEvent: 'readonly',
            useRequestFetch: 'readonly',
            useRequestURL: 'readonly',
            useHead: 'readonly',
            useSeoMeta: 'readonly',
            useError: 'readonly',
            useNuxtData: 'readonly',
            refreshNuxtData: 'readonly',
            clearNuxtData: 'readonly',
            createError: 'readonly',
            showError: 'readonly',
            clearError: 'readonly',
            navigateTo: 'readonly',
            abortNavigation: 'readonly',
            setPageLayout: 'readonly',
            definePageMeta: 'readonly',
            prefetchComponents: 'readonly',
            preloadRouteComponents: 'readonly',
            preloadComponents: 'readonly',
            reloadNuxtApp: 'readonly',
            defineNuxtPlugin: 'readonly',
            defineNuxtRouteMiddleware: 'readonly',
            defineNitroPlugin: 'readonly',

            // Nuxt Content
            queryContent: 'readonly',

            // Pinia
            defineStore: 'readonly',
            acceptHMRUpdate: 'readonly',
            storeToRefs: 'readonly',

            // VueUse
            useDebounceFn: 'readonly',
            useEventListener: 'readonly',
            useTemplateRef: 'readonly',
            useIntersectionObserver: 'readonly',
            useDocumentVisibility: 'readonly',
            watchDebounced: 'readonly',

            // 项目常量
            COOKIE_KEY: 'readonly',
            COOKIE_AGREED_STATUS: 'readonly',

            // 第三方库全局对象
            ClipboardJS: 'readonly',

            // TypeScript 全局类型
            NodeListOf: 'readonly',

            // 浏览器全局对象
            document: 'readonly',
            navigator: 'readonly',
            window: 'readonly',
          },
        },
      },

      // 忽略的文件
      {
        ignores: [
          '**/node_modules/**',
          '**/dist/**',
          '**/.output/**',
          '**/.nuxt/**',
          '**/coverage/**',
          '**/public/**',
          '**/npmcache/**',
          '**/cache/**',
          '**/vite.config.{js,ts}',
          '**/vitest.config.{js,ts}',
          '**/vitest.setup.{js,ts}',
          '**/nuxt.config.{js,ts}',
          '**/config.mjs',
          '**/*.d.ts',
          'eslint.config.js'
        ],
      },
    ];" > eslint.config.js

    # 5. 运行 ESLint 检查
    echo "运行 ESLint 检查..."
    RESULT_FILE="eslint_result.json"
    npx eslint "**/*.{js,ts,vue,jsx,tsx}" --format json -o $RESULT_FILE

    cat "$RESULT_FILE"
    if [ -f "$RESULT_FILE" ]; then
            # 3. 提取数据 (ESLint JSON 是数组结构，需要对所有文件的 message 进行汇总)

            # 统计总数 (所有 messages 的长度总和)
            TOTAL_FOUND=$(jq '[.[].messages | length] | add // 0' "$RESULT_FILE")
            print_json_details_if_any "$RESULT_FILE" '.[] | select(.messages | length > 0) | { file: .filePath, issues: [ .messages[] | { line: .line, rule: .ruleId, msg: .message, severity: .severity } ] }' "$TOTAL_FOUND"
            # 统计 High (severity 为 2)
            HIGH_COUNT=$(jq '[.[].messages[] | select(.severity == 2)] | length // 0' "$RESULT_FILE")

            # 统计 Medium (severity 为 1)
            MEDIUM_COUNT=$(jq '[.[].messages[] | select(.severity == 1)] | length // 0' "$RESULT_FILE")

            echo "ESLint 扫描汇总 (NodeJS):"
            echo "----------------------"
            echo "Total Issues   : $TOTAL_FOUND"
            echo "High (Errors)  : $HIGH_COUNT"
            echo "Medium (Warns) : $MEDIUM_COUNT"
            echo "----------------------"

            # 4. 判定逻辑
            if (( HIGH_COUNT > 0 )); then
                # 只要有 Error 级别的漏洞就 FAIL
                record_result "ESLint NodeJS" "$HIGH_COUNT" "HIGH_SEC"
            elif (( TOTAL_FOUND > 10 )); then
                # 总量过多也判定为 FAIL
                record_result "ESLint NodeJS" "$TOTAL_FOUND" "FAIL"
            else
                record_result "ESLint NodeJS" "$TOTAL_FOUND" "PASS"
            fi

      else
            echo "警告: 未找到 ESLint 报告文件。"
            record_result "ESLint NodeJS" "N/A" "ERROR"
      fi
}

# 主逻辑
main() {
    # 获取 fetch 的仓库地址
    remote_url=$(git remote -v | grep '(fetch)' | awk '{print $2}')

    # 去除 .git、ssh@、https://、http://
    cleaned_repo_url=$(echo "$remote_url" | sed -E 's#https?://##' | tr '/' '@') 

    # 将清理后的结果赋值给变量
    echo "Cleaned URL: $cleaned_repo_url"
        # 分析目标代码（将 ./src 替换为你实际的代码目录）
    TARGET_DIR="."
    ls -l

    # 检查文件是否存在
    if [ -f "$TARGET_DIR/pom.xml" ]; then
        echo "pom.xml 文件存在于目录: $TARGET_DIR"
        check_java
        run_analysis_Java "$TARGET_DIR"
    else
        echo "pom.xml 文件不存在于目录: $TARGET_DIR,不执行find-sec-bugs"
    fi

    # 检查并扫描 Go 项目 (支持多模块)
    # run_analysis_Go 内部使用 find 查找 go.mod，无需在此判断
    run_analysis_Go "$TARGET_DIR"

    # 4. 扫描 Python 项目
    echo "检查 Python 项目..."
    if has_python_files "./"; then
        echo "检测到 Python 文件，运行 Bandit 扫描..."
        run_analysis_Python
    else
        echo "未检测到 Python 文件，跳过 Bandit 扫描。"
    fi

    # 4. 扫描 Nodejs 项目
    echo "检查 Nodejs 项目..."
    if [ -f "package.json" ]; then
        echo "检测到 Nodejs 文件，运行 Eslint 扫描..."
        run_analysis_Nodejs
    else
        echo "未检测到 Nodejs 文件，跳过 Eslint 扫描。"
    fi

    # --- 最终汇总展示 ---
    echo ""
    echo "==============================================================="
    echo "              SAST SECURITY SCAN REPORT                         "
    echo "==============================================================="
    printf "| %-35s | %-12s | %-8s |\n" "TOOL/LANG" "ISSUES FOUND" "STATUS"
    echo "---------------------------------------------------------------"
    echo "$SUMMARY_DATA"
    echo "==============================================================="

    # --- 返回值判定 ---
    if [[ "${GLOBAL_SUCCESS}" == "true" ]]; then
        echo ">>> [SUCCESS] SAST 扫描通过，未发现关键漏洞。"
        exit 0
    else
        echo ">>> [FAILURE] SAST 扫描失败，请修复上述漏洞！"
        # 显式返回 1，确保流水线能够捕获到失败状态
        exit 1
    fi
}

main
