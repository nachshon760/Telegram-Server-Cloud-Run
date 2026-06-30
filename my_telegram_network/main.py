import os
import socket
import sys
from core.database import get_user_stars
from core.protocol import generate_server_keys, process_telegram_packet

def start_custom_telegram_server():
    # גוגל קלאוד מזרים את הפורט למשתנה הסביבה PORT באופן אוטומטי
    port = int(os.environ.get("PORT", 8080))
    host = "0.0.0.0"
    
    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    
    try:
        server.bind((host, port))
        server.listen(500)
        print(f"🚀 YOUR PRIVATE TELEGRAM SERVER is running on cloud port: {port}")
    except Exception as e:
        print(f"❌ Core Error: Could not bind network port: {e}")
        sys.exit(1)
        
    # הפקת המפתחות הפרטיים של הרשת שלך
    private_key, public_key = generate_server_keys()
    print(f"[Security] Server Public Key generated (Hex): {public_key.hex()[:32]}...")
    
    try:
        while True:
            client_socket, addr = server.accept()
            data = client_socket.recv(8192)
            if data:
                cmd_id = process_telegram_packet(data)
                print(f"[Traffic] Received official packet from App. Command ID: {cmd_id}")
                # כאן השרת מחזיר את הנתונים המשונים (כמו כמות הכוכבים מבסיס הנתונים)
            client_socket.close()
    except KeyboardInterrupt:
        print("\n🛑 Server shutting down.")
    finally:
        server.close()

if __name__ == '__main__':
    start_custom_telegram_server()
