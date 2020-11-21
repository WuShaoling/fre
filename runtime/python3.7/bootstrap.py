import http.client
import importlib
import json
import sys
import time
import traceback

# 记录服务启动时的时间
process_run_timestamp = int(round(time.time() * 1000000))

# 解析参数
param = json.loads(sys.argv[-1])

# 加载执行函数代码
result = {
    "id": param["id"],
    "containerProcessRunAt": process_run_timestamp
}
try:
    sys.path.append(param["codePath"])
    handler = importlib.import_module(param["handler"])
    result["functionRunTimestamp"] = int(round(time.time() * 1000000))
    result["functionResult"] = handler.handler(param["params"])
except Exception as e:
    traceback.format_exc()
    result["error"] = str(e)
result["functionEndTimestamp"] = int(round(time.time() * 1000000))

# 上报结果
address = "127.0.0.1:" + param["servePort"]
conn = http.client.HTTPConnection(address)
conn.request(
    "PUT",
    "/inner/function/end",
    json.dumps(result, default=lambda obj: obj.__dict__),
    {'content-type': "application/json"})
