#!/usr/bin/env bash
echo "Setting up user's rootpath $OWNCLOUD_ROOTPAH.." #>> modifikasi
sed -i.bak 's#);#'$OWNCLOUD_ROOTPAH');#' /var/www/owncloud/config/config.php #>> modifikasi
case "${OWNCLOUD_BACKGROUND_MODE}" in
  "ajax")
    echo "Enabling ajax background..."
    occ background:ajax
    ;;
  "cron")
    echo "Enabling cron background..."
    occ background:cron
    ;;
  "webcron")
    echo "Enabling webcron background..."
    occ background:webcron
    ;;
esac

if [[ ${OWNCLOUD_CROND_ENABLED} == "true" ]]
then
  echo "Touching cron configs..."
  touch /etc/cron.d/*

  echo "Starting cron daemon..."
  /usr/sbin/cron -l
fi

true
