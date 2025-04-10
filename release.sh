#!/bin/bash

set -e

# 检查是否有 CHANGELOG.md 文件
if [ ! -f CHANGELOG.md ]; then
  echo "CHANGELOG.md not found."
  exit 1
fi

# 提取 CHANGELOG.md 中的最新版本号（兼容 macOS/Linux）
CURRENT_VERSION=$(grep -E '^##[[:space:]]+\[?v?([0-9]+\.[0-9]+\.[0-9]+)\]?' CHANGELOG.md | head -1 | sed -E 's/[^0-9]*([0-9]+\.[0-9]+\.[0-9]+).*/\1/')

if [ -z "$CURRENT_VERSION" ]; then
  echo "Failed to parse version from CHANGELOG.md"
  exit 1
fi

# 确认当前分支是 main 或 master
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [[ "$CURRENT_BRANCH" != "main" && "$CURRENT_BRANCH" != "master" ]]; then
  echo "Current branch is '$CURRENT_BRANCH'. Are you sure? [y/N]"
  read -r confirm
  if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo "Release aborted."
    exit 1
  fi
fi

# 添加 CHANGELOG.md 并提交
git add *
git commit -m "chore: release v$CURRENT_VERSION"

# 创建 Git tag
git tag -a "v$CURRENT_VERSION" -m "Release v$CURRENT_VERSION"

# 推送代码和 tag 到两个仓库
echo "Pushing to origin..."
git push origin "$CURRENT_BRANCH" --tags

echo "Pushing to github..."
git push github "$CURRENT_BRANCH" --tags

echo "Release v$CURRENT_VERSION completed!"
