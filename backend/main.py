from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse
from fastapi.templating import Jinja2Templates
from influxdb_client import InfluxDBClient

app = FastAPI()
templates = Jinja2Templates(directory="backend/templates")

client = InfluxDBClient(
    url="http://localhost:8086",
    token="9VCJVflR7sYxZEFtJ_KwZdul8NHzvoDzMpoVAofPSjcvk7SSQQcyb2JweA7mSZAdLt6CKCV7SMmYbhyBYzUR2w==",
    org="parkme"
)
query_api = client.query_api()

@app.get("/", response_class=HTMLResponse)
def read_root(request: Request):
    return templates.TemplateResponse("index.html", {"request": request})

@app.get("/spots")
def get_spots():
    query = '''
    from(bucket: "lidar_data")
      |> range(start: -1m)
      |> filter(fn: (r) => r._measurement == "parking")
      |> filter(fn: (r) => r._field == "occupied")
      |> last()
    '''
    result = query_api.query(org="parkme", query=query)
    spots = []
    for table in result:
        for record in table.records:
            spots.append({
                "spot_id": record.values["spot_id"],
                "occupied": bool(record.get_value())
            })
    return spots

