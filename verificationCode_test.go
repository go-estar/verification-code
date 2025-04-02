package verificationCode

import (
	"github.com/go-estar/redis"
	"testing"
	"time"
)

var TestCode = New(&Config{
	Redis: redis.New(&redis.Config{
		Addr:     "127.0.0.1:6379",
		Password: "GBkrIO9bkOcWrdsC",
	}),
	Name:         "test-code",
	ExpireTime:   600,
	IntervalTime: 60,
})

func TestCreate(t *testing.T) {
	code1, err := TestCode.Create("a")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("code1", code1)
	time.Sleep(time.Minute)
	code2, err := TestCode.Create("a")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("code2", code2)
}

func TestCreateWithMustRegen(t *testing.T) {
	code1, err := TestCode.Create("a")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("code1", code1)
	time.Sleep(time.Minute)
	code2, err := TestCode.Create("a", CreateWithMustRegen())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("code2", code2)
}

func TestCreateWithCode(t *testing.T) {
	code1, err := TestCode.Create("a", CreateWithCode("1234"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("code1", code1)
	time.Sleep(time.Minute)
	code2, err := TestCode.Create("a")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("code2", code2)
}

func TestVerify(t *testing.T) {
	if err := TestCode.Verify("a", "1234"); err != nil {
		t.Fatal(err)
	}
	t.Log("Succeed")
}

func TestCreateVerify(t *testing.T) {
	code, err := TestCode.Create("a")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("code", code)
	if err := TestCode.Verify("a", code); err != nil {
		t.Fatal(err)
	}
	t.Log("Succeed")
}
