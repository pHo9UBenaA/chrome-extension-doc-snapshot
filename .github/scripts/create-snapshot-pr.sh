#!/bin/bash

# Gitユーザー設定
git config --local user.email "41898282+github-actions[bot]@users.noreply.github.com"
git config --local user.name "github-actions[bot]"

# ブランチ名を設定（日付を含める）
BRANCH_NAME="snapshot-update-$(date +%Y%m%d-%H%M%S)"

# 新しいブランチを作成
git switch -c $BRANCH_NAME

# 変更をコミット
git add __snapshot__
git commit -m "Auto update snapshot $(date +%Y-%m-%d)"

# ブランチをプッシュ
git push origin $BRANCH_NAME

# プルリクエスト作成
gh pr create --title "Auto update snapshot $(date +%Y-%m-%d)" --body "This PR is automatically created by github actions." --base master --head $BRANCH_NAME 
