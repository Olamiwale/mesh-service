package main
import ( "encoding/json" "fmt" "log" "math/rand" "net/http" "os" "sync" "time" )

    type Payment struct {
      PaymentID     string    `json:"paymentId"`
      OrderID       string    `json:"orderId"`
      UserID        string    `json:"userId"`
      Amount        float64   `json:"amount"`
      Status        string    `json:"status"`
      ProcessedAt   time.Time `json:"processedAt"`
      TransactionID string    `json:"transactionId"`
    }

    type PaymentRequest struct {
      OrderID string  `json:"orderId"`
      UserID  string  `json:"userId"`
      Amount  float64 `json:"amount"`
    }

    var (
      payments = make(map[string]Payment)
      mu       sync.RWMutex
    )

    func init() {
      rand.Seed(time.Now().UnixNano())
    }

    func main() {
      http.HandleFunc("/payments", handlePayments)
      http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"status": "healthy - working fine"})
      })
      http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
      })

      port := os.Getenv("PORT")
      if port == "" { port = "8080" }
      log.Printf("Payment Service on port %s", port)
      http.ListenAndServe(":"+port, nil)
    }

    func handlePayments(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "application/json")
      if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", 405)
        return
      }

      var req PaymentRequest
      if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), 400)
        return
      }

      time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
      success := rand.Float32() > 0.1

      payment := Payment{
        PaymentID:     fmt.Sprintf("PAY-%d", time.Now().UnixNano()),
        OrderID:       req.OrderID,
        UserID:        req.UserID,
        Amount:        req.Amount,
        ProcessedAt:   time.Now(),
        TransactionID: fmt.Sprintf("TXN-%d", time.Now().UnixNano()),
      }

      if success {
        payment.Status = "completed"
      } else {
        payment.Status = "failed"
        w.WriteHeader(402)
      }

      mu.Lock()
      payments[payment.PaymentID] = payment
      mu.Unlock()

      json.NewEncoder(w).Encode(payment)
    }