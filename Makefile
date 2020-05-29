PROTO_DIRS= "testdata" "testdata/initdb" 
TESTS= "mapper" "plugin" "templates"

gen:
	go install .
	for i in $(PROTO_DIRS); do \
		protoc --map_out="sql=$$i/sql:$(GOPATH)/src"   \
			--go_out="plugins=grpc:$(GOPATH)/src"   \
			-I=. \
			$$i/*.proto ; \
	done

test:
	# runs all tests
	bazel test $(foreach var,$(TESTS), //$(var):go_default_test ) --verbose_failures --test_output=all

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
