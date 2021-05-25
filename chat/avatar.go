package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrNoAvatarはAvatarインスタンスがアバターのURLを返すことができない
// 場合に発生するエラー。
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")

// Avatarはユーザのプロフィール画像を表す型です。
type Avatar interface {
	// GetAvatarURLは指定されたクライアントのアバターのURLを返す。
	// 問題が発生した場合にはエラーを返す。特に、URLを取得できなかった
	// 場合にはErrNoAvatarを返す。
	GetAvatarURL(c *client) (string, error) // 接頭辞の"Get"は必要なければつけない方が望ましい。

}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			println("===========")
			print(urlStr)
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(emailStr))
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}
	return "", ErrNoAvatarURL
}
