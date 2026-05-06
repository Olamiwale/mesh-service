
from flask import Flask, request, jsonify
import os
    
app = Flask(__name__)

inventory = {
      '1': {'productId': '1', 'stock': 50, 'reserved': 0},
      '2': {'productId': '2', 'stock': 200, 'reserved': 0},
      '3': {'productId': '3', 'stock': 150, 'reserved': 0}
    }
    
@app.route('/inventory/<product_id>')
def get_inventory(product_id):
 item = inventory.get(product_id)
 if not item:
    return jsonify({'error': 'Not found'}), 404
    return jsonify(item)
    
@app.route('/inventory/reserve', methods=['POST'])
def reserve_inventory():
      data = request.json
      product_id = data.get('productId')
      quantity = data.get('quantity', 0)
      
      if product_id not in inventory:
        return jsonify({'error': 'Product not found'}), 404
      
      item = inventory[product_id]
      available = item['stock'] - item['reserved']
      
      if available < quantity:
        return jsonify({'error': 'Insufficient stock'}), 400
      
      item['reserved'] += quantity
      return jsonify({'success': True, 'reserved': quantity})
    
@app.route('/health')
def health():
      return jsonify({'status': 'healthy', 'service': 'inventory'})
    
if __name__ == '__main__':
     app.run(host='0.0.0.0', port=int(os.getenv('PORT', 8080)))
