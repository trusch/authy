authy
=====

authy is a HTTP service which provides signed JSON Web Token to its users

## Functionality

* Manage Users (Create, Retrieve, Update, Delete, List)
* Manage Auth Schemas (Create, Retrieve, Update, Delete, List)
* proper patch requests for users and schemas
* login and request some claims
* verify a given token
* arbitary additional user properties
* arbitary claim structure

## Idea

* Every user has `id`, a `password` and one `authschema`
  * Other fields are possible but not manadatory
  * `id` is autogenerated if not supplied
  * `authshema` is the id of a known auth schema
* On `put`, `post` or `patch` of a user the password field is hashed using bcrypt
* Auth-Schemas are JSON-Schemas which represents what claims a user is able to aquire
* On login `id` and `password` are checked and `claims` are validated against the `authschema` of the user
* If everthing is ok, the `claims` are signed and packed as a JSON Web Token and send to the user
* The target service can now verify the token by itself using the public key of authy or send the token to authy to get the contained claims back if the token is valid.

## Setup

### Create Keys
```
go get -v github.com/trusch/pki/pkitool
pkitool -c P521 -r 0 init
pkitool -c P521 -r 0 issue server authy
```

### Start authy
```
docker run --rm -d \
  --name authy \
  --volume $(pwd)/pki/authy.key:/srv/authy/authy.key \
  --volume $(pwd)/pki/authy.crt:/srv/authy/authy.crt \
  --volume /srv/authy/storage \
  -p 80:80 \
  trusch/authy
```

### Interact
```
# create admin schema (admin can aquire role admin|user|guest)
curl -XPUT -d'{"properties":{"role":{"pattern":"^admin$|^user$|^guest$"}}}' http://localhost/api/v1/schema/admin

# create admin user
curl -XPOST -d'{"id":"admin","password":"admin","authschema":"admin"}' http://localhost/api/v1/user/

# login
token=$(curl -XPOST -d '{"id":"admin","password":"admin","claims":{"role":"admin"}}' http://localhost/api/v1/login)

# verify
curl -XPOST -d"$token" http://localhost/api/v1/verify
-> {"role":"admin"}
```
