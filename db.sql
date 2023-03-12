create database operation;

use operation;

# tb_system_menu 菜单表
CREATE TABLE `tb_system_menu`(
                                 `id` int(11) NOT NULL Primary key  AUTO_INCREMENT COMMENT '主键',
                                 `title` varchar(128) NOT NULL COMMENT '菜单标题',
                                 `icon` varchar(128) DEFAULT NULL COMMENT '图标',
                                 `path` varchar(128) binary DEFAULT NULL COMMENT '路径',
                                 `name` varchar(128) DEFAULT NULL COMMENT '菜单name',
                                 `type` varchar(128) DEFAULT NULL COMMENT '菜单类型',
                                 `permission` varchar(128) DEFAULT NULL COMMENT '指令',
                                 `method` varchar(128) DEFAULT NULL COMMENT '接口请求方式',
                                 `component` varchar(128) DEFAULT NULL COMMENT '组件地址',
                                 `redirect` varchar(128) DEFAULT NULL COMMENT '重定向地址',
                                 `parent_id` int(11) NOT NULL COMMENT '父级菜单ID',
                                 `is_hidden` tinyint(1) DEFAULT 0 COMMENT '是否隐藏',
                                 `is_cache` tinyint(1) DEFAULT 0 COMMENT '是否缓存页面',
                                 `weight` int(11) DEFAULT 0 COMMENT '菜单权重',
                                 `operator` varchar(128) NOT NULL COMMENT '操作人员',
                                 `operator_id` int(11) NOT NULL COMMENT '操作人员ID',
                                 `created_at` int(11) DEFAULT NULL COMMENT '创建时间',
                                 `updated_at` int(11) DEFAULT NULL COMMENT '修改时间'
) ENGINE=InnoDB CHARSET=utf8 COLLATE=utf8_unicode_ci;

# tb_system_role 角色表
CREATE TABLE `tb_system_role` (
                                  `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
                                  `parent_id` int(11) not null COMMENT '父角色id',
                                  `name` varchar(128) NOT NULL COMMENT '角色名称',
                                  `keyword` varchar(128) binary NOT NULL COMMENT '角色关键字',
                                  `status` tinyint(1) NOT NULL COMMENT '角色状态',
                                  `weight` int(11) DEFAULT '0' COMMENT '角色权重',
                                  `description` varchar(300) DEFAULT NULL COMMENT '角色备注',
                                  `data_scope` varchar(128) NOT NULL COMMENT '数据权限',
                                  `team_ids` text COMMENT '自定义权限部门id',
                                  `operator` varchar(128)  NOT NULL COMMENT '操作人员',
                                  `operator_id` int(11) NOT NULL COMMENT '操作人员ID',
                                  `created_at` int(11) DEFAULT NULL,
                                  `updated_at` int(11) DEFAULT NULL,
                                  UNIQUE  index  (`keyword`)
) ENGINE=InnoDB  CHARSET=utf8 COLLATE=utf8_unicode_ci;

# tb_system_role_menu 菜单表
CREATE TABLE `tb_system_role_menu` (
                                       `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
                                       `role_id` int(11) NOT NULL COMMENT '角色ID',
                                       `menu_id` int(11) NOT NULL COMMENT '菜单ID',
                                       `operator` varchar(128) NOT NULL COMMENT '操作人员',
                                       `operator_id` int(11) NOT NULL COMMENT '操作人员ID',
                                       `created_at` int(11) DEFAULT NULL,
                                       `updated_at` int(11) DEFAULT NULL,
                                       FOREIGN KEY (`role_id`) REFERENCES `tb_system_role` (`id`) ON DELETE CASCADE,
                                       FOREIGN KEY (`menu_id`) REFERENCES `tb_system_menu` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB CHARSET=utf8 COLLATE=utf8_unicode_ci;

# tb_system_team 部门表
CREATE TABLE `tb_system_team` (
                                  `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
                                  `name` varchar(128) NOT NULL COMMENT '部门名称',
                                  `description` varchar(300) DEFAULT NULL COMMENT '部门备注',
                                  `parent_id` int(11) NOT NULL COMMENT '上级ID',
                                  `operator` varchar(128) NOT NULL COMMENT '操作人员',
                                  `operator_id` int(11) NOT NULL COMMENT '操作人员ID',
                                  `created_at` int(11) DEFAULT NULL,
                                  `updated_at` int(11) DEFAULT NULL,
                                  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB CHARSET=utf8 COLLATE=utf8_unicode_ci;

# tb_system_user 用户表
CREATE TABLE `tb_system_user` (
                                  `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
                                  `role_id` int(11) NOT NULL COMMENT '角色ID',
                                  `team_id` int(11) NOT NULL COMMENT '部门ID',
                                  `nickname` varchar(128)  NOT NULL COMMENT '用户昵称',
                                  `name` varchar(32) NOT NULL COMMENT '用户姓名',
                                  `phone` varchar(32) binary NOT NULL COMMENT '用户电话',
                                  `avatar` varchar(128) NOT NULL COMMENT '用户头像',
                                  `email` varchar(128) binary NOT NULL COMMENT '用户邮箱',
                                  `sex` tinyint(1) NOT NULL COMMENT '用户性别',
                                  `password` varchar(300) binary NOT NULL COMMENT '用户密码',
                                  `status` tinyint(1) NOT NULL COMMENT '用户状态', # 用户状态封禁 ｜ 启用
                                      `disable_desc` varchar(128) DEFAULT NULL COMMENT '禁用原因',
                                  `last_login` int(11) DEFAULT NULL COMMENT '最后登陆时间',
                                  `operator` varchar(128) NOT NULL COMMENT '操作人员',
                                  `operator_id` int(11) NOT NULL COMMENT '操作人员ID',
                                  `created_at` int(11) DEFAULT NULL,
                                  `updated_at` int(11) DEFAULT NULL,
                                  unique index(phone),
                                  unique index(email),
                                  FOREIGN KEY (`role_id`) REFERENCES `tb_system_role` (`id`),
                                  FOREIGN KEY (`team_id`) REFERENCES `tb_system_team` (`id`)
) ENGINE=InnoDB CHARSET=utf8 COLLATE=utf8_unicode_ci;

# tb_system_login_log 登陆日志表
CREATE TABLE `tb_system_login_log` (
                                       `id` int(11) NOT NULL  PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
                                       `phone` varchar(128) binary NOT NULL COMMENT '用户账号',
                                       `ip` char(32) NOT NULL COMMENT 'IP地址',
                                       `address` varchar(256) NOT NULL COMMENT '登陆地址',
                                       `browser` varchar(128) NOT NULL COMMENT '浏览器',
                                       `device` varchar(128) NOT NULL COMMENT '登录设备',
                                       `status` tinyint(1) NOT NULL COMMENT '登录状态',
                                       `code` int(11) NOT NULL COMMENT '错误码',
                                       `description` varchar(256) NOT NULL COMMENT '登录备注',
                                       `created_at` int(11) DEFAULT NULL COMMENT '登陆时间',
                                       index(created_at),
                                       index(phone)
) ENGINE=InnoDB  CHARSET=utf8;

# 创建菜单表
INSERT INTO `tb_system_menu` VALUES (1,'根菜单',NULL,'/',NULL,'R',NULL,NULL,NULL,NULL,0,0,0,0,'system',0,1234567,1234567),(2,'系统管理','settings','/system','System','M','','','Layout','/system/user',1,0,1,0,'方伟业',1,1676616568,1676616568),(4,'基本接口','apps','/baseApi','baseApi','M','','','','',1,1,0,1,'方伟业',1,1676684535,1676685106),(5,'系统管理基础接口','apps','/baseApi','baseApi','M','','','','',4,1,0,0,'方伟业',1,1676684709,1676684897),(6,'获取当前用户信息','','/api/system/user','','A','baseApi','GET','','',5,0,0,0,'方伟业',1,1676684954,1676685013),(7,'获取当前用户菜单','','/api/system/user/menus','','A','baseApi','GET','','',5,0,0,0,'方伟业',1,1676685004,1676685004),(8,'获取系统部门信息','','/api/system/teams','','A','baseApi','GET','','',5,0,0,0,'方伟业',1,1676685078,1676685078),(9,'菜单管理','menu','menu','systemMenu','M','','','system/menu/index','',2,0,1,0,'方伟业',1,1676685381,1676685381),(10,'查看菜单','','/api/system/menus','','A','system:menu:query','GET','','',9,0,0,0,'方伟业',1,1676685528,1676685528),(11,'新增菜单','','/api/system/menu','','A','system:menu:add','POST','','',9,0,0,0,'方伟业',1,1676685599,1676685599),(12,'修改菜单','','/api/system/menu','','A','system:menu:update','PUT','','',9,0,0,0,'方伟业',1,1676685632,1676685632),(13,'删除菜单','','/api/system/menu','','A','system:menu:delete','DELETE','','',9,0,0,0,'方伟业',1,1676685657,1676685657),(14,'部门管理','user-group','team','sysTeam','M','','','system/team/index','',2,0,0,0,'方伟业',1,1676686013,1676691441),(15,'新增部门','','/api/system/team','','A','system:team:add','POST','','',14,0,0,0,'方伟业',1,1676686055,1676686055),(16,'修改部门','','/api/system/team','','A','system:team:update','PUT','','',14,0,0,0,'方伟业',1,1676686086,1676686086),(17,'删除部门','','/api/system/team','','A','system:team:delete','DELETE','','',14,0,0,0,'方伟业',1,1676686120,1676686120),(18,'角色管理','safe','role','sysRole','M','','','system/role/index','',2,0,0,0,'方伟业',1,1676686294,1676691447),(19,'查看角色','','/api/system/roles','','A','system:role:query','GET','','',18,0,0,0,'方伟业',1,1676686334,1676686334),(20,'新增角色','','/api/system/role','','A','system:role:add','POST','','',18,0,0,0,'方伟业',1,1676686390,1676691345),(21,'修改角色','','/api/system/role','','A','system:role:update','PUT','','',18,0,0,0,'方伟业',1,1676686414,1676691354),(22,'删除角色','','/api/system/role','','A','system:role:delete','DELETE','','',18,0,0,0,'方伟业',1,1676686455,1676691361),(23,'用户管理','user','user','sysUser','M','','','system/user/index','',2,0,0,0,'方伟业',1,1676686506,1676687642),(24,'查看用户','','/api/system/users','','A','system:user:query','GET','','',23,0,0,0,'方伟业',1,1676686542,1676686542),(25,'新增用户','','/api/system/user','','A','system:user:add','POST','','',23,0,0,0,'方伟业',1,1676686578,1676686648),(26,'修改用户','','/api/system/user','','A','system:user:update','PUT','','',23,0,0,0,'方伟业',1,1676686603,1676686643),(27,'删除用户','','/api/system/user','','A','system:user:delete','DELETE','','',23,0,0,0,'方伟业',1,1676686637,1676686637),(28,'登陆日志','history','login_log','sysLoginLog','M','','','system/login_log/index','',2,0,0,0,'方伟业',1,1676686716,1676687647),(29,'查看登陆日志','','/api/system/login/log','','A','system:login:log:query','GET','','',28,0,0,0,'方伟业',1,1676686777,1676686777),(30,'修改角色菜单','','','','G','system:role:menu:update','','','',18,0,0,0,'方伟业',1,1676687115,1676687115),(31,'获取角色的菜单id','','/api/system/role/menu_ids','','A','system:role:menu:query','GET','','',30,0,0,0,'方伟业',1,1676687183,1676687183),(32,'修改角色菜单','','/api/system/role/menu','','A','system:role:menu:update','PUT','','',30,0,0,0,'方伟业',1,1676687240,1676687240),(33,'数据展板','dashboard','/dashboard','Dashboard','M','','','Layout','/dashboard/workplace',1,0,0,1,'方伟业',1,1676687974,1676688048),(34,'系统数据','dashboard','workplace','Workplace','M','','','dashboard/workplace/index','',33,0,1,0,'方伟业',1,1676688043,1676688107);


# 创建部门数据
INSERT INTO `tb_system_team` VALUES (1,'青岑云科技','青岑云科技',0,'system',0,1673682377,1673682377);

# 创建角色数据
INSERT INTO `tb_system_role` VALUES (1,0,'超级管理员','superAdmin',1,1,'超级管理员','ALLTEAM',NULL,'system',0,1676619290,1676619290);

# 创建用户数据
INSERT INTO `tb_system_user` VALUES (1,1,1,'超级管理员','方伟业','18288888888','','128@qq.com',1,'$2a$10$ogNytaIiyFwlRmfz.Xj1kO9vBRohNk.G6qrVsi4P7Lk7vV6KFqpAK',1,'',1677341092,'system',0,1676626494,1677341092);

