# 筋トレ日記アプリ インフラ設計書

## 1. 概要

本ドキュメントは、Flutter / Go / PostgreSQL を使用した「筋トレ日記アプリ」の本番運用を想定した
GCP 上のインフラ構成についてまとめたものである。

- フロントエンド：Flutter（iOS アプリ） → App Store 配信
- バックエンド：Go
- データベース：PostgreSQL
- インフラ基盤：Google Cloud Platform

---

## 2. 要件

- Flutter アプリから HTTPS 経由で API にアクセスできること
- 認証済みユーザーのみが記録データにアクセスできること(JWT を用いた認証を Go API 側で検証し、Flutter アプリから送信する想定)
- ログは Cloud Logging に集約し、アプリのログを確認できること
- CD を導入し、main ブランチにマージしたタイミングで自動デプロイが走るようにすること
- DB パスワードなどの秘匿情報は Secret Manager で管理する。

## 3. 採用する GCP サービスと役割

| サービス                   | 役割                                   |
| -------------------------- | -------------------------------------- |
| Cloud Run                  | Go バックエンド API の実行環境         |
| Artifact Registry          | バックエンドのコンテナイメージ格納     |
| Neon                       | アプリデータ永続化（外部データベース） |
| Secret Manager             | DB パスワード等の秘匿情報管理          |
| Cloud Logging / Monitoring | アプリ・インフラのログ・メトリクス管理 |

---

## 4. アーキテクチャ概要

<img width="990" height="416" alt="image" src="https://github.com/user-attachments/assets/cc0da8f0-56d6-4432-8ab0-bd3a2b1decac" />

---

## 5. デプロイフロー

### 5-1. 手動デプロイ（初期）

1. Go アプリの Docker イメージビルド
2. Artifact Registry へ push
3. `gcloud run deploy` コマンドで Cloud Run サービスを更新
4. 動作確認（ヘルスチェック用エンドポイント / 本番アプリからのアクセス）

### 5-2. 将来の自動デプロイ（構想）

- GitHub Actions などを用いて、
  - `main` ブランチにマージされたら自動で
    - Docker イメージビルド
    - Artifact Registry へ push
    - Cloud Run へデプロイ
- 上記は、初期構築後に段階的に導入する。

---

## 6. 運用・監視

- Cloud Logging
  - Cloud Run の標準出力 / エラー出力を Cloud Logging に集約。
- Cloud Monitoring
  - リクエスト数・エラー率・レイテンシを確認。
