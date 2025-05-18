# Blaze

## How to run
1. Make sure you have Go installed. You can download it from [here](https://go.dev/dl/).
2. Clone the repository
3. Take a look at `config.yaml` and make necessary changes
4. Run
   ```bash
   go run ./
   ```

## Generating SSL certificates for local HTTPS

The frontend server needs `wss` to work. To do this, you need to generate SSL certificates. 
You can use the following command to generate self-signed certificates:

```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=127.0.0.1"
```

To generate certificates for a specific network (e.g., Wi-Fi address 192.168.0.104), replace CN with the desired address.
Replace the same address in the `config.yaml` file.

If the websocket is not connecting, visit the URL in your browser and accept the certificate. 
Since the certificate is self-signed, the browser will show a warning. You can ignore it and proceed to the site.
