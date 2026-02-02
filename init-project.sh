#!/bin/bash

# 项目初始化脚本 - 自动替换项目名称
# 用法: ./init-project.sh [new-project-name]
# 如果不提供参数，则使用当前文件夹名称作为项目名

set -e

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 原始项目名称（模板名称）
OLD_MODULE_NAME="go-zero-template"      # Go 模块名称
OLD_SERVICE_NAME="go_zero_template-api" # API 服务名称（下划线分隔）
OLD_SERVICE_FILE="goZeroTemplate-Api"   # 服务文件名（驼峰命名）
OLD_API_FILE="go_zero_template"         # API 文件名（下划线分隔，不含 -api）

# 确定新项目名称
if [ -z "$1" ]; then
    # 如果没有提供参数，使用当前文件夹名称
    NEW_MODULE_NAME=$(basename "$PWD")
    echo "未提供项目名称，使用当前文件夹名称: $NEW_MODULE_NAME"
else
    NEW_MODULE_NAME="$1"
    echo "使用指定的项目名称: $NEW_MODULE_NAME"
fi

# 从模块名生成各种命名格式
# 例如输入: my-new-service
NEW_SERVICE_NAME="${NEW_MODULE_NAME//-/_}-api"           # my_new_service-api（API 服务名）
NEW_API_FILE="${NEW_MODULE_NAME//-/_}"                   # my_new_service（API 文件名）

# 生成驼峰命名的服务文件名（例如: my-new-service -> myNewService-Api）
# 1. 将连字符分隔的单词转为驼峰格式
# 2. 第一个单词首字母小写，其他单词首字母大写
IFS='-' read -ra PARTS <<< "$NEW_MODULE_NAME"
NEW_SERVICE_FILE=""
for i in "${!PARTS[@]}"; do
    PART="${PARTS[$i]}"
    if [ $i -eq 0 ]; then
        # 第一个单词：首字母小写
        NEW_SERVICE_FILE="${PART}"
    else
        # 其他单词：首字母大写
        PART_CAPITALIZED="$(tr '[:lower:]' '[:upper:]' <<< ${PART:0:1})${PART:1}"
        NEW_SERVICE_FILE="${NEW_SERVICE_FILE}${PART_CAPITALIZED}"
    fi
done
NEW_SERVICE_FILE="${NEW_SERVICE_FILE}-Api"              # myNewService-Api

# 检查新旧名称是否相同
if [ "$OLD_MODULE_NAME" = "$NEW_MODULE_NAME" ]; then
    echo "项目名称与模板名称相同，无需替换"
    exit 0
fi

echo "开始替换项目名称:"
echo "  模块名: $OLD_MODULE_NAME -> $NEW_MODULE_NAME"
echo "  服务名: $OLD_SERVICE_NAME -> $NEW_SERVICE_NAME"
echo "  服务文件: $OLD_SERVICE_FILE -> $NEW_SERVICE_FILE"
echo "  API 文件: $OLD_API_FILE -> $NEW_API_FILE"
echo "-------------------------------------------"

# 需要替换 Go 模块名称的文件列表
MODULE_NAME_FILES=(
    "go.mod"
    "${OLD_SERVICE_FILE}.go"
    "internal/middleware/authMiddleware.go"
    "internal/svc/redis.go"
    "internal/svc/serviceContext.go"
    "internal/svc/db.go"
    "internal/handler/routes.go"
    "internal/logic/ping/pingUserServiceLogic.go"
    "internal/handler/system/healthHandler.go"
    "internal/handler/ping/pingUserServiceHandler.go"
    "internal/logic/system/healthLogic.go"
    "internal/request/user.go"
    "internal/request/request.go"
)

# 需要替换 API 服务名称的文件列表
SERVICE_NAME_FILES=(
    "${OLD_API_FILE}.api"
    "etc/${OLD_SERVICE_FILE}.yaml"
)

# 需要更新文件名引用的文件
FILE_REFERENCE_FILES=(
    "makefile"
)

# 统计替换次数
TOTAL_REPLACEMENTS=0

# 函数: 替换文件中的内容
replace_in_file() {
    local file=$1
    local old_pattern=$2
    local new_pattern=$3
    
    if [ -f "$file" ]; then
        # 使用 sed 进行替换
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS 使用 BSD sed
            sed -i '' "s|$old_pattern|$new_pattern|g" "$file"
        else
            # Linux 使用 GNU sed
            sed -i "s|$old_pattern|$new_pattern|g" "$file"
        fi
        
        # 检查文件是否包含新名称（验证替换成功）
        if grep -q "$new_pattern" "$file"; then
            echo "✓ $file ($old_pattern -> $new_pattern)"
            return 0
        fi
    else
        echo "⚠ 文件不存在: $file"
        return 1
    fi
}

# 替换 Go 模块名称
echo ""
echo "替换 Go 模块名称..."
for file in "${MODULE_NAME_FILES[@]}"; do
    if replace_in_file "$file" "$OLD_MODULE_NAME" "$NEW_MODULE_NAME"; then
        TOTAL_REPLACEMENTS=$((TOTAL_REPLACEMENTS + 1))
    fi
done

# 替换 API 服务名称
echo ""
echo "替换 API 服务名称..."
for file in "${SERVICE_NAME_FILES[@]}"; do
    if replace_in_file "$file" "$OLD_SERVICE_NAME" "$NEW_SERVICE_NAME"; then
        TOTAL_REPLACEMENTS=$((TOTAL_REPLACEMENTS + 1))
    fi
done

echo ""
echo "重命名文件..."

# 重命名 API 文件
if [ -f "${OLD_API_FILE}.api" ]; then
    mv "${OLD_API_FILE}.api" "${NEW_API_FILE}.api"
    echo "✓ ${OLD_API_FILE}.api -> ${NEW_API_FILE}.api"
    TOTAL_REPLACEMENTS=$((TOTAL_REPLACEMENTS + 1))
fi

# 重命名服务主文件
if [ -f "${OLD_SERVICE_FILE}.go" ]; then
    mv "${OLD_SERVICE_FILE}.go" "${NEW_SERVICE_FILE}.go"
    echo "✓ ${OLD_SERVICE_FILE}.go -> ${NEW_SERVICE_FILE}.go"
    TOTAL_REPLACEMENTS=$((TOTAL_REPLACEMENTS + 1))
fi

# 重命名配置文件
if [ -f "etc/${OLD_SERVICE_FILE}.yaml" ]; then
    mv "etc/${OLD_SERVICE_FILE}.yaml" "etc/${NEW_SERVICE_FILE}.yaml"
    echo "✓ etc/${OLD_SERVICE_FILE}.yaml -> etc/${NEW_SERVICE_FILE}.yaml"
    TOTAL_REPLACEMENTS=$((TOTAL_REPLACEMENTS + 1))
fi

echo ""
echo "更新文件名引用..."

# 更新主文件中的配置文件路径引用
if [ -f "${NEW_SERVICE_FILE}.go" ]; then
    replace_in_file "${NEW_SERVICE_FILE}.go" "etc/${OLD_SERVICE_FILE}.yaml" "etc/${NEW_SERVICE_FILE}.yaml"
fi

# 更新 makefile 中的文件名引用
for file in "${FILE_REFERENCE_FILES[@]}"; do
    if [ -f "$file" ]; then
        # 替换 API 文件名
        replace_in_file "$file" "${OLD_API_FILE}.api" "${NEW_API_FILE}.api"
        # 替换服务主文件名
        replace_in_file "$file" "${OLD_SERVICE_FILE}.go" "${NEW_SERVICE_FILE}.go"
    fi
done

echo ""
echo "-------------------------------------------"
echo "替换完成！已处理 $TOTAL_REPLACEMENTS 项"
echo ""
echo "后续步骤:"
echo "1. 运行 'make mt' 整理 Go 模块依赖"
echo "2. 运行 'make gen' 重新生成代码"
echo "3. 运行 'make run' 测试服务启动"
