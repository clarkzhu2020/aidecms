#!/bin/bash

# AideCMS Web 部署脚本
# 用法: ./deploy.sh [production|staging]

set -e

ENV=${1:-production}

echo "🚀 开始部署 AideCMS 官网 ($ENV 环境)..."

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 安装依赖
echo -e "${YELLOW}📦 安装依赖...${NC}"
npm install

# 构建
echo -e "${YELLOW}🔨 构建项目...${NC}"
npm run build

# 检查构建结果
if [ ! -d "dist" ]; then
    echo "❌ 构建失败：dist 目录不存在"
    exit 1
fi

echo -e "${GREEN}✅ 构建成功！${NC}"

# 根据环境部署
case $ENV in
    production)
        echo -e "${YELLOW}📤 部署到生产环境...${NC}"
        # 示例：使用 rsync 部署
        # rsync -avz --delete dist/ user@server:/var/www/aidecms-web/dist/
        
        # 示例：使用 SCP 部署
        # scp -r dist/* user@server:/var/www/aidecms-web/dist/
        
        # 示例：部署到 Vercel
        # npm install -g vercel
        # vercel --prod
        
        echo "📝 请配置您的部署命令"
        ;;
        
    staging)
        echo -e "${YELLOW}📤 部署到测试环境...${NC}"
        # 测试环境部署命令
        # rsync -avz --delete dist/ user@staging-server:/var/www/aidecms-web/dist/
        echo "📝 请配置您的测试环境部署命令"
        ;;
        
    *)
        echo "❌ 未知环境: $ENV"
        echo "用法: ./deploy.sh [production|staging]"
        exit 1
        ;;
esac

echo -e "${GREEN}🎉 部署完成！${NC}"
echo ""
echo "📋 构建信息:"
echo "  - 构建时间: $(date)"
echo "  - 环境: $ENV"
echo "  - 构建大小: $(du -sh dist | cut -f1)"
echo ""
echo "🌐 下一步:"
echo "  1. 检查 dist 目录内容"
echo "  2. 配置 Nginx 或其他 Web 服务器"
echo "  3. 访问您的域名验证部署"
