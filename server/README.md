# Chuyên đề Backend
### Setup your local dev environment
After checkout this repo. You can go through the following steps to get stared:
1. Setup `direnv` for managing your local environment varibles\
Mac: `brew install direnv`\
Linux: `sudo apt -y install direnv`\
Hook: https://direnv.net/docs/hook.html
2. Setup docker compose for running dependencies. In our project, we need mysql in order to run the app.
`curl -fsSL https://get.docker.com | bash`
3. Start docker compose
`docker-compose -f build/docker-compose.yaml up`
4. Run db migration
Install: https://atlasgo.io/getting-started/
Mac + Linux: `curl -sSf https://atlasgo.sh | sh` \
Run `make migration-up`
5. Start app (already setup for vscode -> press F5)

### How to migrate database
We use [Entgo](https://entgo.io/) as our ORM lib and [Atlas](https://atlasgo.io/getting-started#installation) for database migration.
1. If you want to create new migration, please run this command `$ make migrate name=<MIGRATION_NAME>`
2. To apply the migrations: `$ make migrate-up`
