package models

import (
	"testing"
	"time"
)

func TestSetScope(t *testing.T) {
	token := NewToken(1, time.Duration(24))

	token.SetScope("scope")
	if token.Scopes != "scope" {
		t.Fatalf("Scope ожидался: %s, получен: %s", "scope", token.Scopes)
	}

	token.SetScope("two_scope")
	if token.Scopes != "scope;two_scope" {
		t.Fatalf("Scope ожидался: %s, получен: %s", "scope;two_scope", token.Scopes)
	}

	token.SetScope("three_scope")
	if token.Scopes != "scope;two_scope;three_scope" {
		t.Fatalf("Scope ожидался: %s, получен: %s", "scope;two_scope;three_scope", token.Scopes)
	}
}

func TestSetScopes(t *testing.T) {
	token := NewToken(1, time.Duration(24))

	token.SetScope("scope")
	token.SetScope("two_scope")
	if token.Scopes != "scope;two_scope" {
		t.Fatalf("Scope ожидался: %s, получен: %s", "scope;two_scope", token.Scopes)
	}

	token.SetScopes([]string{"three_scope", "four_scope"})
	if token.Scopes != "scope;two_scope;three_scope;four_scope" {
		t.Fatalf("Scope ожидался: %s, получен: %s", "scope;two_scope;three_scope;four_scope", token.Scopes)
	}

	token = NewToken(1, time.Duration(24))

	token.SetScopes([]string{"three_scope", "four_scope"})
	if token.Scopes != "three_scope;four_scope" {
		t.Fatalf("Scope ожидался: %s, получен: %s", "three_scope;four_scope", token.Scopes)
	}

}

func TestGetScope(t *testing.T) {
	token := NewToken(1, time.Duration(24))

	token.SetScope("scope")
	if token.GetScope("scope") != true {
		t.Fatalf("scope не найден")
	}
	if token.GetScope("error_scope") != false {
		t.Fatalf("Должна возвращать false, т.к. error_scope нет")
	}

	token.SetScope("two_scope")
	if token.Scopes != "scope;two_scope" {
		t.Fatalf("Ожидаемый scope: %s, получен: %s", "scope;two_scope", token.Scopes)
	}
	if token.GetScope("scope") != true {
		t.Fatalf("scope не найден")
	}
	if token.GetScope("two_scope") != true {
		t.Fatalf("two_scope не найден")
	}
	if token.GetScope("error_scope") != false {
		t.Fatalf("Должна возвращать false, т.к. error_scope нет")
	}

}

func TestValid(t *testing.T) {
	token := NewToken(1, time.Duration(24))

	token.Expires = time.Now().Add(time.Second * 24).Unix()
	if token.Valid() == false {
		t.Fatalf("Токен валиден, но Valid возвращает - false")
	}

	token.Expires = time.Now().Add(time.Second * -100).Unix()
	if token.Valid() == true {
		t.Fatalf("Токен просрочен, но Valid возвращает - true")
	}

}
