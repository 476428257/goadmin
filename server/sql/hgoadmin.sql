SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for h_admin
-- ----------------------------
DROP TABLE IF EXISTS `h_admin`;
CREATE TABLE `h_admin` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `password` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码',
  `nickname` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '昵称',
  `email` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '邮箱',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态：1-启用，0-禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `version` int(10) DEFAULT '0' COMMENT '版本号',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `avatar` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '头像',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  KEY `idx_email` (`email`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员表';

-- ----------------------------
-- Records of h_admin
-- ----------------------------
INSERT INTO `h_admin` VALUES ('1', 'admin', '$2a$10$ZmybnfeAUBbDIYnooJmIFOuLK/v/k/WsxymieJxZo6L7PD2AK351G', '超级管理员', 'admin@1qq.com', '1', '2025-09-12 14:15:06', '2025-10-28 16:17:56', '28', '18888888888', '/image/head.png');
INSERT INTO `h_admin` VALUES ('2', 'admin2', '$2a$10$AKHF0Z0sZOZWS33s4hn/4OqE5EkDCu5b6JNZOfOlEhtfDTllL8Duu', '管理员', 'admin2@qq.com', '1', '2025-09-12 15:48:17', '2025-09-24 15:39:05', '4', '18888888888', '/image/head.png');

-- ----------------------------
-- Table structure for h_article
-- ----------------------------
DROP TABLE IF EXISTS `h_article`;
CREATE TABLE `h_article` (
  `id` int(12) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(1000) NOT NULL COMMENT '标题',
  `image` varchar(255) NOT NULL COMMENT '图片',
  `is_hot` tinyint(1) NOT NULL COMMENT '是否热门:0=否,1=是',
  `is_top` tinyint(1) NOT NULL COMMENT '是否置顶:0=否,1=是',
  `status` tinyint(1) NOT NULL COMMENT '状态:0=禁用,1=启用',
  `content` longtext NOT NULL COMMENT '内容',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of h_article
-- ----------------------------
INSERT INTO `h_article` VALUES ('1', '123', '/image/head.png', '1', '1', '1', '<p>123415</p>', '2025-09-09 10:16:13', '2025-09-24 13:58:18');

-- ----------------------------
-- Table structure for h_auth_role
-- ----------------------------
DROP TABLE IF EXISTS `h_auth_role`;
CREATE TABLE `h_auth_role` (
  `id` int(12) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `pid` int(12) unsigned NOT NULL DEFAULT '0' COMMENT '父级ID',
  `title` varchar(100) NOT NULL COMMENT '角色名称',
  `rule` text NOT NULL COMMENT '权限',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='权限角色表';

-- ----------------------------
-- Records of h_auth_role
-- ----------------------------
INSERT INTO `h_auth_role` VALUES ('1', '0', '超级管理员', '*', '1');
INSERT INTO `h_auth_role` VALUES ('2', '1', '管理员', '*', '1');

-- ----------------------------
-- Table structure for h_auth_rule
-- ----------------------------
DROP TABLE IF EXISTS `h_auth_rule`;
CREATE TABLE `h_auth_rule` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `pid` bigint(20) unsigned DEFAULT '0' COMMENT '父规则ID',
  `title` varchar(100) NOT NULL COMMENT '规则标题',
  `route` varchar(200) NOT NULL COMMENT '路由路径',
  `is_menu` tinyint(1) unsigned NOT NULL COMMENT '是否菜单',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态:0=禁用,1=启用',
  `icon` varchar(100) DEFAULT NULL COMMENT '图标',
  `weigh` int(10) NOT NULL DEFAULT '0' COMMENT '权重',
  `pagepath` varchar(100) DEFAULT NULL COMMENT '页面路径',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_auth_rule_route` (`route`)
) ENGINE=InnoDB AUTO_INCREMENT=67 DEFAULT CHARSET=utf8mb4 COMMENT='权限路由表';

-- ----------------------------
-- Records of h_auth_rule
-- ----------------------------
INSERT INTO `h_auth_rule` VALUES ('2', '20', '添加管理员', '/api/v1/admin/add', '0', '1', null, '0', null);
INSERT INTO `h_auth_rule` VALUES ('3', '20', '删除管理员', '/api/v1/admin/del', '0', '1', null, '0', null);
INSERT INTO `h_auth_rule` VALUES ('4', '20', '获取信息', '/api/v1/admin/info', '0', '1', null, '0', null);
INSERT INTO `h_auth_rule` VALUES ('5', '20', '管理员列表', '/api/v1/admin/list', '0', '1', '', '97', null);
INSERT INTO `h_auth_rule` VALUES ('6', '20', '修改管理员', '/api/v1/admin/update', '0', '1', null, '0', null);
INSERT INTO `h_auth_rule` VALUES ('8', '21', '添加角色组', '/api/v1/auth/role/add', '0', '1', null, '0', null);
INSERT INTO `h_auth_rule` VALUES ('9', '21', '删除角色组', '/api/v1/auth/role/del', '0', '1', null, '0', null);
INSERT INTO `h_auth_rule` VALUES ('10', '21', '查看信息', '/api/v1/auth/role/info', '0', '1', null, '0', null);
INSERT INTO `h_auth_rule` VALUES ('11', '21', '角色组列表', '/api/v1/auth/role/list', '0', '1', null, '96', null);
INSERT INTO `h_auth_rule` VALUES ('12', '21', '修改角色组', '/api/v1/auth/role/update', '0', '1', null, '0', null);
INSERT INTO `h_auth_rule` VALUES ('14', '22', '添加权限规则', '/api/v1/auth/rule/add', '0', '1', null, '0', null);
INSERT INTO `h_auth_rule` VALUES ('15', '22', '删除权限规则', '/api/v1/auth/rule/del', '0', '1', null, '0', null);
INSERT INTO `h_auth_rule` VALUES ('16', '22', '获取信息', '/api/v1/auth/rule/info', '0', '1', null, '0', null);
INSERT INTO `h_auth_rule` VALUES ('17', '22', '权限规则列表', '/api/v1/auth/rule/list', '0', '1', null, '95', null);
INSERT INTO `h_auth_rule` VALUES ('18', '22', '修改权限规则', '/api/v1/auth/rule/update', '0', '1', null, '0', null);
INSERT INTO `h_auth_rule` VALUES ('19', '0', '系统首页', '/dashboard', '1', '1', 'Odometer', '105', 'dashboard');
INSERT INTO `h_auth_rule` VALUES ('20', '0', '用户管理', '/system-user', '1', '1', 'User', '50', 'system/user');
INSERT INTO `h_auth_rule` VALUES ('21', '0', '角色管理', '/system-role', '1', '1', 'Grape', '49', 'system/role');
INSERT INTO `h_auth_rule` VALUES ('22', '0', '菜单管理', '/system-menu', '1', '1', 'MagicStick', '48', 'system/menu');
INSERT INTO `h_auth_rule` VALUES ('52', '0', '操作日志', '/czrz', '1', '1', 'Memo', '47', 'system/aclog');
INSERT INTO `h_auth_rule` VALUES ('53', '52', '日志列表', '/api/v1/system/aclog', '0', '1', '', '50', '');
INSERT INTO `h_auth_rule` VALUES ('54', '0', '文章管理', '/wenzhangguanli', '1', '1', 'Film', '46', 'system/article');
INSERT INTO `h_auth_rule` VALUES ('55', '54', '添加文章', '/api/v1/article/add', '0', '1', '', '0', '');
INSERT INTO `h_auth_rule` VALUES ('56', '54', '删除文章', '/api/v1/article/del', '0', '1', '', '0', '');
INSERT INTO `h_auth_rule` VALUES ('57', '54', '获取文章信息', '/api/v1/article/info', '0', '1', '', '0', '');
INSERT INTO `h_auth_rule` VALUES ('58', '54', '文章列表', '/api/v1/article/list', '0', '1', '', '97', '');
INSERT INTO `h_auth_rule` VALUES ('59', '54', '修改文章', '/api/v1/article/update', '0', '1', '', '0', '');
INSERT INTO `h_auth_rule` VALUES ('60', '54', '修改文章状态', '/api/v1/article/updatestatus', '0', '1', '', '0', '');
INSERT INTO `h_auth_rule` VALUES ('63', '19', '仪表盘', '/api/v1/system/dashboard', '0', '1', '', '105', '');
INSERT INTO `h_auth_rule` VALUES ('64', '0', '系统配置', '/config', '1', '1', 'Monitor', '104', 'system/config');
INSERT INTO `h_auth_rule` VALUES ('65', '64', '更新配置内容', '/api/v1/config/update', '0', '1', '', '104', '');

INSERT INTO `h_auth_rule` VALUES ('66', '64', '获取配置内容', '/api/v1/config/list', '0', '1', '', '104', '');

-- ----------------------------
-- Table structure for h_auth_user
-- ----------------------------
DROP TABLE IF EXISTS `h_auth_user`;
CREATE TABLE `h_auth_user` (
  `id` int(12) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` int(10) unsigned NOT NULL COMMENT '角色id',
  `admin_id` int(10) unsigned NOT NULL COMMENT '用户id',
  PRIMARY KEY (`id`),
  KEY `admin_id` (`admin_id`) USING BTREE,
  KEY `role_id` (`role_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of h_auth_user
-- ----------------------------
INSERT INTO `h_auth_user` VALUES ('1', '1', '1');
INSERT INTO `h_auth_user` VALUES ('2', '2', '2');

-- ----------------------------
-- Table structure for h_config
-- ----------------------------
DROP TABLE IF EXISTS `h_config`;
CREATE TABLE `h_config` (
  `id` int(12) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(100) NOT NULL COMMENT '变量名称',
  `title` varchar(100) NOT NULL COMMENT '标题',
  `value` text COMMENT '变量值',
  `group` varchar(100) NOT NULL COMMENT '分组',
  `type` varchar(100) NOT NULL COMMENT '类型',
  `content` text COMMENT '扩展内容',
  `extend` varchar(255) DEFAULT NULL COMMENT '扩展属性',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=108 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of h_config
-- ----------------------------
INSERT INTO `h_config` VALUES ('1', 'title', '网站标题', '后台管理系统', '系统配置', 'string', null, null);
INSERT INTO `h_config` VALUES ('2', 'logo', 'LOGO', '/image/logo.png', '系统配置', 'upload', null, '{\"multiple\":false}');
INSERT INTO `h_config` VALUES ('3', 'loginbg', '登录背景图', '/image/login-bg.png', '系统配置', 'upload', null, '{\"multiple\":false}');
INSERT INTO `h_config` VALUES ('4', 'userbg', '个人中心背景图', '/image/ucenter-bg.jpg', '系统配置', 'upload', '', '{\"multiple\":false}');
INSERT INTO `h_config` VALUES ('5', 'userhead', '默认头像', '/image/head.png', '系统配置', 'upload', '', '{\"multiple\":false}');
INSERT INTO `h_config` VALUES ('100', 'string', '字符', '', '配置示例', 'string', null, null);
INSERT INTO `h_config` VALUES ('101', 'number', '数字', '0', '配置示例', 'number', null, null);
INSERT INTO `h_config` VALUES ('102', 'select', '下拉列表', '', '配置示例', 'select', '[\"大\",\"中\",\"小\"]', null);
INSERT INTO `h_config` VALUES ('103', 'date', '日期', '', '配置示例', 'date', null, '{\"format\":\"YYYY-MM-DD HH:ii:ss\",\"type\":\"datetime\"}');
INSERT INTO `h_config` VALUES ('104', 'daterange', '日期区间', '', '配置示例', 'daterange', null, '{\"format\":\"YYYY-MM-DD HH:ii:ss\",\"type\":\"datetimerange\"}');
INSERT INTO `h_config` VALUES ('105', 'switch', '开关', '0', '配置示例', 'switch', null, '{\"activeValue\":\"1\",\"inactiveValue\":\"0\",\"activeText\":\"开启\",\"inactiveText\":\"关闭\"}');
INSERT INTO `h_config` VALUES ('106', 'upload', '文件上传', '', '配置示例', 'upload', '', '{\"multiple\":false}');
INSERT INTO `h_config` VALUES ('107', 'editor', '富文本', '&lt;p&gt;11&lt;/p&gt;', '配置示例', 'editor', null, '{\"height\":\"300px\"}');

-- ----------------------------
-- Table structure for h_operation_log
-- ----------------------------
DROP TABLE IF EXISTS `h_operation_log`;
CREATE TABLE `h_operation_log` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `admin_id` bigint(20) unsigned NOT NULL COMMENT '管理员ID',
  `username` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '管理员用户名',
  `path` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '请求路径',
  `ip` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP地址',
  `request_data` text COLLATE utf8mb4_unicode_ci COMMENT '请求数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_admin_id` (`admin_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

