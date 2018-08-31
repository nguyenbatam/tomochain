rm -Rf node*/tomo*
rm -Rf logs/*.txt
mongo governance --eval "db.dropDatabase();"
redis-cli flushall
