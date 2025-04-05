#!/bin/bash -eu
# -e: 遇到错误时立即退出脚本（除非通过 try-catch 的方式处理）
# -u: 使用未定义变量时视为错误

# 设置时区为上海时间，便于生成的提交时间更符合本地习惯
export TZ='Asia/Shanghai'

# 初始化返回结果变量
result=0

echo "Using API_KEY: ${API_KEY:0:4}****"  # 打印前几位调试用，注意避免泄露

# 执行 collect-k8s-official-vulns 程序，参数说明：
# -u：更新漏洞
# -t：翻译漏洞
# -c：指定缓存目录
# 如果 collect-k8s-official-vulns 执行失败（退出码非0），将错误码赋值给 result 而不终止脚本
./collect-k8s-official-vulns -u -t -c "$REPOSITORY"  || result=$?

# 如果 collect-k8s-official-vulns 执行出错
if [ $result -ne 0 ]; then
  echo "[Err] Revert changes" >&2 # 打印错误信息到标准错误输出
  cd "$REPOSITORY" && git reset --hard HEAD # 回滚仓库的所有变更
  exit 1 # 脚本以失败退出，通知 GitHub Actions 工作流
fi

# 进入目标仓库目录，如果失败就退出脚本
cd "$REPOSITORY" || exit 1

# 检查仓库是否有变更（git status --porcelain 会输出简洁格式的状态）
if [[ -n $(git status --porcelain) ]]; then
  git add .  # 添加所有变更的文件
  git commit -m "update at $(date +'%Y-%m-%d %H:%M:%S')" # 提交变更，附带当前时间
  git push # 推送提交到远程仓库
fi