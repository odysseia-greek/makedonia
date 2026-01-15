PROTO_DIRS := antigonos eukleides hefaistion parmenion perdikkas ptolemaios filippos

.PHONY: all
all: generate docs antigonos

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

.PHONY: tidy tidy-go
tidy: tidy-go

tidy-go:
	@echo "==> Running 'go mod tidy' in all modules..."
	@mods=$$(find . -name go.mod -print0 | xargs -0 -n1 dirname | sort -u); \
	if [[ -z "$$mods" ]]; then \
		echo "No go.mod files found."; \
		exit 0; \
	fi; \
	for d in $$mods; do \
		echo "==> go mod tidy in $$d"; \
		( cd "$$d" && go mod tidy && go fmt ./... ); \
	done

# Convenience: list all module dirs found
.PHONY: mods
mods:
	@find . -name go.mod -print0 | xargs -0 -n1 dirname | sort -u


.PHONY: images-dev images-prod

images-dev:
	OWNER="$(OWNER)" ROOT="$(ROOT)" GROUP=dev \
		./bump-images.sh

images-prod:
	OWNER="$(OWNER)" ROOT="$(ROOT)" GROUP=prod \
		./bump-images.sh