from flask import Flask, request, jsonify
import requests
import uuid
import os
from datetime import datetime

app = Flask(__name__)

# External services (use Kubernetes service names or env vars)
PAYMENT_SERVICE = os.getenv('PAYMENT_SERVICE', 'http://payment-service:8080')
INVENTORY_SERVICE = os.getenv('INVENTORY_SERVICE', 'http://inventory-service:8080')
SHIPPING_SERVICE = os.getenv('SHIPPING_SERVICE', 'http://shipping-service:8080')
NOTIFICATION_SERVICE = os.getenv('NOTIFICATION_SERVICE', 'http://notification-service:8080')

orders = {}


@app.route('/orders', methods=['POST'])
def create_order():
    try:
        data = request.json
        user_id = data.get('userId')
        items = data.get('items', [])

        if not user_id or not items:
            return jsonify({'error': 'Missing userId or items'}), 400

        total = sum(item['price'] * item['quantity'] for item in items)
        order_id = str(uuid.uuid4())

        order = {
            'id': order_id,
            'userId': user_id,
            'items': items,
            'total': total,
            'status': 'pending',
            'createdAt': datetime.utcnow().isoformat()
        }

        # ---- Payment Call ----
        try:
            payment_response = requests.post(
                f"{PAYMENT_SERVICE}/api/payments",
                json={'orderId': order_id, 'amount': total, 'userId': user_id},
                timeout=5
            )
            if payment_response.status_code != 200:
                order['status'] = 'payment_failed'
                orders[order_id] = order
                return jsonify({'error': 'Payment failed', 'order': order}), 400

            order['paymentId'] = payment_response.json().get('paymentId')

        except Exception as e:
            order['status'] = 'payment_failed'
            orders[order_id] = order
            return jsonify({'error': f'Payment error: {str(e)}', 'order': order}), 500

        # Payment successful
        order['status'] = 'confirmed'
        orders[order_id] = order

        # ---- Notification call (non-blocking) ----
        try:
            requests.post(
                f"{NOTIFICATION_SERVICE}/api/notifications",
                json={'userId': user_id, 'type': 'order_confirmed', 'orderId': order_id},
                timeout=1
            )
        except:
            pass

        return jsonify(order), 201

    except Exception as e:
        return jsonify({'error': f'Order failed: {str(e)}'}), 500


@app.route('/orders/<order_id>', methods=['GET'])
def get_order(order_id):
    order = orders.get(order_id)
    if not order:
        return jsonify({'error': 'Order not found'}), 404
    return jsonify(order)


@app.route('/health', methods=['GET'])
def health():
    return jsonify({'status': 'healthy', 'service': 'order'})


@app.route('/ready', methods=['GET'])
def ready():
    return jsonify({'status': 'ready', 'service': 'order'})


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=int(os.getenv('PORT', 8080)))
