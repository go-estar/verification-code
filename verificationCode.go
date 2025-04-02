package verificationCode

import (
	"context"
	"errors"
	baseError "github.com/go-estar/base-error"
	"github.com/go-estar/redis"
	"github.com/thoas/go-funk"
	"time"
)

var (
	ErrorVerificationCodeIdEmpty    = baseError.NewSystemCode("2701", "id can't be empty")
	ErrorVerificationCodeInterval   = baseError.NewCode("2702", "请在[%s]秒后再试")
	ErrorVerificationCodeVerifyId   = baseError.NewSystemCode("2703", "id can't be empty")
	ErrorVerificationCodeVerifyCode = baseError.NewSystemCode("2704", "code can't be empty")
	ErrorVerificationCodeExpired    = baseError.NewCode("2705", "验证码不正确")
	ErrorVerificationCodeVerify     = baseError.NewCode("2706", "验证码错误")
)

type Config struct {
	Redis        *redis.Redis
	Name         string
	Chars        string
	Length       int
	ExpireTime   int
	IntervalTime int
}

func New(c *Config) *VerificationCode {
	if c == nil {
		panic("config必须设置")
	}
	if c.Name == "" {
		panic("Name必须设置")
	}
	if c.Redis == nil {
		panic("Redis必须设置")
	}
	if c.Chars == "" {
		c.Chars = "012346789"
	}
	if c.Length == 0 {
		c.Length = 6
	}
	if c.ExpireTime == 0 {
		c.ExpireTime = 300
	}
	if c.IntervalTime == 0 {
		c.IntervalTime = 60
	}
	return &VerificationCode{c}
}

type VerificationCode struct {
	*Config
}
type CreateConfig struct {
	Code      string
	MustRegen bool
}
type CreateOption func(*CreateConfig)

func CreateWithCode(val string) CreateOption {
	return func(c *CreateConfig) {
		c.Code = val
	}
}
func CreateWithMustRegen() CreateOption {
	return func(c *CreateConfig) {
		c.MustRegen = true
	}
}

func (vc *VerificationCode) Create(id string, opts ...CreateOption) (string, error) {
	if id == "" {
		return "", ErrorVerificationCodeIdEmpty
	}

	conf := &CreateConfig{
		Code:      "",
		MustRegen: false,
	}
	for _, apply := range opts {
		apply(conf)
	}

	remainTime, err := vc.Redis.TTL(context.Background(), vc.Name+":i:"+id).Result()
	if err != nil {
		return "", err
	}
	if remainTime > 0 {
		return "", ErrorVerificationCodeInterval.Clone().SetMsgArgs(remainTime)
	}

	code := ""
	if conf.Code != "" {
		code = conf.Code
	} else {
		if !conf.MustRegen {
			val, err := vc.Redis.Get(context.Background(), vc.Name+":"+id).Result()
			if err != nil && !errors.Is(err, redis.Nil) {
				return "", err
			}
			if val != "" {
				code = val
			}
		}
		if code == "" {
			code = funk.RandomString(vc.Length, []rune(vc.Chars))
		}
	}

	ctx := context.Background()
	multi := vc.Redis.TxPipeline()
	multi.Set(ctx, vc.Name+":"+id, code, time.Duration(vc.ExpireTime)*time.Second)
	multi.Set(ctx, vc.Name+":i:"+id, time.Now(), time.Duration(vc.IntervalTime)*time.Second)
	_, err = multi.Exec(ctx)

	return code, nil
}

func (vc *VerificationCode) Verify(id string, code string) error {
	if id == "" {
		return ErrorVerificationCodeVerifyId
	}
	if code == "" {
		return ErrorVerificationCodeVerifyCode
	}

	val, err := vc.Redis.Get(context.Background(), vc.Name+":"+id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrorVerificationCodeExpired
		}
		return err
	}

	if val != code {
		return ErrorVerificationCodeVerify
	}

	vc.Remove(id)
	return nil
}

func (vc *VerificationCode) Remove(id string) error {
	ctx := context.Background()
	multi := vc.Redis.TxPipeline()
	multi.Del(ctx, vc.Name+":"+id)
	multi.Del(ctx, vc.Name+":i:"+id)
	_, err := multi.Exec(ctx)
	return err
}
