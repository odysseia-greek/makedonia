PROTO_DIRS := antigonos eukleides hefaistion parmenion perdikkas

.PHONY: all
all: generate docs proto

.PHONY: generate
generate:
	@for dir in $(PROTO_DIRS); do \
		echo "Generating Protobuf files in $$dir..."; \
		(cd $$dir && \
		 protoc --go_out=. --go_opt=paths=source_relative \
		        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		        proto/$$dir.proto); \
	done

.PHONY: docs
docs:
	@for dir in $(PROTO_DIRS); do \
		echo "Generating docs in $$dir..."; \
		docker run --rm \
			-v $$PWD/$$dir/docs:/out \
			-v $$PWD/$$dir/proto:/protos \
			localproto:latest --doc_opt=html,docs.html; \
		docker run --rm \
			-v $$PWD/$$dir/docs:/out \
			-v $$PWD/$$dir/proto:/protos \
			localproto:latest --doc_opt=markdown,docs.md; \
	done
	cd ./alexandros/docs && spectaql -c  spectaql.yaml

.PHONY: proto
proto:
	buf generate --template filippos/buf.gen.yaml filippos
	buf generate --template antigonos/buf.gen.yaml antigonos
