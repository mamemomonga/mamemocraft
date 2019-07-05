# mamemocraft

まめもくらふと ランチャー・ビルダー

* [起動](bin/run.sh)
* [停止](bin/stop.sh)
* [メンテナンスモード開始](bin/maintenance.sh)
* [起動完了監視](bin/launch_watcher.sh)
* [mcrcon Linuxバイナリ](bin/mcrcon)

# メモ

GCE Preemptive Instance の終了通知は、インスタンスのカスタムメタデータとして、**キー: shutdown-script 値: スクリプト** を与えればよい。ただし、許容されたシャットダウン時間は10秒である。

# LICENSE

MIT

except [mcrcon](https://github.com/Tiiffi/mcrcon) ( zlib License )

