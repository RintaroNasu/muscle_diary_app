# 筋トレ日記アプリ

## アプリ概要

**筋トレを記録してみんなと共有してモチベーションを上げるアプリです。<br>
種目・重量・回数をサクっと入力し、カレンダーとグラフで振り返れます。**

### 機能概要

1. **筋トレの回数を記録機能**: 種目・重量・回数・日付を登録<br>
2. **カレンダー記録表示機能**: 記録がある日にドット表示し、日付をタップすると当日の記録一覧が出ます。<br>
3. **グラフ記録表示機能**: 期間を指定して回数・重量などの推移を折れ線／棒グラフで可視化します。<br>
4. **ランキング表示機能**: 月間のジム日数ランキングを可視化して、みんなのモチベーションの向上を計ります。<br>

---

## デモ画像

<table>
　 <tr>
    <th>ログインページ</th>
    <th>新規登録ページ</th>
  </tr>
   <tr>
    <td><img width="500" alt="image" src="https://github.com/user-attachments/assets/3a37cd62-166f-4bf4-849c-fec91d67bf4b" /></td>
    <td><img width="500" alt="image" src="https://github.com/user-attachments/assets/82e38012-f79e-4944-aa63-b4627d56e715" /></td>
  </tr>
  <tr>
    <th>Home</th>
    <th>記録の保存画面</th>
  </tr>
  <tr>
    <td><img width="500" alt="image" src="https://github.com/user-attachments/assets/92f17f8a-924f-4c1f-befa-c45ad915d18f" /></td>
    <td><img width="500" alt="image" src="https://github.com/user-attachments/assets/2d65ade8-4998-4631-91e0-4edd1ff78a55" /></td>
  </tr>
  <tr>
    <th>カレンダー表示画面</th>
    <th>記録編集モーダル</th>
  </tr>
  <tr>
    <td><img width="500" alt="image" src="https://github.com/user-attachments/assets/dcfe3254-8a5a-4fa4-ab8e-246d0bf6d858" /></td>
    <td><img width="500" alt="image" src="https://github.com/user-attachments/assets/371a14a5-75ac-465a-89a9-5e4fc04988d7" /></td>
  </tr>
  <tr>
    <th>グラフ表示画面</th>
    <th>プロフィール画面</th>
  </tr>
  <tr>
    <td><img width="500" height="700" alt="image" src="https://github.com/user-attachments/assets/af594c33-252e-45a1-b3f4-159bdd749e2b" /></td>
    <td><img width="500" height="700" alt="image" src="https://github.com/user-attachments/assets/63cceed7-6e28-4e16-996c-137a9e12a3af" /></td>
  </tr>
  <tr>
    <th>ランキング画面</th>
  </tr>
  <tr>
    <td><img width="500" height="700" alt="image" src="https://github.com/user-attachments/assets/f1165d89-e130-42bf-accf-1ee5ee084adf" /></td>
  </tr>
</table>

## 技術スタック

### 使用言語

・**[Dart](https://dart.dev/docs)**<br>
・**[Go](https://go.dev/doc/)**

### フロントエンド

・**[flutter](https://docs.flutter.dev/)**: ユーザーインターフェースの構築<br>
・**[Riverpod](https://riverpod.dev/ja/docs/introduction/getting_started)**: 状態管理と依存性注入<br>
・**[fl_chart](https://pub.dev/documentation/fl_chart/latest/)**: 折れ線／棒グラフの描画<br>
・**[table_calendar](https://pub.dev/documentation/table_calendar/latest/)**: 月間カレンダー表示とイベントドット<br>

### バックエンド

・**[echo](https://echo.labstack.com/docs)**: 高速な Web フレームワーク<br>
・**[Gorm](https://gorm.io/ja_JP/docs/index.html)**: データベースアクセスのための ORM<br>
・**[JWT](https://jwt.io/)**: JWT を使用した認証管理<br>

### データベース

・**[PostgresSQL](https://www.postgresql.org/docs/)**: リレーショナルデータベース管理システム

### インフラ

・**[Google Cloud Platform](https://cloud.google.com/)**: Google が提供するクラウドコンピューティングサービス<br>
・**[Neon](https://neon.com/)**: サーバーレスで自動スケーリング可能なクラウドデータベースサービス<br>
・**[Docker](https://docs.docker.com/)**: コンテナ化プラットフォームで環境構築を効率化<br>

### アーキテクチャ概要

バックエンドは GitHub Actions から Docker イメージを Artifact Registry に push し、
Cloud Run にデプロイしています。データベースは Neon を利用しています。

<img width="990" height="416" alt="image" src="https://github.com/user-attachments/assets/cc0da8f0-56d6-4432-8ab0-bd3a2b1decac" />

より詳細な構成は [`docs/infra/infra_design.md`](docs/infra/infra_design.md) を参照

---

## 本番環境

### **フロントエンド(デプロイ準備中)**

・ フロントはローカル環境からでも Cloud Run の本番 API を直接叩けます。

### **バックエンド(公開済)**

##### API ベース URL

- https://muscle-diary-backend-259090965328.asia-northeast1.run.app

---

## 開発環境のセットアップ手順

ローカル環境で開発サーバーを起動するための手順は以下の通りです。

1. リポジトリをクローン

```bash
git clone https://github.com/RintaroNasu/muscle_diary_app.git
```

### 前提環境

本プロジェクトでは以下の環境が必要です：

- Flutter 開発環境（FVM / Android Studio / Xcode）
  → 初めての方は https://www.sejuku.net/blog/123973 を推奨

- Go（バックエンド）
  → Go の公式インストール手順：https://go.dev/doc/install

### フロントエンド側セットアップ

2. フロントエンドディレクトリへ移動

```bash
cd frontend
```

3. 依存関係を取得

```bash
fvm install        # .fvmrc のバージョンを取得
fvm flutter pub get
```

4. サーバー立ち上げ

```bash
flutter devices #接続可能デバイスを確認
flutter run -d 〇〇
```

### バックエンド側セットアップ

5. バックエンドディレクトリへ移動

```bash
cd backend
```

6. ルートディレクトリに .env ファイルを作成し、以下の内容を追加

```bash
POSTGRES_USER="your_postgres_user"
POSTGRES_PASSWORD="your_postgres_password"
POSTGRES_DB="muscle_diary"
POSTGRES_HOST="localhost"
POSTGRES_PORT="5433"

JWT_SECRET="your_jwt_secret"
```

7. 依存関係の取得

```bash
go mod tidy
```

8. Docker コンテナを起動

```bash
docker compose up -d
```

9. バックエンドサーバーを起動

```bash
go run cmd/server/main.go
```

## フロントのローカル環境で本番 API を使用する方法

通常はローカル API が使われますが、以下のように --dart-define をつけて起動することで
本番 API に切り替えて起動できます。

```bash
flutter run \
  --dart-define=API_BASE_URL=https://muscle-diary-backend-259090965328.asia-northeast1.run.app
```
