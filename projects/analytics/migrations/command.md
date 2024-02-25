# Migration command

## Create Table

docker exec -it bagaig-db psql -U postgres -c "$(<migrations/001_create_table.sql)"

## Select from events table

docker exec -it bagaig-db psql -U postgres -c "SELECT \* FROM events;"
