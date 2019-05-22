
CREATE DATABASE IF NOT EXISTS testdb;

USE testdb;

CREATE TABLE IF NOT EXISTS `role` (
  `uid` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `username` varchar(20) DEFAULT '' COMMENT '名称',
  `nickname` varchar(20) DEFAULT '' COMMENT '昵称',
  `level` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '等级',
  `registerDatetime` datetime COMMENT '注册时间',
  `lastLoginDatetime` datetime COMMENT '最后一次登录时间',
  `data` blob COMMENT '二进制数据',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT '角色表';


-- DROP TABLE IF EXISTS mail;
CREATE TABLE IF NOT EXISTS mail (
  `gid` BIGINT UNSIGNED NOT NULL UNIQUE COMMENT '',
  `rolegid` BIGINT UNSIGNED NOT NULL COMMENT '',
  `data` MEDIUMBLOB NOT NULL COMMENT '',
  INDEX(`rolegid`)
) ENGINE = InnoDB CHARACTER SET = utf8 COMMENT '';

-- 将语句的结束符号从分号;临时改为两个$$
DELIMITER $$

-- 创建存储过程
DROP PROCEDURE IF EXISTS add_mail $$
CREATE PROCEDURE add_mail (IN param1 BIGINT UNSIGNED, IN param2 BIGINT UNSIGNED, IN param3 MEDIUMBLOB)
BEGIN
  INSERT INTO mail(`gid`, `rolegid`, `data`) VALUES(param1, param2, param3);
END $$

DROP PROCEDURE IF EXISTS get_mail $$
CREATE PROCEDURE get_mail (IN param1 INT, IN param2 INT)
BEGIN
  SELECT `gid`, `rolegid`, `data` FROM mail LIMIT param1, param2;
END $$

DELIMITER ;


