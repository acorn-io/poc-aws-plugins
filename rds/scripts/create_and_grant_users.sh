#!/bin/bash
set -e

while ! mysql --connect-timeout=2 -h ${MYSQL_HOST} -u ${MYSQL_ADMIN_USER} -p"${MYSQL_ADMIN_PASSWORD}" -e "SELECT 1"; do
    sleep .5
done

sql_file="./create_grant.sql"

# Reset the file
echo "use ${MYSQL_DATABASE};" > ${sql_file}

echo "CREATE USER IF NOT EXISTS '${MYSQL_USER}'@'%' IDENTIFIED BY '${MYSQL_PASSWORD}';">> ${sql_file}
echo "GRANT ALL PRIVILEGES on ${MYSQL_DATABASE}.* TO '${MYSQL_USER}'@'%';">> @{sql_file}

mysql -u"${MYSQL_ADMIN_USER}" -p"${MYSQL_ADMIN_PASSWORD}" -h"${MYSQL_HOST}" < ${sql_file}