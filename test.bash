go run ./main.go \
    -package=csicinder \
    -go-header-file=/home/wozy/projects/go/src/github.com/kubermatic/kubermatic/hack/boilerplate/ce/boilerplate.go.txt \
    -src=./examples/cinder-csi-resources.yaml > ./output/test_resources.go

go run ./main.go \
    -package=csicinder \
    -go-header-file=/home/wozy/projects/go/src/github.com/kubermatic/kubermatic/hack/boilerplate/ce/boilerplate.go.txt \
    -kubermatic-interfaces=true \
    -src=./examples/cinder-csi-resources.yaml > ./output/test_interfaces.go

go run ./main.go \
    -package=csicinder \
    -go-header-file=/home/wozy/projects/go/src/github.com/kubermatic/kubermatic/hack/boilerplate/ce/boilerplate.go.txt \
    -src=./examples/ws-crd.yaml > ./output/test_crd.go
