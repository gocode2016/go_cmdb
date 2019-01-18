#-*-coding:utf8-*-

import os
import logging
import logging.handlers

def init_logger(log_file):
    dir_path = os.path.dirname(log_file)
    if not os.path.exists(dir_path):
        os.makedirs(dir_path)

    handler = logging.handlers.RotatingFileHandler(log_file, maxBytes=20 * 1024 * 1024, backupCount=5)
    fmt = '%(asctime)s %(filename)s[line:%(lineno)d] %(levelname)s %(message)s'
    #fmt = '%(asctime)s %(levelname)s %(message)s'
    formatter = logging.Formatter(fmt)
    handler.setFormatter(formatter)
    logger_instance = logging.getLogger('logs')
    logger_instance.addHandler(handler)
    logger_instance.setLevel(logging.DEBUG)
    return logger_instance
