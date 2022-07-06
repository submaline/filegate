aws lambdaがよくわからず、学習自体が時間の無駄だったのでgoでなんとかした。  
ec2?のちっさいやつ750時間/月（登録から１２ヶ月間）無料らしいので。

実際の運用はこれのアップロードごとに認証つけて、s3にアップロードしたのち、lambdaを使って加工するのが良いのだろうか。

```shell
docker build -t <NAME> .
```

```shell
docker run \
  -v ~/.aws/config:/root/.aws/config:ro \
  -v ~/.aws/credentials:/root/.aws/credentials:ro \
  -v ~/**/<ThisProjectRoot>/.env:/go/src/app/.env:ro \
  -p 80:8080 <NAME>:latest
```