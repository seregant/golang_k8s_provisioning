FROM owncloud/server:10.3
RUN rm /etc/owncloud.d/20-config.sh
COPY 20-config.sh /etc/owncloud.d/20-config.sh
RUN chmod +x /etc/owncloud.d/20-config.sh
ENV APACHE_LISTEN 80
EXPOSE 80

ENTRYPOINT ["/usr/bin/entrypoint"]
CMD ["/usr/bin/owncloud", "server"]
