module github.com/dimsonson/pswmanager

go 1.20

// replace log => github.com/rs/zerolog v1.29.0

require (
	github.com/MashinIvan/rabbitmq v0.1.0
	github.com/derailed/tcell/v2 v2.3.1-rc.3
	github.com/derailed/tview v0.8.1
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2 v2.0.0-rc.3
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.0-rc.3
	github.com/jackc/pgx/v5 v5.3.1
	github.com/mutecomm/go-sqlcipher v0.0.0-20190227152316-55dbde17881f
	github.com/redis/go-redis/v9 v9.0.3
	github.com/rs/zerolog v1.29.0
	github.com/streadway/amqp v1.0.0
	google.golang.org/grpc v1.54.0
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/rivo/uniseg v0.4.3 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/term v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230331144136-dcfb400f0633 // indirect
// github.com/rs/zerolog v1.29.0
)

replace github.com/MashinIvan/rabbitmq => github.com/dimsonson/rabbitmq v0.0.0-20230521180522-7f3c8f0d2f7e

//replace github.com/MashinIvan/rabbitmq => github.com/dimsonson/rabbitmq latest
