# コミットガイドライン
このリポジトリに貢献していただきありがとうございます！ コミットメッセージの一貫性を保つために、以下のガイドラインに従ってください。
## コミットメッセージのフォーマット
コミットメッセージは以下の形式に従ってください：
```
<タイプ>(<スコープ>): <概要>
<空行>
<詳細な説明>
<空行>
<フッター>
```
### タイプ
- `feat`: 新機能の追加
- `fix`: バグ修正
- `docs`: ドキュメントの変更
- `style`: フォーマットの変更（コードの意味に影響しない
- `refactor`: リファクタリング（バグ修正や機能追加を伴わないコード変更）
- `test`: テストの追加や修正
- `chore`: その他の変更（ビルドプロセスや補助ツールの変更など）
### スコープ
スコープは変更が影響を与える部分を示します。例えば、`UI`、`API`、`データベース`などです。スコープが特にない場合は省略可能です。
### 概要
概要は50文字以内で、変更内容を簡潔に説明します。動詞の原形で始め、ピリオドで終わらないでください。
### 詳細な説明
詳細な説明は72文字以内の行で、変更の理由や背景を説明します。必要に応じて複数行に分けることができます。
### フッター
フッターには、関連するIssue番号やBREAKING CHANGEの情報を含めることができます。例えば：
```
Closes #123
BREAKING CHANGE: APIのエンドポイントが変更されました。
```
## 例
```
feat(UI): ユーザープロフィールページを追加
ユーザープロフィールページを新たに作成し、ユーザーが自分の情報を編集できるようにしました。
Closes #45
```
```fix(API): 認証エンドポイントのバグ修正
認証エンドポイントで発生していたトークン生成のバグを修正しました。
```
## 注意事項
- コミットメッセージは日本語で書くことを推奨します。
- 一つのコミットには一つの目的を持たせるようにしてください。
- 大きな変更は複数の小さなコミットに分けることを検討してください。
これらのガイドラインに従うことで、プロジェクトの履歴が整理され、他の貢献者が変更内容を理解しやすくなります。ご協力ありがとうございます！
## References
- [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)
- [Git Commit Message Conventions](https://chris.beams.io/posts/git-commit/)
- [How to Write a Git Commit Message](https://www.freecodecamp.org/news/how-to-write-a-git-commit-message/)
- [コミットメッセージの書き方](https://qiita.com/jyoshiki/items/4f0b2d1f8c6e5c3f4c1d)
- [良いコミットメッセージの書き方](https://zenn.dev/okazuki/articles/commit-message)
- [Gitのコミットメッセージのベストプラクティス](https://www.atlassian.com/git/tutorials/communicating-with-git/commit-messages)
