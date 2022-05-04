module github.com/digitalocean/openvswitch_exporter

go 1.15

require (
	github.com/digitalocean/go-openvswitch v0.0.0-20211105174344-2a0f99c6436b
	github.com/go-kit/kit v0.10.0 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/mdlayher/genetlink v1.2.0 // indirect
	github.com/mdlayher/socket v0.2.3 // indirect
	github.com/prometheus/client_golang v1.12.1
	github.com/prometheus/common v0.34.0 // indirect
	github.com/prometheus/prometheus v2.2.1-0.20180315085919-58e2a31db8de+incompatible
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/digitalocean/go-openvswitch => github.com/gecio/go-openvswitch v0.0.0-20220504095400-3dad490fcf3b
