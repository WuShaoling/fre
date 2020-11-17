import time


# import numpy
# import scipy
# import pandas
# import django
# import matplotlib


def handler(event):
    # print("log from lambda: ", event)
    return time.time_ns()
