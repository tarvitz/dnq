package common

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/tarvitz/dnq/pkg/telegram"
)

func TestAuth_Method(t *testing.T) {
	token := "1337"
	methodName := "Test"
	expected := fmt.Sprintf("%s%s/%s", telegram.BotAPIURL, token, methodName)
	if result := (&Auth{Token: token}).Method(methodName); result != expected {
		t.Errorf("expected: `%v`, got: `%v`", expected, result)
	}
}

func TestAuth_GetClient(t *testing.T) {
	t.Run("blank", func(in *testing.T) {
		cmd := Auth{}
		client := cmd.GetClient()
		if !reflect.DeepEqual(client, cmd.client) {
			in.Errorf("clients are expected to be the same.")
		}
	})

	t.Run("not-blank", func(in *testing.T) {
		client := telegram.NewClient("")
		cmd := Auth{client: client}

		if result := cmd.GetClient(); !reflect.DeepEqual(result, cmd.client) {
			in.Errorf("clients are expected to be the same.")
		}
	})
}

func TestAuth_SetClient(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		cmd := Auth{}
		client := telegram.NewClient("")
		cmd.SetClient(client)

		if !reflect.DeepEqual(cmd.client, client) {
			t.Errorf("clients are expected to be the same.")
		}
	})
}
