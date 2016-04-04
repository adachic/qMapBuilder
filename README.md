# qMapBuilder

qMapBuilderは、マップ自動生成ツールです。

オープンシンプレックス法で材質・勾配を計算し、自然なマップを無限に生成します。

adachicは、従来、qMapEditorで全てのマップを１パネルずつクラフトしていましたが、このツールを用いることにより、自動で大量のマップを生成し、その中から好きなものを選択できます。

## 関連ツール

- 出力されたメタデータは、qEnemyGeneratorに入力することによってさらに、敵情報を自動生成し、ステージ敵情報もセットで生成することができます。
- 出力されたマップ形状の一部を変更したい場合は、メタデータをqMapEditorで編集することで、お好みのマップへ変更することができます。


## 出力形式

- png
- pngに対応するjson（qMapEditor、qEnemyGenerator互換）


## 以下のマップの種類に対応

- 平原／洞窟／火山／毒の沼地／雪原／城内／遺跡

## 使い方

- go get  github.com/adachic/qMapBuilder
- go build
- ./qMapBuilder
- outputディレクトリ以下にファイルが出力されます


## 出力されるマップの例

![](15ee7878-9752-4bf4-9ef2-2e675e95164e.png)

![](2996d28d-b094-4e84-bfce-ffcbec6bda14.png)

![](c5f24c75-1fa0-4d4f-b7cf-3590a0337ea4.png)

