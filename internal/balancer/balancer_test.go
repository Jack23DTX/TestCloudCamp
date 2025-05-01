package balancer

import (
	"net/url"
	"testing"
)

func TestServerPool_GetNextBackend(t *testing.T) {
	backends := []*Backend{
		{URL: &url.URL{Host: "server1"}, Alive: true},
		{URL: &url.URL{Host: "server2"}, Alive: true},
		{URL: &url.URL{Host: "server3"}, Alive: false},
	}
	pool := NewServerPool(backends)

	// Тест round-robin с "живыми" бекендами
	expectedOrder := []string{"server2", "server1", "server2"}
	for _, expected := range expectedOrder {
		backend := pool.GetNextBackend()
		if backend.URL.Host != expected {
			t.Errorf("Expected %s, got %s", expected, backend.URL.Host)
		}
	}

	// Тест если все бекенды "мертвы"
	backends[0].Alive = false
	backends[1].Alive = false
	if pool.GetNextBackend() != nil {
		t.Error("Expected nil when all servers are dead")
	}
}

// Проверка на несуществующий сервер
func TestHealthCheck(t *testing.T) {
	backend, _ := NewBackend("http://localhost:12345")
	pool := NewServerPool([]*Backend{backend})

	pool.checkBackend(backend)
	if backend.IsAlive() {
		t.Error("Dead server should be marked as not alive")
	}
}

func BenchmarkGetNextBackend(b *testing.B) {
	backends := make([]*Backend, 10)
	for i := 0; i < 10; i++ {
		backends[i] = &Backend{URL: &url.URL{Host: "server"}, Alive: true}
	}
	pool := NewServerPool(backends)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.GetNextBackend()
	}
}
