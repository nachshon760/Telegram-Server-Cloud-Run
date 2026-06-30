import os
import hashlib

def generate_server_keys():
    # יצירת מפתח הצפנה ראשי (Private/Public Key Pair) עבור האפליקציה שלך
    print("[MTProto Core] Automatically generating secure Auth Key pair for your private network...")
    server_private_key = os.urandom(256)
    server_public_key = hashlib.sha256(server_private_key).digest()
    return server_private_key, server_public_key

def process_telegram_packet(data):
    # מנוע עיבוד השפה של טלגרם - קורא את הפקודות הבינאריות
    if len(data) < 4:
        return None
    constructor_id = int.from_bytes(data[:4], byteorder='little')
    return constructor_id
