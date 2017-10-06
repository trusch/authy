FROM alpine

ADD authy /bin/authy

CMD /bin/authy \
  --http-address :80 \
  --keyfile /srv/authy/key.pem \
  --webroot /srv/authy/webroot \
  --storage xz+file:///srv/authy/storage
