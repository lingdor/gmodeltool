# gmodeltool
the tools of gmodel, you can automatically  generate codes from go:generate, or commands. If your data table fields chanaged, you only need to rerun the command.

# Installaion

```shell
    go install github.com/lingdor/gmodeltool@latest
```

# Configuation

write configuration to your project root (gmodel.yml):
```yaml
gmodel:
  connection:
    default:
      dsn: mysql://localhost@bobby:abc123/dbname
    user:
      dsn: mysql://localhost@bobby:abc123/userdb
```
# Generate gmodel schema code to a code file:
```go
//go:generate gomodeltool gen schema --tables "tb_user"
//or
//go:generate gomodeltool gen schema --tables "tb_%"
```

# Generate entity from database table

```go
//go:generate gmodeltool gen entity --tables "tb_user%"

//or generate a entity for gorm 
//go:generate gmodeltool gen entity --tables "tb_user" --gorm
```

# Generate codes in shell easily
```shell
    gmodeltool gen entity --tables "tb_%" --tofiles
    #or
    gmodeltool gen schema --tables "%" --tofiles
    #or
    gmodeltool gen schema --tables "tb_user,tb_school" --tofiles
    #or
    gmodeltool gen schema --tables "tb_user" --dry-run --verbose
    #or
    gmodeltool gen schema --tables "tb_user" --dsn "mysql://user:pass@tcp(localhost:3306)/db1" --tofiles
    
```
