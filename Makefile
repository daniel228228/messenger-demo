include .env
export

restore_db_messenger:
	docker exec -i $(MESSENGER_DB_HOST) psql "postgres://${MESSENGER_DB_USER}:${MESSENGER_DB_PASSWORD}@localhost:${MESSENGER_DB_PORT}/${MESSENGER_DB_NAME}" < ./migrations/messenger/messenger.sql

restore_db_users:
	docker exec -i $(USERS_DB_HOST) psql "postgres://${USERS_DB_USER}:${USERS_DB_PASSWORD}@localhost:${USERS_DB_PORT}/${USERS_DB_NAME}" < ./migrations/users/users.sql