package poker

//type ContextStub struct {
//}
//
//func (c ContextStub) Deadline() (time.Time, bool) {
//	return time.Time{}, false
//}
//
//func (c ContextStub) Done() <-chan struct{} {
//	return nil
//}
//
//func (c ContextStub) Err() error {
//	return nil
//}
//
//func (c ContextStub) Value(key interface{}) interface{} {
//	return nil
//}
//
//func TestServerServe(t *testing.T) {
//	t.Run("Serve", func(t *testing.T) {
//
//		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			w.WriteHeader(http.StatusOK)
//		}))
//
//		ctxStub := &ContextStub{}
//		//serve(ctxStub, server)
//		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
//
//	})
//}
