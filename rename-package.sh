#!/bin/bash

# 检查参数
if [ $# -ne 1 ]; then
  echo "Usage: $0 <new-package-name>"
  echo "Example: $0 xxxx/xxx-api"
  exit 1
fi

OLD_PKG="cornyk/gin-template"
NEW_PKG=$1

# 获取脚本所在目录（项目根目录）
PROJECT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

echo "Renaming package from $OLD_PKG to $NEW_PKG in $PROJECT_ROOT"

# 1. 替换所有Go文件中的包引用
find "$PROJECT_ROOT" -type f \( -name "*.go" -o -name "go.mod" -o -name "*.yml" -o -name "*.yaml" -o -name "*.json" \) \
  -exec sed -i '' "s|$OLD_PKG|$NEW_PKG|g" {} +

# 2. 重命名go.mod文件中的模块名
if [ -f "$PROJECT_ROOT/go.mod" ]; then
  sed -i '' "s|module $OLD_PKG|module $NEW_PKG|g" "$PROJECT_ROOT/go.mod"
fi

# 3. 替换目录结构（如果新旧包路径深度一致）
OLD_PKG_PATH="${OLD_PKG//\//\/}"
NEW_PKG_PATH="${NEW_PKG//\//\/}"

if [[ "$OLD_PKG_PATH" != "$NEW_PKG_PATH" ]]; then
  echo "Warning: Package path structure changed. You may need to manually move directories."
  echo "Old path: $OLD_PKG_PATH"
  echo "New path: $NEW_PKG_PATH"
fi

# 4. 清理go.sum并重新下载依赖
if [ -f "$PROJECT_ROOT/go.sum" ]; then
  rm "$PROJECT_ROOT/go.sum"
fi

echo "Running go mod tidy..."
(cd "$PROJECT_ROOT" && go mod tidy)

echo "Rename completed!"
echo "Please check the changes with 'git status' before committing."
