CREATE USER 'sync'@'%' IDENTIFIED BY 'syncpass';
GRANT REPLICATION SLAVE,REPLICATION CLIENT ON *.* TO 'sync'@'%';
FLUSH PRIVILEGES;
CHANGE MASTER TO MASTER_HOST='master02',MASTER_PORT=3306,MASTER_USER='sync',MASTER_PASSWORD='syncpass';
/* START SLAVE; */
SHOW MASTER STATUS;
SHOW SLAVE STATUS;