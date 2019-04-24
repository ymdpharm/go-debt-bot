## usage
- line グループに招待して半角整数値を送りつけるだけで、貸し借りの累積差し引き額を覚えててくれる
- special commands: 
    - `iam [string]`: 呼び名きめる/変える
    - `check`: 今の差し引き額を聞く
    - `reset`: 差し引きゼロにする
    - `help`: 今使えるコマンド

## requires
- active line bot channel
- heroku environment with Redis(add-on) backend
```
heroku config:set CHANNEL_SECRET=[line channel secret] -a [appname]
heroku config:set CHANNEL_ACCESS_TOKEN=[line channel access token] -a [appname]
```