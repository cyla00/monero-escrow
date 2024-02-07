FROM alpine
EXPOSE 28080 28082
WORKDIR /app
COPY / ./monero-wallet-rpc
RUN chmod +x ./monero-wallet-rpc
ENTRYPOINT ["monero-wallet-rpc", "--rpc-bind-port 28082", "--stagenet", "password walletpassword", "--rpc-login xmruser:xmrpassword"]