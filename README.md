# Scalable e-commerce
This project comprises a few microservices, each with it's own codebase, connected by gRPC protocols.

# Technology stack
 - gRPC
 - sqlite3

# Used packages
 - `cleanenv` for reading config
 - `slog` for logging
 - `golang-migrate` for database migrations 
 - `gofakeit` for creating random passwords and emails for testing

# TO DO:
- [ ] Make authotization service
- [ ] Make product catalog service
- [ ] Make shopping cart service
- [ ] Make order service
- [ ] Make payment service
- [ ] Make notification service
- [ ] Add external handle and router
- [ ] Containerize
- [ ] Build pipelines with github jobs
- [ ] Add priviledges handling in auth service
