FROM alpine:latest
EXPOSE 28080 28082
WORKDIR /app
COPY / ./monero-wallet-rpc
RUN chmod +x ./monero-wallet-rpc
RUN ./monero-wallet-rpc
# ENTRYPOINT ["./monero-wallet-rpc", "--stagenet", "--rpc-bind-port 28082", "--wallet-dir monero/wallets/main/"]