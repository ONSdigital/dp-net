module github.com/ONSdigital/dp-net/v2

go 1.19

replace "github.com/ONSdigital/dp-renderer/v2" => "/Users/jon/dev/dp-renderer"

retract (
	v2.7.2 // contains retraction
	v2.7.1 // TODO: rethink application of timeout to standard  http server
	v2.7.0
)

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.252.0
	github.com/ONSdigital/dp-cookies v0.4.0
	github.com/ONSdigital/dp-renderer/v2 v2.4.0
	github.com/ONSdigital/log.go/v2 v2.4.1
	github.com/aws/aws-sdk-go v1.44.76
	github.com/gorilla/mux v1.8.0
	github.com/justinas/alice v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/smartystreets/goconvey v1.8.1
	github.com/stretchr/testify v1.8.0
	golang.org/x/net v0.14.0
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/ONSdigital/dp-healthcheck v1.6.0 // indirect
	github.com/c2h5oh/datasize v0.0.0-20220606134207-859f65c6625b // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/gosimple/slug v1.13.1 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/nicksnyder/go-i18n/v2 v2.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/smarty/assertions v1.15.1 // indirect
	github.com/unrolled/render v1.6.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
