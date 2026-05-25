.PHONY: test

# Run all Go tests across workspace modules.
# NOTE: "go test ./..." is not supported in this go.work layout because
# orm-crud and sa-token are parent directories, not module roots.
test:
	go test ./admin/... ./go-common/... ./orm-crud/api/... ./orm-crud/gormc/... ./orm-crud/pagination/... ./sa-token/rueidis/...
