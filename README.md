- 夫婦間での立て替えた金額の累積の差を管理するだけのLINE bot
- 高機能なアプリを使うのがイヤだったので作った
- 複数グループ/複数人数に対応しているはずだけど細かい動作は確認してない
- LINE Messaging API/golang/heroku/Redis

## usage
- preparation: 
    - グループに招待して使う
- basic words:
    - `[int]`: "貸し" を登録
- special words: 
    - `iam [string]`: 自分の呼ばれ方きめる/変える
    - `check`: グループ内の貸し借りを確認する
    - `reset`: グループ内の貸し借りを精算する
    - `help`: help

## requires
- active line bot channel
- active heroku app with Redigo(add-on) backend
```
heroku create
heroku config:set CHANNEL_SECRET=[line channel secret] -a [appname]
heroku config:set CHANNEL_ACCESS_TOKEN=[line channel access token] -a [appname]
```
