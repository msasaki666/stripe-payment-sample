#!/usr/bin/expect

# タイムアウト時間を設定
set timeout 20

# stripe loginコマンドを実行
spawn stripe login --interactive

# プロンプトを待つ
expect "Enter your API key:"
# 認証コードを入力
send "$env(STRIPE_API_KEY)\r"
expect "How would you like to identify this device in the Stripe Dashboard?"
send "\r"

# 処理が完了するまで待つ
expect eof
