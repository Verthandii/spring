package ies

import (
	"context"
	"testing"
)

func TestNewES(t *testing.T) {
	client, err := Init(&Config{
		Address:  "https://es-cn-j4g3vsd13000fcmgp.public.elasticsearch.aliyuncs.com:9200",
		Username: "elastic",
		Password: "xxxx",
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Ping().Do(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp)
}
