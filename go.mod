module github.com/TimothyYe/godns

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/bwmarrin/discordgo v0.25.0
	github.com/fatih/color v1.13.0
	github.com/google/uuid v1.3.0
	github.com/json-iterator/go v1.1.12
	github.com/linode/linodego v1.6.0
	github.com/miekg/dns v1.1.50
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e
	golang.org/x/oauth2 v0.0.0-20220622183110-fd043fe589d2
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/go-resty/resty/v2 v2.7.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/sys v0.0.0-20220624220833-87e55d714810 // indirect
	golang.org/x/tools v0.1.11 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)
go 1.17
replace github.com/phiwatec/godns/internal/provider/strato => ./internal/provider/strato

