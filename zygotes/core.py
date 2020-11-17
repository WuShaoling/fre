import importlib
import os
import sys
import time
import traceback

import syscall


# import numpy
# import scipy
# import pandas
# import django
# import matplotlib


def do_exec(param):
    param["afterForkTime"] = time.time_ns()

    # 加载依赖包
    package_path = ['/usr/local/lib/python37.zip', '/usr/local/lib/python3.7',
                    '/usr/local/lib/python3.7/lib-dynload', '/usr/local/lib/python3.7/site-packages']
    for p in package_path:
        sys.path.append(p)
    sys.path.append(param["lambdaPath"])

    result = {
        "id": param["id"],
    }
    try:
        f = importlib.import_module(param["entrypoint"])
        result["result"] = f.handler(param["event"])
    except Exception as e:
        traceback.format_exc()
        result["error"] = e

    param["afterHandlerTime"] = time.time_ns()

    # TODO 结果写回 server
    print("firstFork=" + str((param["afterFirstForkTime"] - param["startTime"]) / 1e6) + ", " +
          "unshare=" + str((param["afterUnshareTime"] - param["afterFirstForkTime"]) / 1e6) + ", " +
          "chroot=" + str((param["afterChrootTime"] - param["afterUnshareTime"]) / 1e6) + ", " +
          "fork=" + str((param["afterForkTime"] - param["afterChrootTime"]) / 1e6) + ", " +
          "handler=" + str((param["afterHandlerTime"] - param["afterForkTime"]) / 1e6) + ", " +
          "total=" + str((param["afterHandlerTime"] - param["afterFirstForkTime"]) / 1e6)
          )

    sys.exit(0)


def new_container(param):
    try:
        param["afterFirstForkTime"] = time.time_ns()

        # unshare
        res = syscall.unshare()
        if res != 0:
            raise Exception("syscall.unshare return non zero status " + res)

        param["afterUnshareTime"] = time.time_ns()

        # # set cgroup
        # cur_pid = str(os.getpid())
        # for cgroup in param["cgroupFileList"]:
        #     f = open(cgroup, 'w')
        #     f.write(cur_pid)
        #     f.close()

        # chroot
        root_fd = os.open(param["rootPath"], os.O_RDONLY)
        os.fchdir(root_fd)
        os.chroot(".")
        os.close(root_fd)

        param["afterChrootTime"] = time.time_ns()

        # fork
        pid = os.fork()
        if pid == 0:  # child, 正式进入容器环境中
            do_exec(param)
        else:  # parent
            os.waitpid(pid, 0)
            # TODO 上报子进程退出
            sys.exit(0)
    except Exception as e:
        traceback.format_exc()
        # TODO 上报异常
        print(e)
        sys.exit(-1)


container_param = {
    "id": "aaa",
    "logFile": "/go/src/container/aaa.log",
    "lambdaPath": "/lambda",
    "entrypoint": "f",
    "cgroupFileList": ["/sys/fs/cgroup/memory/test1/tasks", "/sys/fs/cgroup/cpu/test1/tasks"],
    "rootPath": "/go/src/rootfs",
    "event": {},
    "startTime": time.time_ns()
}

first_pid = os.fork()
if first_pid == 0:
    new_container(container_param)
else:
    os.waitpid(first_pid, 0)
