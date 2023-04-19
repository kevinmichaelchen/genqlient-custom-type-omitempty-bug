FLAGS = --skip-update-check --project hasura --admin-secret myadminsecretkey
MIGRATE_FLAGS = --database-name postgres
NPX = bunx

.PHONY: start
start:
	docker compose up -d

.PHONY: stop
stop:
	docker compose stop

.PHONY: console
console:
	hasura $(FLAGS) console

.PHONY: migrate
migrate:
	hasura $(FLAGS) $(MIGRATE_FLAGS) \
	  migrate apply

.PHONY: status
status:
	hasura $(FLAGS) $(MIGRATE_FLAGS) \
	  migrate status

.PHONY: md-apply
md-apply:
	hasura $(FLAGS) \
	  metadata apply

.PHONY: md-reload
md-reload:
	hasura $(FLAGS) \
	  metadata reload

## pull-schema      : introspects and pulls Hasura's schema (requires a locally running Hasura instance)
.PHONY: pull-schema
pull-schema:
	$(NPX) graphqurl http://localhost:8080/v1/graphql -H 'x-hasura-admin-secret: myadminsecretkey' --introspect > hasura.sdl.graphqls

## format                 : Format code
.PHONY: format
format:
	$(NPX) prettier --write "internal/**/*.graphql"

## gen-client             : auto-generates GraphQL client code
.PHONY: gen-client
gen-client: format
	(cd internal/client/graphql && go run github.com/Khan/genqlient genqlient.yaml)
