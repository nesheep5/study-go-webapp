package main

import "testing"

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
