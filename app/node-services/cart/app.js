    const express = require('express');
    const axios = require('axios');
    const app = express();
    app.use(express.json());

    const PORT = process.env.PORT || 8080;
    const PRODUCT_SERVICE = 'http://product-service:8080';
    const INVENTORY_SERVICE = 'http://inventory-service:8080';

    const carts = new Map();

    app.post('/cart/:userId/items', async (req, res) => {
      try {
        const { userId } = req.params;
        const { productId, quantity } = req.body;

        const productResponse = await axios.get(`${PRODUCT_SERVICE}/api/products/${productId}`, 
        { timeout: 3000 });
        
        if (!carts.has(userId)) {
          carts.set(userId, { items: [], userId });
        }

        const cart = carts.get(userId);
        const existingItem = cart.items.find(item => item.productId === productId);

        if (existingItem) {
          existingItem.quantity += quantity;
        } else {
          cart.items.push({
            productId,
            name: productResponse.data.name,
            price: productResponse.data.price,
            quantity
          });
        }

        res.json(cart);
      } catch (error) {
        res.status(500).json({ error: 'Failed to add item', details: error.message });
      }
    });


    app.get('/cart/:userId', (req, res) => {
      const cart = carts.get(req.params.userId) || { items: [], userId: req.params.userId };

      const total = cart.items.reduce((sum, item) => sum + (item.price * item.quantity), 0);

      res.json({ ...cart, total });
    });

    app.delete('/cart/:userId/items/:productId', (req, res) => {
      const { userId, productId } = req.params;
       if (!carts.has(userId)) 
        return res.status(404).json({ error: 'Cart not found' });
      
      const cart = carts.get(userId);
      cart.items = cart.items.filter(item => item.productId !== productId);
      res.json(cart);
    });
    
    app.get('/health', (req, res) => res.json({ status: 'healthy', service: 'cart' }));
    app.get('/ready', (req, res) => res.json({ status: 'ready', service: 'cart' }));


    app.listen(PORT, () => console.log(`Cart Service on port ${PORT}`));