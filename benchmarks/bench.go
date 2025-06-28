func BenchmarkOpaqueP256(b *testing.B) {
    client := opaque.NewClient()
    server := opaque.NewServer()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        initMsg, _ := client.SerializeRegistrationRequest()
        server.DeserializeRegistrationRequest(initMsg)
        // Complete protocol...
    }
}
