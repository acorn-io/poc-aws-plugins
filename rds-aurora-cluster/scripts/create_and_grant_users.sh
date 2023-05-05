#!/bin/bash

set -e

# Checking if it's a networking thing
sleep 10
echo "HOST: ${MYSQL_HOST}"

sql_file="/acorn/create_grant.sql"
user_dir="/acorn/users"

# Reset the file
echo "use ${MYSQL_DATABASE};" > ${sql_file}

for u in $(ls ${user_dir}); do
echo "CREATE USER IF NOT EXISTS '${u}'@'%' IDENTIFIED BY '$(<${user_dir}/${u})';">> ${sql_file}
echo "GRANT ALL PRIVILEGES on ${MYSQL_DATABASE}.* TO '$u'@'%';">> @{sql_file}
done

mysql -u"${MYSQL_ADMIN_USER}" -p"${MYSQL_PASSWORD}" -h"${MYSQL_HOST}" < ${sql_file}