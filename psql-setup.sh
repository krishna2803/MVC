#!/bin/ash

export $(grep -v '^#' .env | xargs)

mkdir /run/postgresql && chown postgres: /run/postgresql
chown -R postgres /var/lib/postgresql/data
su-exec postgres initdb -D /var/lib/postgresql/data

su-exec postgres pg_ctl start -D /var/lib/postgresql/data -l /var/lib/postgresql/logfile.log
su-exec postgres createuser $DB_USER
su-exec postgres createdb $DB_NAME
su-exec postgres psql -c "ALTER USER $DB_USER WITH ENCRYPTED PASSWORD '$DB_PASS';"
su-exec postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME to $DB_USER;"
su-exec postgres pg_ctl stop -D /var/lib/postgresql/data -m smart
