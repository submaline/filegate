環境変数
- R2_BUCKET_NAME
- R2_ACCOUNT_ID
- R2_ACCESS_KEY_ID
- R2_ACCESS_KEY_SECRET

cloudflare R2にファイルをアップロードするためのサーバ。

```shell
docker build -t <NAME> .
```

```shell
docker run \
  -v ~/<ThisProjectRoot>/.env:/go/src/app/.env:ro \
  -p 80:8080 <NAME>:latest
```