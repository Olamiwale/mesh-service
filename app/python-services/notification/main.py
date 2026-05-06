from flask import Flask, request, jsonify
import os
from datetime import datetime

app = Flask(__name__)
notifications = []


@app.route('/api/notifications', methods=['POST'])
def send_notification():
    data = request.json
    notification = {
        'id': len(notifications) + 1,
        'userId': data.get('userId'),
        'type': data.get('type'),
        'orderId': data.get('orderId'),
        'message': f"Order {data.get('orderId')} has been {data.get('type')}",
        'timestamp': datetime.utcnow().isoformat(),
        'sent': True
    }
    notifications.append(notification)
    print(f"Notification sent: {notification['message']}")
    return jsonify(notification)


@app.route('/api/notifications/user/<user_id>')
def get_notifications(user_id):
    user_notifications = [n for n in notifications if n['userId'] == user_id]
    return jsonify({'notifications': user_notifications})


@app.route('/health')
def health():
    return jsonify({'status': 'healthy', 'service': 'notification'})


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=int(os.getenv('PORT', 8080)))