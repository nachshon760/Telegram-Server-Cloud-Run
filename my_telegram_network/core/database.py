# Custom Server Database and Controls
import os

# בסיס נתונים וירטואלי - כאן השרת שלך שולט בהכל!
USERS_DATABASE = {
    "user_id_1": {"name": "Admin", "stars": 999999, "role": "Owner"},
    "user_id_2": {"name": "Guest", "stars": 50, "role": "User"}
}

def get_user_stars(user_id):
    # פונקציה המאפשרת לשרת שלך להחליט כמה כוכבים יש לכל משתמש באפליקציה
    user = USERS_DATABASE.get(user_id, {"stars": 0})
    return user["stars"]

def modify_user_stars(user_id, new_stars_count):
    if user_id in USERS_DATABASE:
        USERS_DATABASE[user_id]["stars"] = new_stars_count
        return True
    return False
