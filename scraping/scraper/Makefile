include .env
export

scrape:
	go build -o main main.go
	./main

clear_cache:
	rm $(TMP_DIR)/vlr_cache.db

