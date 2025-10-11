# Setup

## Database
Db name for inventory-service: inventories
Db name for order-service: orders

Please specify your own *user* and *password* for connection string.
To migrate inventory db:
```
cd inventories-service/
migrate -database postgres://postgres:postgres@localhost:5432/inventories?sslmode=disable -path sql/migrations up
```

To migrate order db:
```
cd inventories-service/
migrate -database postgres://postgres:postgres@localhost:5432/orders?sslmode=disable -path sql/migrations up
```

![inventory db](https://github.com/rianprayoga/synp-challenge/doc/main/db-inventory.png?raw=true)
![order db](https://github.com/rianprayoga/synp-challenge/doc/main/db-order.png?raw=true)


## Running service
By default inventory-service will run on port 8081 for rest and 8082 for grpc.
The order-service will run on port 8082.
You can specify the port by using the available flag. The same applied for db connection to run the app.

From root directory:
```
make run-inventories
```
or 
```
make run-orders
```

# Description
The inventory only consist of 2 tables, items and transactions. For simplicity no SKU to track number of items.
Using Optimistic Locking in inventory service, no resource blocking but increase the number of conflict between clients. 
To lock I use *version* column. This version will increase everytimes changes occur on the specified record.

The transactions table is to record the changes, if there is failed order we can revert it by orderId.

