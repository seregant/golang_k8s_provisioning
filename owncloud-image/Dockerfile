FROM owncloud/server:10.3
RUN rm /etc/owncloud.d/55-cron.sh
COPY 55-cron.sh /etc/owncloud.d/55-cron.sh
RUN chmod +x /etc/owncloud.d/55-cron.sh
ENV APACHE_LISTEN 80
EXPOSE 80

ENTRYPOINT ["/usr/bin/entrypoint"]
CMD ["/usr/bin/owncloud", "server"]
