/*
 Navicat Premium Data Transfer

 Source Server         : Shinewell
 Source Server Type    : MySQL
 Source Server Version : 80016 (8.0.16)
 Source Host           : localhost:3306
 Source Schema         : blog

 Target Server Type    : MySQL
 Target Server Version : 80016 (8.0.16)
 File Encoding         : 65001

 Date: 12/05/2026 13:04:39
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for article_categories
-- ----------------------------
DROP TABLE IF EXISTS `article_categories`;
CREATE TABLE `article_categories`  (
  `article_id` bigint(20) UNSIGNED NOT NULL COMMENT '文章ID',
  `category_id` bigint(20) UNSIGNED NOT NULL COMMENT '分类ID',
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`article_id`, `category_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '文章分类关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of article_categories
-- ----------------------------
INSERT INTO `article_categories` VALUES (1, 1, '2026-04-22 15:14:49.587');
INSERT INTO `article_categories` VALUES (2, 1, '2026-04-22 15:16:21.693');
INSERT INTO `article_categories` VALUES (3, 3, '2026-04-22 15:22:01.798');
INSERT INTO `article_categories` VALUES (4, 2, '2026-04-22 15:20:13.557');
INSERT INTO `article_categories` VALUES (5, 1, '2026-04-22 15:11:54.769');
INSERT INTO `article_categories` VALUES (6, 1, '2026-04-22 15:24:40.698');
INSERT INTO `article_categories` VALUES (7, 9, '2026-04-22 15:25:59.221');
INSERT INTO `article_categories` VALUES (8, 1, '2026-05-11 22:26:43.125');
INSERT INTO `article_categories` VALUES (8, 2, '2026-05-11 22:26:43.125');
INSERT INTO `article_categories` VALUES (8, 3, '2026-05-11 22:26:43.125');
INSERT INTO `article_categories` VALUES (8, 4, '2026-05-11 22:26:43.125');
INSERT INTO `article_categories` VALUES (8, 5, '2026-05-11 22:26:43.125');
INSERT INTO `article_categories` VALUES (9, 6, '2026-04-22 15:32:22.995');
INSERT INTO `article_categories` VALUES (10, 8, '2026-04-22 15:14:34.667');
INSERT INTO `article_categories` VALUES (12, 1, '2026-05-11 22:28:08.897');
INSERT INTO `article_categories` VALUES (12, 2, '2026-05-11 22:28:08.897');
INSERT INTO `article_categories` VALUES (12, 3, '2026-05-11 22:28:08.897');
INSERT INTO `article_categories` VALUES (12, 4, '2026-05-11 22:28:08.897');

-- ----------------------------
-- Table structure for articles
-- ----------------------------
DROP TABLE IF EXISTS `articles`;
CREATE TABLE `articles`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '文章ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '作者用户ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '标题',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '正文',
  `summary` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '摘要',
  `cover_image` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '封面图',
  `status` tinyint(4) NULL DEFAULT 0 COMMENT '状态 0草稿,1待审核,2已发布,3已驳回,4已封禁',
  `view_count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '浏览量',
  `like_count` bigint(20) NULL DEFAULT 0 COMMENT '点赞数',
  `favorite_count` bigint(20) NULL DEFAULT 0 COMMENT '收藏数',
  `comment_count` bigint(20) NULL DEFAULT 0 COMMENT '评论数',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `reject_reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '驳回原因',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_status_created`(`status` ASC, `created_at` ASC) USING BTREE,
  INDEX `idx_articles_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_articles_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '文章表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of articles
-- ----------------------------
INSERT INTO `articles` VALUES (1, 2, 'Gin 项目分层实践', '虽然大模型的技术更新的非常之快啊，但是在多年的技术进步与竞争淘汰中，有大量的技术和名词被沉淀了下来，这些名词大概率在未来是不会发生很大的变化的。因此，掌握它们是一件非常有必要的事情。理解这些名词的，不仅可以让你更加快速与方便的看懂各个大佬与 KOL 写的文章，也能让你更加清楚的知道当你想要使用这类技术的时候，应该怎么去搜或者是怎么询问大模型。毕竟正所谓人们不知道他们不知道的东西。\r\n\r\n2.1 什么是大模型——一栋庞大的摩天大楼\r\n大模型其实就是一个非常复杂的深度学习模型，它是一段占用大量内存的电脑程序。我们往往使用“参数”这个词来形容深度学习模型的大小，而大模型顾名思义就是参数量非常庞大的深度学习模型。而在如今这个时代，对于大部分的普通民众来说，他们基本上已将「大模型、深度学习、深度神经网络」和「AI」划上了等号，但其实前者之间更多是个包含的关系，不过这对于我们理解后续的内容而言并不重要。\r\n\r\n不理解什么是「非常复杂的深度学习模型」的朋友，我可以打一个比方。比如你们所在城市的摩天大厦，这个摩天大厦的每一层的每一个小隔间其实都不是很复杂，你甚至可以记住每一个小隔间的桌子椅子大概是怎么摆放的，并在一些模拟游戏里面将它复刻出来。但是整栋摩天大厦的复杂度就非常非常高，给你无限的材料，你几乎很难复原出整个摩天大楼里面的每一个细节。你可以将这里摩天大厦的每一层比作是一个非常基础的函数。比如说它可以是你们高中所学习到的1元2次函数，或者是一个稍微复杂一点的矩阵与矩阵操作的函数。而「复杂的深度学习模型」就是若干层这样的函数的堆叠，就如同这个复杂的摩天大厦，是无数个小隔间组成的无数个楼层的堆叠一样。稍微有点数学基础的朋友都知道函数的堆叠也是函数，所以你完全可以将深度学习模型等同于一个很复杂的函数，这两个是一个意思。至于具体的细节的话，感兴趣的朋友可以去学习一下，深度学习这门课。这里推荐大家可以看一下我的知乎专栏 机器学习幼儿园，无论从基础的机器学习概念，还是到稍微复杂一点的深度学习，还是到工程化的深度学习框架原理解析，一应俱全。\r\n\r\n大模型技术的研发最早可以追溯到 2015 年，各项基础技术与理论验证储备的完善是在 2020 年，第一个爆款的商业大模型 ChatGPT 则出现在 2022 年，基础模型的预训练能力在 2025 年被验证基本收敛（这一点目前学界与业界仍然有争议。但是基于目前已有的算力登记和公开高质量数据而言，事实上的通用文本大模型能力在 2025 年年中的时候被验证基本收敛），而能力基本收敛，你可以简单理解成在基础训练范式不变的情况下，未来出现的大模型和现在最厉害的大模型相比，再也不会出现质的上的能力飞跃了。\r\n\r\n而大模型基础能力的基本收敛是业内开始全面转向大模型应用开发，而非基础模型能力研发的一个重要标志，当然，这些都是我们后续需要了解的内容了，「如何将大模型做成 Agent？」，我们先按下不表。', '介绍 Gin 项目的 controller/service/repository 分层。', '/uploads/article/20260422151446_c1abb67d7d26.jpg', 2, 133, 0, 1, 2, '2026-04-02 09:00:00.000', '2026-04-22 15:37:43.280', '');
INSERT INTO `articles` VALUES (2, 3, 'GORM 更新操作避坑指南', '本文详细分析 Save 和 Updates 的区别。', '重点讲解 Save、Updates(struct)、Updates(map) 的差异。', '/uploads/article/20260422151618_95c1df1c6480.jpg', 2, 186, 0, 0, 2, '2026-04-02 09:30:00.000', '2026-04-22 15:37:40.152', '');
INSERT INTO `articles` VALUES (3, 4, 'MySQL 索引设计最佳实践', '围绕联合索引、覆盖索引与最左匹配原则展开。', '适合后端开发者的 MySQL 索引实战总结。', '/uploads/article/20260422152157_6b0b220fb2f5.png', 2, 265, 0, 0, 1, '2026-04-02 10:00:00.000', '2026-04-22 15:37:36.994', '');
INSERT INTO `articles` VALUES (4, 5, 'Vue3 组件通信方式总结', '本文总结 props、emit、provide/inject 等通信方式。这种“一块一块”的布局在前端开发中，最准确的术语通常叫作 卡片式设计。\r\n在设计系统（如 Ant Design、Material Design）中，它就是一个标准的 Card（卡片） 组件。\r\n为什么叫“卡片”\r\n这种设计模仿了现实生活中的名片或扑克牌：\r\n独立性：每一块（比如你截图里的每一篇文章）都是一个独立的容器，包含图片、标题、作者、时间等相关信息。\r\n视觉分割：通过阴影或边框，让这一块内容从背景中“浮”起来，显得更有层次感。\r\n实现这种布局的方式\r\n虽然“卡片”是指那个具体的块，但把它们排列得整整齐齐的布局方式，通常叫作：\r\n列表视图：虽然它是卡片，但因为是一行一行垂直排列的，有时也叫卡片列表。\r\n瀑布流：如果这些卡片的高度不统一，像 Pinterest 那样参差不齐地排列，就叫瀑布流（你截图里的卡片高度看起来比较统一，所以更偏向于普通的卡片列表）。', '快速掌握 Vue3 常见组件通信手段。', '/uploads/article/20260422152007_1fde85dbe430.jpg', 2, 100, 0, 0, 1, '2026-04-02 10:30:00.000', '2026-04-22 15:37:33.132', '');
INSERT INTO `articles` VALUES (5, 6, '从零理解 JWT 认证流程', '介绍 JWT 的签发、解析、续签与安全注意事项。', '适合做登录鉴权的入门与进阶文章。', '/uploads/article/20260422151025_169a5698d928.jpg', 2, 46, 0, 0, 0, '2026-04-02 11:00:00.000', '2026-04-22 15:12:43.384', '');
INSERT INTO `articles` VALUES (6, 7, 'Go 并发模式与常见陷阱', '讲解 goroutine、channel、context 的使用经验。', '帮助你写出更稳定的 Go 并发代码。', '/uploads/article/20260422152438_b8d7a0d066cc.png', 2, 321, 0, 0, 3, '2026-04-02 11:30:00.000', '2026-04-22 15:37:29.424', '');
INSERT INTO `articles` VALUES (7, 8, 'Docker 部署个人博客系统', '包括镜像构建、容器编排和环境配置。', '一篇面向初学者的 Docker 部署教程。', '/uploads/article/20260422152556_a518856cdc3b.jpg', 2, 36, 0, 0, 0, '2026-04-02 12:00:00.000', '2026-04-22 15:37:26.522', '');
INSERT INTO `articles` VALUES (8, 9, 'Redis 缓存击穿与雪崩处理', '介绍缓存穿透、击穿、雪崩的区别与应对方案。', '后端系统常见缓存问题的解决思路。', '/uploads/article/20260422152658_40de40a38d47.jpg', 2, 232, 0, 1, 1, '2026-04-02 12:30:00.000', '2026-04-22 15:37:17.781', '');
INSERT INTO `articles` VALUES (9, 10, '区块链钱包基础概念梳理', '围绕地址、私钥、公钥、签名等概念展开。', '适合区块链初学者入门。', '/uploads/article/20260422153220_0c60abeec6f5.png', 2, 30, 0, 0, 2, '2026-04-02 13:00:00.000', '2026-04-22 15:37:13.479', '');
INSERT INTO `articles` VALUES (10, 2, '写给程序员的技术写作建议', '从选题、结构、表达和复盘几个角度分享经验。', '帮助程序员写出更有价值的技术文章。', '/uploads/article/20260422151430_3e778f395a0a.webp', 2, 254, 1, 1, 6, '2026-04-02 13:30:00.000', '2026-04-22 15:36:58.938', '');
INSERT INTO `articles` VALUES (12, 1, '问问', '<p>七五一额骨<br><span style=\"color: rgb(77, 77, 76); background-color: rgb(238, 238, 238);\">func OrderStatus(status []string) func (db *gorm.DB) *gorm.DB {<br> return func (db *gorm.DB) *gorm.DB {<br> return db.Where(\"status IN (?)\", status)<br> }<br>}<br><br>db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&amp;orders)<br>// 查找所有金额大于 1000 的信用卡订单<br><br>db.Scopes(AmountGreaterThan1000, PaidWithCod).Find(&amp;orders)<br>// 查找所有金额大于 1000 的 COD 订单<br><br>db.Scopes(AmountGreaterThan1000, OrderStatus([]string{\"paid\", \"shipped\"})).Find(&amp;orders)<br>// 查找所有金额大于1000 的已付款或已发货订单<br></span><img src=\"/uploads/article-content/20260511202118_702cce7f1f47.png\" alt=\"image.png\" data-href=\"/uploads/article-content/20260511202118_702cce7f1f47.png\" style=\"\"/></p>', '<p>七五一额骨<br><span style=\"color: rgb(77, 77, 76); background-color: rgb(238, 238, 238);\">func OrderSt', '', 2, 7, 0, 0, 0, '2026-05-11 20:21:49.330', '2026-05-11 20:21:54.893', '');

-- ----------------------------
-- Table structure for categories
-- ----------------------------
DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '分类名称',
  `slug` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'URL 标识',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_categories_slug`(`slug` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '分类表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of categories
-- ----------------------------
INSERT INTO `categories` VALUES (1, '后端开发', 'backend', '2026-04-01 11:00:00.000', '2026-04-01 11:00:00.000');
INSERT INTO `categories` VALUES (2, '前端开发', 'frontend', '2026-04-01 11:01:00.000', '2026-04-01 11:01:00.000');
INSERT INTO `categories` VALUES (3, '数据库', 'database', '2026-04-01 11:02:00.000', '2026-04-01 11:02:00.000');
INSERT INTO `categories` VALUES (4, '云计算', 'cloud', '2026-04-01 11:03:00.000', '2026-04-01 11:03:00.000');
INSERT INTO `categories` VALUES (5, '人工智能', 'ai', '2026-04-01 11:04:00.000', '2026-04-01 11:04:00.000');
INSERT INTO `categories` VALUES (6, '区块链', 'blockchain', '2026-04-01 11:05:00.000', '2026-04-01 11:05:00.000');
INSERT INTO `categories` VALUES (7, '算法', 'algorithm', '2026-04-01 11:06:00.000', '2026-04-01 11:06:00.000');
INSERT INTO `categories` VALUES (8, '随笔', 'essay', '2026-04-01 11:07:00.000', '2026-04-01 11:07:00.000');
INSERT INTO `categories` VALUES (9, '运维', 'devops', '2026-04-01 11:08:00.000', '2026-04-01 11:08:00.000');
INSERT INTO `categories` VALUES (10, '产品思考', 'product', '2026-04-01 11:09:00.000', '2026-04-01 11:09:00.000');

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `article_id` bigint(20) UNSIGNED NOT NULL COMMENT '文章ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '评论用户ID',
  `parent_id` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '父评论ID，一级为空',
  `root_id` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '一级评论ID，一级为空',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '内容',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '0待审核,1已发布,2已拒绝',
  `like_count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
  `reply_count` bigint(20) NULL DEFAULT 0 COMMENT '直接回复数',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_article_status_time`(`article_id` ASC, `status` ASC, `created_at` ASC) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_root_id`(`root_id` ASC) USING BTREE,
  INDEX `idx_comments_article_id`(`article_id` ASC) USING BTREE,
  INDEX `idx_comments_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_comments_parent_id`(`parent_id` ASC) USING BTREE,
  INDEX `idx_comments_root_id`(`root_id` ASC) USING BTREE,
  INDEX `idx_comments_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 24 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '评论表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of comments
-- ----------------------------
INSERT INTO `comments` VALUES (1, 1, 3, NULL, NULL, '这篇文章的分层思路很清晰。', 1, 0, 1, '2026-04-03 09:00:00.000', '2026-04-03 09:00:00.000');
INSERT INTO `comments` VALUES (2, 1, 2, 1, 1, '谢谢认可，后续我会补充更多示例。', 1, 0, 0, '2026-04-03 09:05:00.000', '2026-04-03 09:05:00.000');
INSERT INTO `comments` VALUES (3, 2, 4, NULL, NULL, 'Save 和 Updates 的区别终于搞懂了。', 1, 0, 0, '2026-04-03 09:10:00.000', '2026-04-03 09:10:00.000');
INSERT INTO `comments` VALUES (4, 2, 5, NULL, NULL, '建议再补充一下 Select 和 Omit 的用法。', 1, 0, 0, '2026-04-03 09:15:00.000', '2026-04-03 09:15:00.000');
INSERT INTO `comments` VALUES (5, 3, 6, NULL, NULL, '索引部分总结得很好，受教了。', 1, 0, 0, '2026-04-03 09:20:00.000', '2026-04-03 09:20:00.000');
INSERT INTO `comments` VALUES (6, 4, 7, NULL, NULL, 'Vue3 这篇适合初学者阅读。', 1, 0, 0, '2026-04-03 09:25:00.000', '2026-04-03 09:25:00.000');
INSERT INTO `comments` VALUES (7, 6, 8, NULL, NULL, '并发部分例子很实用。', 1, 0, 1, '2026-04-03 09:30:00.000', '2026-04-03 09:30:00.000');
INSERT INTO `comments` VALUES (8, 6, 7, 7, 7, '特别是 context 的取消传播写得很棒。', 1, 0, 1, '2026-04-03 09:35:00.000', '2026-04-03 09:35:00.000');
INSERT INTO `comments` VALUES (9, 8, 9, NULL, NULL, '缓存击穿和雪崩这两个概念总算分清了。', 1, 0, 0, '2026-04-03 09:40:00.000', '2026-04-03 09:40:00.000');
INSERT INTO `comments` VALUES (10, 10, 10, NULL, NULL, '技术写作建议很有启发。', 1, 0, 2, '2026-04-03 09:45:00.000', '2026-04-03 09:45:00.000');
INSERT INTO `comments` VALUES (11, 10, 6, NULL, NULL, '文章666', 1, 0, 0, '2026-04-21 16:08:52.504', '2026-04-21 16:08:52.504');
INSERT INTO `comments` VALUES (12, 10, 6, 10, 10, '就在codex 大规模普及前，我专心学了图形学、threejs、Blender，期望能构筑一点点技术壁垒，防止公司优化。并且在去年下半年期间，确实在项目上用到了，有了一点点不可替代性。毕竟，在前端方向，做这种Web3D还是有一定的稀缺性的。', 1, 0, 1, '2026-04-21 20:27:07.126', '2026-04-21 20:27:07.126');
INSERT INTO `comments` VALUES (13, 10, 6, 12, 10, '一款冷板凳开发的深度定制WordPress极简博客主题：小熊日记', 1, 0, 0, '2026-04-21 20:28:06.382', '2026-04-21 20:28:06.382');
INSERT INTO `comments` VALUES (15, 10, 1, NULL, NULL, '外耳', 1, 0, 0, '2026-05-11 20:32:51.739', '2026-05-11 20:32:51.739');
INSERT INTO `comments` VALUES (16, 10, 2, NULL, NULL, '赵敬超sy', 1, 0, 0, '2026-05-11 20:38:51.689', '2026-05-11 20:38:51.689');
INSERT INTO `comments` VALUES (18, 10, 2, 10, 10, '未激活', 1, 0, 0, '2026-05-11 20:39:41.746', '2026-05-11 20:39:41.746');
INSERT INTO `comments` VALUES (21, 9, 1, NULL, NULL, '问吾问无为谓', 1, 0, 0, '2026-05-11 22:12:52.237', '2026-05-11 22:12:52.237');
INSERT INTO `comments` VALUES (22, 9, 1, NULL, NULL, '哇哦阿伟', 1, 0, 0, '2026-05-11 22:12:57.079', '2026-05-11 22:12:57.079');
INSERT INTO `comments` VALUES (23, 6, 1, 8, 7, 'channel、context 的使用经验。', 1, 0, 0, '2026-05-11 22:13:23.982', '2026-05-11 22:13:23.982');

-- ----------------------------
-- Table structure for favorites
-- ----------------------------
DROP TABLE IF EXISTS `favorites`;
CREATE TABLE `favorites`  (
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `article_id` bigint(20) UNSIGNED NOT NULL COMMENT '文章ID',
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`, `article_id`) USING BTREE,
  INDEX `idx_article_id`(`article_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '收藏表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of favorites
-- ----------------------------
INSERT INTO `favorites` VALUES (1, 1, '2026-05-04 21:49:10.692');
INSERT INTO `favorites` VALUES (1, 8, '2026-05-11 22:44:03.285');
INSERT INTO `favorites` VALUES (1, 10, '2026-05-11 20:32:08.653');

-- ----------------------------
-- Table structure for images
-- ----------------------------
DROP TABLE IF EXISTS `images`;
CREATE TABLE `images`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '图片ID',
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '上传者ID',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '图片URL',
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_images_user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '图片表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of images
-- ----------------------------
INSERT INTO `images` VALUES (1, 2, '/upload/大笑.jpg', '2026-04-04 10:00:00.000');
INSERT INTO `images` VALUES (2, 3, '/upload/man.png', '2026-04-04 10:01:00.000');
INSERT INTO `images` VALUES (3, 4, '/upload/sakyua1.png', '2026-04-04 10:02:00.000');
INSERT INTO `images` VALUES (4, 5, '/upload/sakyua2.jpg', '2026-04-04 10:03:00.000');
INSERT INTO `images` VALUES (5, 6, '/upload/红温1.jpg', '2026-04-04 10:04:00.000');
INSERT INTO `images` VALUES (6, 7, '/upload/视奸.jpg', '2026-04-04 10:05:00.000');
INSERT INTO `images` VALUES (7, 8, '/upload/问号.jpg', '2026-04-04 10:06:00.000');
INSERT INTO `images` VALUES (8, 9, '/upload/beng.png', '2026-04-04 10:07:00.000');
INSERT INTO `images` VALUES (9, 10, '/upload/甄子丹2.jpg', '2026-04-04 10:08:00.000');
INSERT INTO `images` VALUES (10, 2, '/upload/nn.jpg', '2026-04-04 10:09:00.000');

-- ----------------------------
-- Table structure for likes
-- ----------------------------
DROP TABLE IF EXISTS `likes`;
CREATE TABLE `likes`  (
  `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `target_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '目标类型: article, comment, image',
  `target_id` bigint(20) UNSIGNED NOT NULL COMMENT '目标ID',
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`, `target_type`, `target_id`) USING BTREE,
  INDEX `idx_target`(`target_type` ASC, `target_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '点赞表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of likes
-- ----------------------------
INSERT INTO `likes` VALUES (1, 'article', 10, '2026-05-11 20:29:27.273');

-- ----------------------------
-- Table structure for review_logs
-- ----------------------------
DROP TABLE IF EXISTS `review_logs`;
CREATE TABLE `review_logs`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `article_id` bigint(20) UNSIGNED NOT NULL COMMENT '文章ID',
  `admin_id` bigint(20) UNSIGNED NOT NULL COMMENT '审核人ID',
  `action` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'approve/reject/ban',
  `reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '原因/备注',
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_review_logs_article_id`(`article_id` ASC) USING BTREE,
  INDEX `idx_review_logs_admin_id`(`admin_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 26 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of review_logs
-- ----------------------------
INSERT INTO `review_logs` VALUES (1, 1, 1, 'approve', '内容完整，审核通过', '2026-04-06 09:00:00.000');
INSERT INTO `review_logs` VALUES (2, 2, 1, 'approve', '技术内容清晰，审核通过', '2026-04-06 09:05:00.000');
INSERT INTO `review_logs` VALUES (3, 3, 1, 'approve', '适合发布到首页推荐', '2026-04-06 09:10:00.000');
INSERT INTO `review_logs` VALUES (4, 4, 1, 'approve', '结构清晰，审核通过', '2026-04-06 09:15:00.000');
INSERT INTO `review_logs` VALUES (5, 5, 1, 'reject', '部分段落表达不清，请补充示例', '2026-04-06 09:20:00.000');
INSERT INTO `review_logs` VALUES (6, 6, 1, 'approve', '并发示例较好，审核通过', '2026-04-06 09:25:00.000');
INSERT INTO `review_logs` VALUES (7, 7, 1, 'reject', '图片排版混乱，内容深度不足', '2026-04-06 09:30:00.000');
INSERT INTO `review_logs` VALUES (8, 8, 1, 'approve', '缓存主题实用，审核通过', '2026-04-06 09:35:00.000');
INSERT INTO `review_logs` VALUES (9, 9, 1, 'reject', '草稿内容尚未完善', '2026-04-06 09:40:00.000');
INSERT INTO `review_logs` VALUES (10, 10, 1, 'approve', '适合作为写作经验分享', '2026-04-06 09:45:00.000');
INSERT INTO `review_logs` VALUES (11, 5, 1, 'approve', '', '2026-04-22 15:12:43.389');
INSERT INTO `review_logs` VALUES (12, 10, 1, 'approve', '', '2026-04-22 15:36:58.940');
INSERT INTO `review_logs` VALUES (13, 9, 1, 'approve', '', '2026-04-22 15:37:13.481');
INSERT INTO `review_logs` VALUES (14, 8, 1, 'approve', '', '2026-04-22 15:37:17.782');
INSERT INTO `review_logs` VALUES (15, 7, 1, 'approve', '', '2026-04-22 15:37:26.523');
INSERT INTO `review_logs` VALUES (16, 6, 1, 'approve', '', '2026-04-22 15:37:29.426');
INSERT INTO `review_logs` VALUES (17, 4, 1, 'approve', '', '2026-04-22 15:37:33.134');
INSERT INTO `review_logs` VALUES (18, 3, 1, 'approve', '', '2026-04-22 15:37:36.995');
INSERT INTO `review_logs` VALUES (19, 2, 1, 'approve', '', '2026-04-22 15:37:40.153');
INSERT INTO `review_logs` VALUES (20, 1, 1, 'approve', '', '2026-04-22 15:37:43.280');
INSERT INTO `review_logs` VALUES (21, 12, 1, 'approve', '', '2026-05-11 20:21:54.893');
INSERT INTO `review_logs` VALUES (22, 12, 1, 'category', '', '2026-05-11 20:36:02.108');
INSERT INTO `review_logs` VALUES (23, 12, 1, 'category', '', '2026-05-11 20:36:09.160');
INSERT INTO `review_logs` VALUES (24, 8, 1, 'category', '', '2026-05-11 22:26:43.125');
INSERT INTO `review_logs` VALUES (25, 12, 1, 'category', '', '2026-05-11 22:28:08.897');

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '邮箱',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '头像',
  `bio` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '个人简介',
  `role` tinyint(4) NOT NULL DEFAULT 1 COMMENT '0管理员,1普通用户',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '昵称',
  `status` tinyint(4) NULL DEFAULT 1 COMMENT '状态:1正常,2禁用',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_users_username`(`username` ASC) USING BTREE,
  UNIQUE INDEX `idx_users_email`(`email` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, 'super', '$2a$10$gNe3tDfmqyByjp8xIhThaeccjrlVGAT4lKSs4lN4GcePIQJbKguoG', '2302635277@qq.com', '/uploads/user/20260419154136_3aa4b8fb4e4d.jpg', '随便看看, 随便玩玩', 0, '管理员', 1, '2026-04-01 09:00:00.000', '2026-04-21 15:46:59.552');
INSERT INTO `users` VALUES (3, 'bob', '$2a$10$gNe3tDfmqyByjp8xIhThaeccjrlVGAT4lKSs4lN4GcePIQJbKguoG', 'bob@example.com', '/uploads/user/20260422151550_a2210a70c0fe.gif', '专注 Go 后端开发', 1, '河南彭于晏', 1, '2026-04-01 09:20:00.000', '2026-04-22 15:15:53.083');
INSERT INTO `users` VALUES (4, 'charlie', '$2a$10$gNe3tDfmqyByjp8xIhThaeccjrlVGAT4lKSs4lN4GcePIQJbKguoG', 'charlie@example.com', '/uploads/user/20260422152143_1306986ea4c2.jpg', '喜欢数据库和系统设计', 1, '爱吃香菜', 1, '2026-04-01 09:30:00.000', '2026-04-22 15:21:45.364');
INSERT INTO `users` VALUES (5, 'david', '$2a$10$gNe3tDfmqyByjp8xIhThaeccjrlVGAT4lKSs4lN4GcePIQJbKguoG', 'david@example.com', '/uploads/user/20260422151928_ca795134b72e.webp', '前端与全栈爱好者', 1, '不吃肘击', 1, '2026-04-01 09:40:00.000', '2026-04-22 15:19:33.688');
INSERT INTO `users` VALUES (6, 'dxj', '$2a$10$gNe3tDfmqyByjp8xIhThaeccjrlVGAT4lKSs4lN4GcePIQJbKguoG', '123456789@qq.com', '/uploads/user/20260421155548_7ed7325ec6b9.jpg', '记录生活与技术成长', 1, '黑白魔女', 1, '2026-04-01 09:50:00.000', '2026-04-21 15:56:24.872');
INSERT INTO `users` VALUES (7, 'frank', '$2a$10$gNe3tDfmqyByjp8xIhThaeccjrlVGAT4lKSs4lN4GcePIQJbKguoG', 'frank@example.com', '/uploads/user/20260422152418_f9a63a283201.jpg', '关注性能优化', 1, '忧郁荷包蛋', 1, '2026-04-01 10:00:00.000', '2026-04-22 15:24:20.603');
INSERT INTO `users` VALUES (8, 'grace', '$2a$10$gNe3tDfmqyByjp8xIhThaeccjrlVGAT4lKSs4lN4GcePIQJbKguoG', 'grace@example.com', '/uploads/user/20260422152521_1ab3b9f6e8e5.jpg', '喜欢读书和分享', 1, '番茄鸡蛋仔', 1, '2026-04-01 10:10:00.000', '2026-04-22 15:25:28.102');
INSERT INTO `users` VALUES (9, 'henry', '$2a$10$gNe3tDfmqyByjp8xIhThaeccjrlVGAT4lKSs4lN4GcePIQJbKguoG', 'henry@example.com', '/uploads/user/20260422152644_37835e9c3227.png', '一名测试工程师', 1, '悲伤黄焖鸡', 1, '2026-04-01 10:20:00.000', '2026-04-22 15:26:47.851');
INSERT INTO `users` VALUES (10, 'leon', '$2a$10$gNe3tDfmqyByjp8xIhThaeccjrlVGAT4lKSs4lN4GcePIQJbKguoG', 'ivy@example.com', '/uploads/user/20260422153134_61c724d819f3.jpg', '热爱产品与运营', 1, '李昂', 1, '2026-04-01 10:30:00.000', '2026-04-22 15:31:52.525');

SET FOREIGN_KEY_CHECKS = 1;
