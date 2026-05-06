const express = require('express');
const app = express();
app.use(express.json());
    
const shipments = new Map();



app.post('/shipments', (req, res) => {
  const { orderId, userId, items } = req.body;
  const trackingNumber = `TRACK-${Date.now()}`;
      
  const shipment = { trackingNumber, orderId, userId, items, status: 'processing',
        estimatedDelivery: new Date(Date.now() + 5 * 24 * 60 * 60 * 1000)
     };
      
    shipments.set(trackingNumber, shipment);
    res.json(shipment);
    });
    
app.get('/shipments/:trackingNumber', (req, res) => {
    const shipment = shipments.get(req.params.trackingNumber);
      if (!shipment) return res.status(404).json({ error: 'Not found' });
      res.json(shipment);
    });
    
app.get('/health', (req, res) => { res.json({ status: 'healthy working fine', service: 'shipping' });  });
    
app.listen(8080, () => {
  console.log('Shipping service working')
})