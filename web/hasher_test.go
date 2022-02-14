package web_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Suicide40/test1/web"
)

type dummyHasher struct {
	writeFn func(p []byte) (n int, err error)
	sumFn   func(b []byte) []byte
	resetFn func()
}

func (d dummyHasher) Write(p []byte) (n int, err error) { return d.writeFn(p) }
func (d dummyHasher) Sum(b []byte) []byte               { return d.sumFn(b) }
func (d dummyHasher) Reset()                            { d.resetFn() }
func (d dummyHasher) Size() int                         { panic("implement me") }
func (d dummyHasher) BlockSize() int                    { panic("implement me") }

type dummyClient struct {
	GetFn func(s string) (*http.Response, error)
}

func (d dummyClient) Get(s string) (*http.Response, error) { return d.GetFn(s) }

func TestHasher_Calc(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		resp := []byte("<html></html>")
		md5sum := []byte("1234567890123456")

		client := dummyClient{
			GetFn: func(s string) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader(resp)),
				}, nil
			},
		}
		hasher := dummyHasher{writeFn: func(p []byte) (n int, err error) {
			return len(resp), nil
		},
			resetFn: func() {},
			sumFn: func(b []byte) []byte {
				return md5sum
			},
		}

		obj := web.NewHasher(client, hasher)

		res, err := obj.Calc("www.test.com")
		if err != nil {
			t.Log(err.Error())
			t.FailNow()
		}

		if !bytes.Equal(res, md5sum) {
			t.Logf("md5 sums are different %s != %s", res, md5sum)
			t.FailNow()
		}
	})
}
