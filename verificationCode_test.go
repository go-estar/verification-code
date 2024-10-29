package verificationCode

import (
	"github.com/go-estar/redis"
	"testing"
)

func TestName(t *testing.T) {
	var TestCode = &VerificationCode{
		Redis: redis.New(&redis.Config{
			Addr:     "127.0.0.1:6379",
			Password: "GBkrIO9bkOcWrdsC",
		}),
		Name:         "test-code",
		ExpireTime:   600,
		IntervalTime: 60,
	}
	code, err := TestCode.CreateCode("a", 6)
	if err != nil {
		t.Fatal(err)
	}
	code, err = TestCode.CreateCode("a", 6)
	if err != nil {
		t.Fatal(err)
	}
	if err := TestCode.Verify("a", code); err != nil {
		t.Fatal(err)
	}
	t.Log("Succeed")
}
