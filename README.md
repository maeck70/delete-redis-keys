## delete-redis-keys: A Go tool for managing Redis keys

This Go application provides functionalities for managing keys in a Redis server. It allows you to:

* **Connect to a Redis server** based on provided hostname, port, and password (optional).
* **Create test keys** with sample data and varying expiration times.
* **Delete all keys** regardless of expiration.
* **Delete expired keys** based on a specified expiration duration.
* **Print all keys** and their associated data and update timestamps.

### Usage

This application uses command line flags to control its behavior. 

**Flags:**

* `-h`: Redis server hostname (default: `localhost`)
* `-p`: Redis server port (default: `6379`)
* `-e`: Expiration duration for keys in minutes (default: `15`) used for `-de` flag.
* `-da`: Delete all keys (default: `false`)
* `-de`: Delete expired keys (default: `false`)
* `-c`: Create test keys (default: `false`)
* `-o`: Output all keys with their data and update timestamps (default: `false`)

**Example:**

To connect to a Redis server on `my-redis-server` at port `6380` and delete all keys:

```bash
go run main.go -h my-redis-server -p 6380 -da
```

**Important Note:** Using `-da` will permanently delete all keys, use it with caution.

### Run the program:

   ```bash
   go run main.go
   ```

### Contributing

Feel free to submit pull requests for bug fixes or enhancements. 

### License

This project is licensed under the MIT License. See the LICENSE file for details.
