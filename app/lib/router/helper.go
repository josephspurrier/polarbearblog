package router

// Delete is a shortcut for router.Handle("DELETE", path, handle)
func (m *Mux) Delete(path string, fn HandlerFunc) {
	m.router.Handle("DELETE", path, handler{
		HandlerFunc:     fn,
		CustomServeHTTP: m.customServeHTTP,
	})
}

// Get is a shortcut for router.Handle("GET", path, handle)
func (m *Mux) Get(path string, fn HandlerFunc) {
	m.router.Handle("GET", path, handler{
		HandlerFunc:     fn,
		CustomServeHTTP: m.customServeHTTP,
	})
}

// Head is a shortcut for router.Handle("HEAD", path, handle)
func (m *Mux) Head(path string, fn HandlerFunc) {
	m.router.Handle("HEAD", path, handler{
		HandlerFunc:     fn,
		CustomServeHTTP: m.customServeHTTP,
	})
}

// Options is a shortcut for router.Handle("OPTIONS", path, handle)
func (m *Mux) Options(path string, fn HandlerFunc) {
	m.router.Handle("OPTIONS", path, handler{
		HandlerFunc:     fn,
		CustomServeHTTP: m.customServeHTTP,
	})
}

// Patch is a shortcut for router.Handle("PATCH", path, handle)
func (m *Mux) Patch(path string, fn HandlerFunc) {
	m.router.Handle("PATCH", path, handler{
		HandlerFunc:     fn,
		CustomServeHTTP: m.customServeHTTP,
	})
}

// Post is a shortcut for router.Handle("POST", path, handle)
func (m *Mux) Post(path string, fn HandlerFunc) {
	m.router.Handle("POST", path, handler{
		HandlerFunc:     fn,
		CustomServeHTTP: m.customServeHTTP,
	})
}

// Put is a shortcut for router.Handle("PUT", path, handle)
func (m *Mux) Put(path string, fn HandlerFunc) {
	m.router.Handle("PUT", path, handler{
		HandlerFunc:     fn,
		CustomServeHTTP: m.customServeHTTP,
	})
}
