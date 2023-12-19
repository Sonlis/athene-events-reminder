.PHONY: test
test: compose gotest clean

.PHONY: gotest
gotest:
	go test ./...

.PHONY: compose
compose:
	docker-compose up -d
	sleep 3
	PGPASSWORD=test_pass psql -h localhost -p 5447 -U test_user -d test -c "CREATE TABLE IF NOT EXISTS reminder(chat_id int, event_id int, reminder_time timestamptz);"

.PHONY: clean
clean:
	docker-compose down
