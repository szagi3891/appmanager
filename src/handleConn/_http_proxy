diff --git a/src/backend/manager.go b/src/backend/manager.go
index 44f2a9f..d91ce6d 100644
@@ -14,8 +14,12 @@ import (
     utils "../utils"
     "io/ioutil"
     "fmt"
-    "../handleConn"
+    //"../handleConn"
     configModule "../config"
+    "net/http"
+    "net/http/httputil"
+    "net/url"
+    "time"
 )
 
 
@@ -53,21 +57,57 @@ func Init(config *configModule.File, logrotor *logrotorModule.Manager, logs *log
     
     manager.backend = backend
     
+    addr := "127.0.0.1:" + strconv.FormatInt(int64(config.GetPortMain()), 10)
     
     
-    addr := "127.0.0.1:" + strconv.FormatInt(int64(config.GetPortMain()), 10)
+    backendAddr := backend.GetAddr()
     
-    errStart := handleConn.Start(addr, logs, func() (string, func(), func()) {
+    fmt.Println("backendAddr", backendAddr)
+    
+    backendAddrUrl, errParse := url.Parse(backendAddr)
+    
+    if errParse != nil {
+        return nil, errorStack.From(errParse)
+    }
+    
+    fmt.Println("backendAddrUrl", backendAddrUrl)
+    
+    //wystartuj serwer httpRevProxy
+    //wystartuj serwer http z przełączeniem na ten serwer RevProxy
+    
+    go func() {
         
-        backend := manager.GetActiveBackend()
+        httpProxy := httputil.NewSingleHostReverseProxy(backendAddrUrl)
+        
+        httpFront := &http.Server{
+            Addr:           addr,
+            Handler:        httpProxy,
+            ReadTimeout:    10 * time.Second,
+            WriteTimeout:   10 * time.Second,
+            MaxHeaderBytes: 1 << 20,
+        }
         
+        errStart := httpFront.ListenAndServe()
+        
+        panic(errStart)
+    }()
+    
+    //ready := make(chan bool)
+    
+    /*
+    
+
+    errStart := handleConn.Start(addr, logs, func() (string, func(), func()) {
+
+        backend := manager.GetActiveBackend()
+
         return backend.GetAddr(), backend.Inc, backend.Sub
     })
-    
+        
     if errStart != nil {
         return nil, errStart
     }
-    
+    */
     
     
     return &manager, nil
