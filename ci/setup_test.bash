#!/usr/bin/env bash
set -eux

if [[ "${PGDBVERSION-}" =~ ^[0-9.]+$ ]]
then
  sudo apt-get remove -y --purge postgresql libpq-dev libpq5 postgresql-client-common postgresql-common
  sudo rm -rf /var/lib/postgresql
  wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
  sudo sh -c "echo deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -cs)-pgdg main $PGDBVERSION >> /etc/apt/sources.list.d/postgresql.list"
  sudo apt-get update -qq
  sudo apt-get -y -o Dpkg::Options::=--force-confdef -o Dpkg::Options::="--force-confnew" install postgresql-$PGDBVERSION postgresql-server-dev-$PGDBVERSION postgresql-contrib-$PGDBVERSION
  sudo chmod 777 /etc/postgresql/$PGDBVERSION/main/pg_hba.conf
  echo "local     all         postgres                          trust"    >  /etc/postgresql/$PGDBVERSION/main/pg_hba.conf
  echo "local     all         all                               trust"    >> /etc/postgresql/$PGDBVERSION/main/pg_hba.conf
  echo "host      all         pgdb_user     127.0.0.1/32        md5"      >> /etc/postgresql/$PGDBVERSION/main/pg_hba.conf
  sudo chmod 777 /etc/postgresql/$PGDBVERSION/main/postgresql.conf
  if $(dpkg --compare-versions $PGDBVERSION ge 9.6) ; then
    echo "wal_level='logical'"     >> /etc/postgresql/$PGDBVERSION/main/postgresql.conf
    echo "max_wal_senders=5"       >> /etc/postgresql/$PGDBVERSION/main/postgresql.conf
    echo "max_replication_slots=5" >> /etc/postgresql/$PGDBVERSION/main/postgresql.conf
  fi
  sudo /etc/init.d/postgresql restart

  psql -U postgres -c 'create database orm_test'
  psql -U postgres orm_test -c 'create extension hstore'
  psql -U postgres orm_test -c 'create domain uint64 as numeric(20,0)'
  psql -U postgres -c "create user pgdb_user SUPERUSER PASSWORD 'secret'"
  psql -U postgres -c "create user `whoami`"
fi

if [[ "${PGDBVERSION-}" =~ ^cockroach ]]
then
  wget -qO- https://binaries.cockroachdb.com/cockroach-v20.2.5.linux-amd64.tgz | tar xvz
  sudo mv cockroach-v20.2.5.linux-amd64/cockroach /usr/local/bin/
  cockroach start-single-node --insecure --background --listen-addr=localhost
  cockroach sql --insecure -e 'create database orm_test'
fi

if [ "${CRATEVERSION-}" != "" ]
then
  docker run \
    -p "5432:5432" \
    -d \
    crate:"$CRATEVERSION" \
    crate \
      -Cnetwork.host=0.0.0.0 \
      -Ctransport.host=localhost \
      -Clicense.enterprise=false
fi
