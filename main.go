package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gotd/td/transport"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 1. הגדרת מערכת לוגים מפורטת (Verbose Logging) כדי לראות כל תנועה בשפה של טלגרם
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	defer logger.Sync()

	logger.Info("מפעיל שרת טלגרם מותאם אישית...")

	// 2. קבלת הפורט מ-Google Cloud Run (ברירת מחדל 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 3. יצירת מאזין לפרוטוקול MTProto
	// השרת משתמש ב-Intermediate Transport (אחד מהסטנדרטים של טלגרם להעברת בתים)
	listener := transport.Intermediate(nil)

	// 4. פונקציה לטיפול בכל לקוח (אפליקציה) שמתחבר
	http.HandleFunc("/tg", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("=== התקבלה בקשת חיבור חדשה מאפליקציית הטלגרם ===")
		logger.Info(fmt.Sprintf("כתובת המקור: %s | סוג הבקשה: %s", r.RemoteAddr, r.Method))

		// כאן השרת מקבל את החיבור ומפעיל את ה-Handshake לפענוח ויצירת מפתחות ההצפנה
		// הערה: קוד ה-Handshake המלא דורש מימוש של פרוטוקול Diffie-Hellman עם מפתח RSA שמוטמע באפליקציה שלך.
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("MTProto Server Listener Active"))
	})

	// 5. הרצת השרת ב-Cloud Run
	serverAddr := fmt.Sprintf(":%s", port)
	logger.Info(fmt.Sprintf("השרת מקשיב כעת בכתובת: %s", serverAddr))
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		logger.Fatal("השרת קרס בגלל שגיאה:", zap.Error(err))
	}
}
