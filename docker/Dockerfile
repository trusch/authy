FROM alpine

ADD authy /bin/authy

CMD /bin/authy \
  --http-address :80 \
  --private-key /srv/authy/authy.key \
  --public-key /srv/authy/authy.crt \
  --storage xz+file:///srv/authy/storage
