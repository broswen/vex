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

### Terraform Provider
```hcl
provider vex {
  token = "api token"
  account_id = "account id"
}

resource "vex_account" "main" {
  name = "account name"
  description = "account description"
}

resource "vex_project" "app1" {
  account_id = vex_account.main.id
  name = "project name"
  description = "project description"
}

resource "vex_flag" "feature1" {
  project_id = vex_project.app1.id
  key = "flag key"
  type = "flag type"
  value = "flag raw value"
}

```