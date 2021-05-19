# go-eth-tx-speedup

[![build](https://github.com/moremorefun/go-eth-tx-speedup/workflows/build/badge.svg)](https://github.com/moremorefun/go-eth-tx-speedup/actions?query=workflow%3Abuild)
[![GitHub release](https://img.shields.io/github/tag/moremorefun/go-eth-tx-speedup.svg?label=release)](https://github.com/moremorefun/go-eth-tx-speedup/releases)
[![GitHub release date](https://img.shields.io/github/release-date/moremorefun/go-eth-tx-speedup.svg)](https://github.com/moremorefun/go-eth-tx-speedup/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://github.com/moremorefun/go-eth-tx-speedup/blob/master/LICENSE)
[![blog](https://img.shields.io/badge/blog-@moremorefun-brightgreen.svg)](https://www.jidangeng.com)


## 目录

- [go-eth-tx-speedup](#go-eth-tx-speedup)
  - [目录](#目录)
  - [背景](#背景)
  - [使用说明](#使用说明)
  - [维护者](#维护者)
  - [使用许可](#使用许可)

## 背景

由于部分时候eth的gasPrice过低,上链太慢,于是写了一个加速现有tx的工具

替换tx的原理在于,如果eth的交易已经广播并且还没打包进block,
那么重新广播一个`nonce值相同`并且`gasPrice大于之前交易`的tx,就会用新的tx替换旧的tx

由于程序调用需要使用私钥,所以请确保你知道自己在做什么再使用此工具.


## 使用说明
下载release中自己平台的压缩包
解压进入目录

```
./goethspeedup
Usage of ./gomysql2struct:
  -gas int
    	gas price value in gwei (default 10)
  -h	help message
  -key string
    	eth address private key
  -limit uint
    	gas limit of tx
    	default is 0, when limit is 0, it will keep the old gas limit in the original tx
  -swap string
    	rpc uri of eth
  -txid string
    	txid to speed up
```
   
## 维护者

[@moremorefun](https://github.com/moremorefun)
[那些年我们De过的Bug](https://www.jidangeng.com)

## 使用许可

[MIT](LICENSE) © moremorefun
