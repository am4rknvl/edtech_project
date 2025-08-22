package main

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/99designs/gqlgen/graphql/handler"
    "github.com/am4rknvl/edtech_project/graph"
    "github.com/am4rknvl/edtech_project/graph/generated"
)

func newTestServer() *httptest.Server {
    srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
    mux := http.NewServeMux()
    mux.Handle("/query", srv)
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })
    return httptest.NewServer(mux)
}

func TestHealthEndpoint(t *testing.T) {
    ts := newTestServer()
    defer ts.Close()

    res, err := http.Get(ts.URL + "/health")
    if err != nil {
        t.Fatalf("health GET error: %v", err)
    }
    defer res.Body.Close()
    if res.StatusCode != http.StatusOK {
        t.Fatalf("expected 200 got %d", res.StatusCode)
    }
    b, _ := io.ReadAll(res.Body)
    if string(b) != "ok" {
        t.Fatalf("unexpected body: %s", string(b))
    }
}

func TestGraphQLSubjects(t *testing.T) {
    ts := newTestServer()
    defer ts.Close()

    payload := `{"query":"query { subjects { id name } }"}`
    res, err := http.Post(ts.URL+"/query", "application/json", bytes.NewBufferString(payload))
    if err != nil {
        t.Fatalf("graphql POST error: %v", err)
    }
    defer res.Body.Close()
    if res.StatusCode != http.StatusOK {
        t.Fatalf("expected 200 got %d", res.StatusCode)
    }
    var got map[string]any
    if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
        t.Fatalf("decode response: %v", err)
    }
    data, ok := got["data"].(map[string]any)
    if !ok {
        t.Fatalf("no data in response")
    }
    subjects, ok := data["subjects"].([]any)
    if !ok {
        t.Fatalf("subjects missing or wrong type")
    }
    if len(subjects) == 0 {
        t.Fatalf("expected at least one subject")
    }
    first, ok := subjects[0].(map[string]any)
    if !ok {
        t.Fatalf("unexpected subject shape")
    }
    if name, _ := first["name"].(string); name != "Math" {
        t.Fatalf("expected seeded subject name Math, got %v", first["name"])
    }
}
