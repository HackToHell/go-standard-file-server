package main

import "testing"

func TestSign_in(t *testing.T) {
	SECRET_KEY_BASE = "qA6irmDikU6RkCM4V0cJiUJEROuCsqTa1esexI4aWedSv405v8lw4g1KB1nQVsSdCrcyRlKFdws4XPlsArWwv9y5Xr5Jtkb11w1NxKZabOUa7mxjeENuCs31Y1Ce49XH9kGMPe0ms7iV7e9F6WgnsPFGOlIA3CwfGyr12okas2EsDd71SbSnA0zJYjyxeCVCZJWISmLB"
	test := User{Uuid:"07a1e7e1-408e-4b8f-8960-2c26862a91c3", encrypted_password:"8e58a72e7c2cf845a510b1eab41d48a8908c533a696dd033d84ad62cc3314785"}
	if test.sign_in("test") == "" {
		t.Error("Failed")
	}

}
