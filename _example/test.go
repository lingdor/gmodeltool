package _example

// step1, install command: go install github.com/lingdor/gmodeltool
// step2, edit gmodel.yml, write the right db connection dsn
// step3, chanage the below annotation parameters, --tables to your tables.
//go:generate gmodeltool gen schema --tables t_user
//or
//go:generate gmodeltool gen entity --gorm --tables tb_user
