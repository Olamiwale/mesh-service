package main

import ( 
  "encoding/json" 
  "log" 
  "net/http" 
  "os" 
  "sync" 
  "time" 
)

type Product struct {
      ID          string  `json:"id"`
      Name        string  `json:"name"`
      Description string  `json:"description"`
      Price       float64 `json:"price"`
      Stock       int     `json:"stock"`
      Version     string  `json:"version"`
}

    var (
      products = make(map[string]Product)
      mu       sync.RWMutex
      version  = os.Getenv("VERSION")
    )

    func init() {
      if version == "" {
        version = "v1"
      }
      products["1"] = Product{ID: "1", Name: "Laptop",   Description: "High-performance laptop", Price: 999.99, Stock: 50,  Version: version}
      products["2"] = Product{ID: "2", Name: "Mouse",    Description: "Wireless mouse",          Price: 29.99,  Stock: 200, Version: version}
      products["3"] = Product{ID: "3", Name: "Keyboard", Description: "Mechanical keyboard",     Price: 79.99,  Stock: 150, Version: version}
    }

    func main() {
      http.HandleFunc("/products", handleProducts)
      http.HandleFunc("/products/", handleProductByID)
      http.HandleFunc("/health", handleHealth)
      http.HandleFunc("/ready", handleReady)

      port := os.Getenv("PORT")
      if port == "" {
        port = "8080"
      }

      log.Printf("Product Service %s starting on port %s", version, port)
      if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
      }
    }

    func handleProducts(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("X-Service-Version", version)

      if version == "v2" {
        time.Sleep(50 * time.Millisecond)
      }

      mu.RLock()
      defer mu.RUnlock()

      switch r.Method {
      case http.MethodGet:
        productList := make([]Product, 0, len(products))
        for _, p := range products {
          productList = append(productList, p)
        }
        json.NewEncoder(w).Encode(map[string]interface{}{
          "products": productList,
          "version":  version,
        })
      default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
      }
    }

    func handleProductByID(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "application/json")
      id := r.URL.Path[len("/products/"):]
      mu.RLock()
      product, exists := products[id]
      mu.RUnlock()
      if !exists {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
      }
      json.NewEncoder(w).Encode(product)
    }

    func handleHealth(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "version": version})
    }

    func handleReady(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(map[string]string{"status": "ready", "version": version})
    }