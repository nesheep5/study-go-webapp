package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.AvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、ErrNoAvatarURLを返すべきです")
	}
	// 値をセット
	testURL := "http://url-to-avatar"
	client.userData = map[string]interface{}{"avatar_url": testURL}
	url, err = authAvatar.AvatarURL(client)
	if err != nil {
		t.Error("値が存在する場合、エラーを返すべきではありません")
	} else {
		if url != testURL {
			t.Error("URLが一致しません")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData =
		map[string]interface{}{
			"userid": "0bc83cb571cd1c50ba6f3e8a78ef1346",
		}
	url, err := gravatarAvatar.AvatarURL(client)
	if err != nil {
		t.Error("AvatarURLはエラーを返すべきではありません")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("誤ったURLです URL: %s", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	// テスト用のアバターを作成
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	client := new(client)
	client.userData = map[string]interface{}{"userid": "abc"}
	url, err := fileSystemAvatar.AvatarURL(client)
	if err != nil {
		t.Error("FileSystemAvatar.AvatarURLはエラーを返すべきではありません")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("%s という誤った値を返しました", url)
	}
}
