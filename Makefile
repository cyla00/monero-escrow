main:
	@mkdir -p monero/wallets/main/
	@templ generate
	@go run ./cmd/main.go

pod up:
	@podman compose up

pod down:
	@podman compose down

templ:
	@templ generate

tailwind:
	@./tailwindcss -i input.css -o output.css --watch

stagenet create:
	@./monero-wallet-rpc --stagenet --rpc-bind-port 28082 --wallet-dir monero/wallets/main/ --password walletpassword --rpc-login xmruser:xmrpassword --log-file monero/logs/monero-wallet-rpc.log --max-log-files 2 --trusted-daemon --daemon-address http://node.monerodevs.org:38089 --non-interactive

stagenet run:
	@./monero-wallet-rpc --stagenet --rpc-bind-port 28082 --wallet-file monero/wallets/main/fidexmrwallet --password walletpassword --rpc-login xmruser:xmrpassword --log-file monero/logs/monero-wallet-rpc.log --max-log-files 2 --trusted-daemon --daemon-address http://node.monerodevs.org:38089
