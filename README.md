- お金の貸し借りを累積で覚えててくれるLINE botのバックエンド
- 頻発する夫婦間の貸し借り/立て替えをミニマルに管理したくて作った
- golang on heroku/Redis

## bot usage
- preparation: 
    - botをグループに招待
- post:
    - `[int]`: グループメンバーに払ってほしい額をpost
    - `iam [string]`: 自分の呼ばれ方をきめる/変える
    - `check`: 累積額の差を問い合わせる
    - `reset`: 実績を消す(精算)
    - `help`: help

## setting up server
- requires
    - active line bot channel
    - active heroku app with Redigo(add-on) backend

- set line tokens as env var

```
heroku config:set CHANNEL_SECRET=** -a [appname]
heroku config:set CHANNEL_ACCESS_TOKEN=** -a [appname]
```

- deploy heroku app
