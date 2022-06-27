# Vex

![diagram](vex.png)

### Todo
- [x] handle postgres errors and wrap in custom store errors (ongoing)
- [ ] provision account tokens to cloudflare kv
- [ ] handle docker-compose initialize local postgres with schema
- [ ] handle local provisioning for dockerfile, flag to skip api calls?
- [ ] add a ton of tests
- [ ] incremental config builds
  - store prerendered config in postgres, parse and insert/update flags as needed