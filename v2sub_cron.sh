#!/bin/bash

# 更新全部订阅结果
v2sub -sub updateall

# 节点测速 + 设置最快节点
v2sub -ser setx

# 杀掉当前正在运行的 v2sub / v2ray
v2sub -conn kill

# 启动 v2ray 连接
v2sub -conn start

