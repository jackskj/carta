PROTO_DIRS= "testdata/initdb" 
DBS= "mapper" "plugin" "templates"



testdbs:
	docker run --name carta-postgres-test --env POSTGRES_HOST_AUTH_METHOD=trust -d  -p 5432:5432 postgres
	docker run --name carta-mysql-test -d  --env MYSQL_ALLOW_EMPTY_PASSWORD=yes --env MYSQL_DATABASE=mysql -p 3306:3306 mysql

gen:
	go install .
	for i in $(PROTO_DIRS); do \
		protoc 	--go_out="plugins=grpc:$(GOPATH)/src"   \
			-I=. \
			$$i/*.proto ; \
	done


install:
	# generating map binary in $$GOPATH/bin
	go install .

build:
	bazel build //:protoc-gen-map

gazelle:
	# generates build files
	bazel run //:gazelle_update

repos:
	# generates go repos
	bazel run //:gazelle_update -- update-repos -from_file=go.mod -to_macro=bazel/go_repositories.bzl%go_repositories

fix:
	# fixes deprecated usage of rules
	bazel run //:gazelle_update -- fix

.PHONY: gen test fix repos gazelle build install test
