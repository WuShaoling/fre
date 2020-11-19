import http.client
import importlib
import json
import sys
import time
import traceback

# 解析参数
param = json.loads(sys.argv[-1])

# 加载执行函数代码
result = {
    "functionResult": {},
    "functionStartTime": int(round(time.time() * 1000000))
}
try:
    sys.path.append(param["functionPath"])
    handler = importlib.import_module(param["functionPath"])
    t = handler.handler(param["handler"])
    if t is not None:
        result["functionResult"] = t
except Exception as e:
    traceback.format_exc()
    result["error"] = str(e)
result["functionEndTime"] = int(round(time.time() * 1000000))

# 上报结果
conn = http.client.HTTPConnection(param["server"])
conn.request(
    "PUT",
    "/container/callback/result/" + param["id"],
    json.dumps(result, default=lambda obj: obj.__dict__),
    {'content-type': "application/json"})
print(conn.getresponse())
