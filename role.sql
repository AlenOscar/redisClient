
CREATE DATABASE IF NOT EXISTS testdb;

USE testdb;

CREATE TABLE `role` (
  `uid` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(20) DEFAULT '',
  `age` int(11) DEFAULT '0',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8
