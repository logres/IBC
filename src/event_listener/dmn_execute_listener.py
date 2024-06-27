import asyncio
from time import sleep
import requests
import websockets
import json
import threading
from concurrent.futures import ThreadPoolExecutor

core_url = "http://127.0.0.1:5000/"
chaincode_url = f"""{core_url}api/v1/namespaces/default/apis/event4-test/"""
ipfs_url = "http://127.0.0.1:10207/ipfs/"


async def listen(executor, uri, listen_subscription_name, listen_action):
    async with websockets.connect(uri) as websocket:
        # 连接后发送一条消息
        message_to_send = {
            "type": "start",
            "name": listen_subscription_name,
            "namespace": "default",
            "autoack": True,
        }
        await websocket.send(json.dumps(message_to_send))
        print(f"Sent message: {message_to_send}")

        # 开始监听来自服务器的消息
        while True:
            try:
                message = await websocket.recv()
                print(f"Received message: {message}")
                # 在接收到消息后启动一个线程执行send_post_request
                loop = asyncio.get_event_loop()
                await loop.run_in_executor(executor, listen_action, message)
            except websockets.ConnectionClosed:
                print("Connection closed")
                break


# 创建线程池执行器
executor = ThreadPoolExecutor()
ws_uri = "ws://localhost:5000/ws"  # 替换为你的 WebSocket 服务器地址


def handle_read_dmn(message):
    # 将字符串转换为字典
    data = json.loads(message)
    eventContent = data.get("blockchainEvent").get("output")
    Id = eventContent.get("ID")
    Cid = eventContent.get("CID")
    dmn_content = read_from_ipfs(ipfs_url, Cid)
    invoke_dmn_contract(dmn_content, Id)


def read_from_ipfs(url, cid):
    response = requests.get(f"{url}{cid}")
    dmn_content = response.text
    return dmn_content


def invoke_dmn_contract(dmn_content, Id):
    # 目标URL
    url = f"""{chaincode_url}invoke/InvokeDmnContract"""
    # 构造请求
    response = requests.post(
        url,
        data=json.dumps({"input": {"dmnContent": dmn_content, "dmnId": Id}}),
        headers={"Content-Type": "application/json"},
    )
    print(response.text)


asyncio.get_event_loop().run_until_complete(
    listen(
        executor=executor,
        uri=ws_uri,
        listen_subscription_name="dmn_read4",
        listen_action=handle_read_dmn,
    )
)