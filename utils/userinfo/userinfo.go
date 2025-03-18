package userinfo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
)

// @Title        userinfo.go
// @Description
// @Create       XdpCs 2025-03-18 上午10:43
// @Update       XdpCs 2025-03-18 上午10:43
// 验证签名 -- 确保数据未修改
func VerifySignature(SessionKey, rawData, Signature string) bool {
	//使用会话密钥创建SHA256算法的HMAC实例
	h := hmac.New(sha256.New, []byte(SessionKey))
	//将原始数据写入HMAC实例进行计算
	h.Write([]byte(rawData))
	//生成十六进制格式的签名结果
	computedSign := hex.EncodeToString(h.Sum(nil))
	//比较签名结果与传入的签名是否相等
	return computedSign == Signature
}

// 解码base64
func base64Decode(s string) ([]byte, error) {
	// 使用标准库进行 Base64 解码
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		zlog.Errorf("解码失败: %v", err)
		return nil, err
	}
	return decoded, nil
}

func DecryptData(sessionKey string, encryptedData, iv []byte) ([]byte, error) {
	// Base64解码 sessionKey
	key, err := base64Decode(sessionKey)
	if err != nil {
		zlog.Errorf("解码 sessionKey 失败: %v", err)
		return nil, err
	}
	// Base64解码 iv
	iv, err = base64Decode(string(iv))
	if err != nil {
		zlog.Errorf("解码 iv 失败: %v", err)
		return nil, err
	}

	// 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		zlog.Errorf("创建 AES 密码块失败: %v", err)
		return nil, err
	}

	// 检查加密数据长度是否是 AES 块大小的倍数
	if len(encryptedData)%aes.BlockSize != 0 {
		zlog.Errorf("加密数据长度不符合 AES 块大小")
		return nil, err
	}

	// 创建 CBC 解密模式，用iv初始化
	// CBC模式需要iv来增加随机性，防止同样的明文生成同样的密文
	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encryptedData))
	// 使用CryptBlocks方法进行解密
	mode.CryptBlocks(decrypted, encryptedData)

	// 去除 PKCS#7 填充
	unpadding := int(decrypted[len(decrypted)-1])
	if unpadding < 1 || unpadding > aes.BlockSize {
		zlog.Errorf("填充值无效")
		return nil, err
	}
	return decrypted[:len(decrypted)-unpadding], nil
}
