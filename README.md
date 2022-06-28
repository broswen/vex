# Vex

![diagram](vex.png)

### Todo
- [x] handle postgres errors and wrap in custom store errors (ongoing)
- [ ] provision account tokens to cloudflare kv
- [x] handle docker-compose initialize local postgres with schema
  - still need to handle multiple migrations
- [x] handle local provisioning for dockerfile, flag to skip api calls?
- [ ] add mocks and tests with testify
- [ ] add created_on and modified_on fields to all resources
- [ ] incremental config builds
  - store prerendered config in postgres, parse and insert/update flags as needed
