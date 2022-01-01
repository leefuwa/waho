module appservice

go 1.15

replace waho => ../

require (
	github.com/bwmarrin/snowflake v0.3.0
	github.com/sirupsen/logrus v1.7.0
	github.com/streadway/amqp v1.0.0
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	waho v0.0.0-00010101000000-000000000000
)
