// Copyright 2022 ByteDance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package utils

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/byted-apaas/server-common-go/constants"
)

func AesDecryptText(fieldID int64, realEncryptKey []byte, encryptedText string) (originText string, err error) {
	if len(realEncryptKey) != 32 {
		return "", fmt.Errorf("ilegal length")
	}
	iv, err := getInitialVector(fmt.Sprintf("%d", fieldID))
	if err != nil {
		return "", err
	}
	bs, err := hex2bin(encryptedText)

	if err != nil {

		return "", err
	}
	originBytes, err := aesCbsDecrypt(bs, iv, realEncryptKey)
	if err != nil {
		return "", err
	}

	return string(originBytes), nil
}

func paddingN(text []byte, size int) []byte {
	if len(text) > size {
		return text[:size]
	}

	return append(text, bytes.Repeat([]byte("0"), size-len(text))...)
}

func getInitialVector(str string) ([]byte, error) {
	h := md5.New()
	if _, err := io.WriteString(h, str); err != nil {
		return nil, err
	}
	return h.Sum(nil)[:], nil
}

func hex2bin(bs string) ([]byte, error) {
	if len(bs)%2 != 0 {
		return nil, hex.ErrLength
	}
	src := []byte(bs)
	dst := make([]byte, hex.DecodedLen(len(bs)))

	_, err := hex.Decode(dst, src)
	if err != nil {
		return nil, err

	}
	return dst, nil
}

func aesCbsDecrypt(cipherText, iv []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	plainText = unPaddingN(plainText)
	return plainText, nil
}

func unPaddingN(cipherText []byte) []byte {
	end := cipherText[len(cipherText)-1]
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
}

func TimeMils(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func NowMils() int64 {
	return TimeMils(time.Now())
}

func InvokeFuncWithRetry(retryCount int, retryInterval time.Duration, f func() error) error {
	var (
		err   error
		count = 0
	)
	for {
		if err = f(); err == nil {
			break
		}

		if count >= retryCount {
			break
		}

		time.Sleep(retryInterval)
		count++
	}
	return err
}

// PrintLog unittest to use
func PrintLog(contents ...interface{}) {
	isPrint := false

	for _, content := range contents {
		if content == nil {
			fmt.Println(content)
			isPrint = true
			continue
		}

		typ := reflect.TypeOf(content)
		val := reflect.ValueOf(content)
		if typ.Kind() == reflect.Ptr {
			val = val.Elem()
			typ = typ.Elem()
		}

		switch typ.Kind() {
		case reflect.String:
			fmt.Println(content)
			isPrint = true
		default:
			content, err := json.Marshal(content)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(content))
			isPrint = true
		}
	}

	if isPrint {
		fmt.Println()
	}
}

func GetEventID(ctx context.Context) string {
	persistHeaders, ok := ctx.Value(constants.PersistAPaaSKeySummarized).(map[string]string)
	if !ok || persistHeaders == nil {
		return ""
	}
	return persistHeaders[constants.HttpHeaderKeyEventID]
}

func GetLegacyLoggerDisabledFromCtx(ctx context.Context) bool {
	optionHeaders, ok := ctx.Value(constants.HTTPInvokeOptionsHeader).(map[constants.OptionKey]string)
	if !ok || optionHeaders == nil {
		return false
	}
	disabled, ok := optionHeaders[constants.gDisableLegacyLogger]
	if !ok {
		return false
	}
	return disabled == "true"
}
