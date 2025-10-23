PROTO_DIRS := antigonos eukleides hefaistion parmenion perdikkas ptolemaios filippos

.PHONY: all
all: generate docs proto antigonos

.PHONY: generate
generate:
	@for dir in $(PROTO_DIRS); do \
		echo "Generating Protobuf files in $$dir..."; \
		buf generate --template $$dir/buf.gen.yaml $$dir; \
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

.PHONY: antigonos-dev
antigonos-dev:
	docker build -f antigonos/Dockerfile.dev -t ghcr.io/odysseia-greek/antigonos:dev .
	docker push ghcr.io/odysseia-greek/antigonos:dev

.PHONY: dev
dev: antigonos-dev
