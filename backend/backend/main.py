from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
import sqlite3

app = FastAPI()

# Allow frontend access
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/spots")
def get_spots():
    conn = sqlite3.connect("backend/db.sqlite")
    cursor = conn.cursor()
    cursor.execute("SELECT spot_id, occupied FROM parking_status ORDER BY spot_id")
    data = cursor.fetchall()
    conn.close()
    return [{"spot_id": row[0], "occupied": bool(row[1])} for row in data]
