ifeq (,$(wildcard ./config/app.env))
  $(error "Create a app.env file based on app.env.example")
endif
include ./config/app.env
export

POSTGRES_USER := $(POSTGRES_USER)
POSTGRES_PASSWORD := $(POSTGRES_PASSWORD)
DB_PORT := $(DB_PORT)
DB_NAME := $(DB_NAME)
SSL_ENABLE := $(SSL_ENABLE)
DB_URI := $(DB_URI)


.PHONY: postgres
postgres:
	docker run \
	--name postgres \
	-p $(DB_PORT):$(DB_PORT) -e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -d postgres

.PHONY: createdb
createdb:
	docker exec -it postgres createdb \
	--username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(DB_NAME)

.PHONY: dropdb
dropdb:
	docker exec -it postgres dropdb $(DB_NAME)

.PHONY: migrateup
migrateup: 
	migrate -path db/migrate \
	-database postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_URI):$(DB_PORT)/$(DB_NAME)$(SSL_ENABLE) \
	-verbose up

.PHONY: migratedown
migratedown:
	migrate -path db/migrate \
	-database postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_URI):$(DB_PORT)/$(DB_NAME)$(SSL_ENABLE) \
	-verbose down

.PHONY: migrateup1
migrateup1: 
	migrate -path db/migrate \
	-database postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_URI):$(DB_PORT)/$(DB_NAME)$(SSL_ENABLE) \
	-verbose up 1

.PHONY: migratedown1
migratedown1:
	migrate -path db/migrate \
	-database postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_URI):$(DB_PORT)/$(DB_NAME)$(SSL_ENABLE) \
	-verbose down 1

.PHONY: full_restart
full_restart:
	@read -p "Are you sure? This will delete all data. Type 'yes' to continue: " confirm && [ "$$confirm" = "yes" ] && \
	$(MAKE) dropdb && \
	$(MAKE) createdb && \
	$(MAKE) migrateup || \
	echo "Operation cancelled."

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: test
test:
	@if command -v gotestsum > /dev/null; then \
		gotestsum --debug --format testname; \
	else \
		go test ./...; \
	fi
    
.PHONY: testv
testv:
	@if command -v gotestsum > /dev/null; then \
		gotestsum --debug --format standard-verbose; \
	else \
		go test -v ./...; \
	fi

.PHONY: gen-cert
gen-cert:
	@mkdir -p config/certs
	@echo "Generating self-signed TLS certificate (RSA 2048, 365 days)" 
	@openssl req -x509 -newkey rsa:2048 -nodes \
	 -keyout config/certs/server.key \
	 -out config/certs/server.crt \
	 -days 365 \
	 -subj "/C=US/ST=Local/L=Local/O=Dev/OU=Dev/CN=localhost"
	@echo "Done: config/certs/server.crt and server.key created"

.PHONE: mockdb
mockdb:
	mockgen -package mockdb \
	-destination db/mock/store.go github.com/mauzec/user-api/db/sqlc Store