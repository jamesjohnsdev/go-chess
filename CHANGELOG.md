# Changelog

## [0.1.1](https://github.com/jamesjohnsdev/go-chess/compare/v0.1.0...v0.1.1) (2026-08-20)


### Features

* add resignation and draw offers ([ed14268](https://github.com/jamesjohnsdev/go-chess/commit/ed142686e835c75b93cf703d7ce4441754d8d2bc))
* add web server and static client placeholder ([178c5c6](https://github.com/jamesjohnsdev/go-chess/commit/178c5c63f973b94a625756cf7745c8b6e8e14bf9))
* **cmd/chess:** add -fen flag and fen command ([aaa427a](https://github.com/jamesjohnsdev/go-chess/commit/aaa427a9e2e89d15f7d7d83f2304ce9c2e4e5328))
* **cmd/chess:** add terminal chess CLI ([d46c615](https://github.com/jamesjohnsdev/go-chess/commit/d46c615cbc5e7ebd956d477ce81f56ca171ed887))
* **cmd/chess:** play against the computer ([3c98ae6](https://github.com/jamesjohnsdev/go-chess/commit/3c98ae656c7b1bc9d8f4c8cd665bf5f73bc58bd5))
* **cmd/chess:** report all draw types ([9cfe34a](https://github.com/jamesjohnsdev/go-chess/commit/9cfe34ad4f3af3f74256bdb4cca300a4a5685714))
* **cmd/chess:** validate moves and report game state ([2aca959](https://github.com/jamesjohnsdev/go-chess/commit/2aca959a5bc146b014b0998f54906e69102e7621))
* **engine/ai:** add computer opponent ([2594795](https://github.com/jamesjohnsdev/go-chess/commit/25947958bca736e8acba9da84d669dde9ee1afc8))
* **engine/ai:** add piece-square tables to evaluation ([49e55ea](https://github.com/jamesjohnsdev/go-chess/commit/49e55ea89676bf5a536b8e91b3539076e43bc728))
* **engine:** add board, piece, square and move types ([0d79efa](https://github.com/jamesjohnsdev/go-chess/commit/0d79efa0a1a6f4fcf8f1bdd0fbb6f62151b8ef51))
* **engine:** add draw detection ([84ec0d8](https://github.com/jamesjohnsdev/go-chess/commit/84ec0d891581103e4441f10b1a92f9ea5959ea24))
* **engine:** add FEN import/export ([4f097e3](https://github.com/jamesjohnsdev/go-chess/commit/4f097e386aa9c44a9dafc32fb8e7a900f8e3a4d1))
* **engine:** add legal move generation ([b0d37f1](https://github.com/jamesjohnsdev/go-chess/commit/b0d37f1f6ca2782c5a57bf44b287271cacd99247))
* **internal/game:** add computer opponent to live sessions ([d3e6fc7](https://github.com/jamesjohnsdev/go-chess/commit/d3e6fc7f7cf9f580235c59771fb37cc845a30b4b))
* **internal/game:** add live game session model ([1bdd36a](https://github.com/jamesjohnsdev/go-chess/commit/1bdd36ace70fc441595232ef9b0f2e686a19f3be))
* **internal/game:** add live game session store ([9f11752](https://github.com/jamesjohnsdev/go-chess/commit/9f1175253e7e997062d956cee2f5707ea784a7a3))
* **internal/game:** reject moves once a game is over ([8652467](https://github.com/jamesjohnsdev/go-chess/commit/8652467d01d705e1c72795abac75a9a0fbe09c4f))
* **internal/game:** surface FEN in game state ([8487d87](https://github.com/jamesjohnsdev/go-chess/commit/8487d87bc6a1932ed4d45d3126790f7834060223))
* **internal/httpserver:** add live game and WebSocket endpoints ([e83e3ee](https://github.com/jamesjohnsdev/go-chess/commit/e83e3eee7b6822317421b1dfb6b5a2d661954870))
* **internal/httpserver:** create live games vs the computer ([ad9d690](https://github.com/jamesjohnsdev/go-chess/commit/ad9d690aa078033328bf1d6fb1482de022c81366))
* **internal/httpserver:** switch to chi router and huma API ([02072ce](https://github.com/jamesjohnsdev/go-chess/commit/02072ce29c7aaee5a4ef350418ef072e4e0bf62a))


### Bug Fixes

* **internal/httpserver:** make create-game request body optional ([0cec7e2](https://github.com/jamesjohnsdev/go-chess/commit/0cec7e26aa13f9ab4876c8d35abec6622a1a7730))
