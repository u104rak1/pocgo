# Infrastructure層について
この層ではDB、ORM及びDomain層に依存する事ができる。pocgoではDBにPostgreSQL、ORMにbunを採用している。

## ユニットテストについて
レコードが無い場合のエラーなどpackageに依存するエラーが外部に漏れないように調整する。