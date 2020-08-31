# weworkapi_cplusplus
official lib of wework api https://work.weixin.qq.com/api/doc

# 注意事项

* 1.回调sdk json版本

* 2.wxbizjsonmsgcrypt.go文件中声明并实现了WXBizJsonMsgCrypt类，提供用户接入企业微信的三个接口。sample.go文件提供了如何使用这三个接口的示例。

* 3.WXBizJsonMsgCrypt类封装了VerifyURL, DecryptMsg, EncryptMsg三个接口，分别用于开发者验证回调url，收到用户回复消息的解密以及开发者回复消息的加密过程。使用方法可以参考sample.go文件。

* 4.加解密协议请参考企业微信官方文档。