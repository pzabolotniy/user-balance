# user-balance

# Build and run

```bash
docker-compose up --build -d --remove-orphans
```

# Stop
```bash
docker-compose down
```

* The application's api will listen the `:8080` port
* The solution contains only one test user - `de152cc3-9cbf-45c6-9081-7dff96708254`
* you'll receive `404 Not Found` if you send unknown user_id

# Test requests

* transaction ID should be passed in the URL
* transaction ID passed in payload will be ignored

```curl
curl -vvv -s \
    -H 'Source-Type: game' \
    -H 'Content-type: application/json' \
    -d '{"state": "win", "amount": "11"}' \
    http://localhost:8080/api/v1/users/de152cc3-9cbf-45c6-9081-7dff96708254/transactions/qwe1
```

* will receive `409 Conflict` if you send same transaction id more than once

# Transaction cancelation
* **cancelation.interval** is used to set up the cancelation launching
* **cancelation.txs_per_round** is used to set up the number of odd records to cancel
* you should restart (see Build and run section) the application to apply config changes

# linter
* you may run `golangci-lint` linter over the sources. linter settings are included
