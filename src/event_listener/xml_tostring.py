import json
import urllib.parse

import requests

xml_string = """<?xml version="1.0" encoding="UTF-8"?>
<definitions xmlns="http://www.omg.org/spec/DMN/20151101/dmn.xsd" xmlns:camunda="http://camunda.org/schema/1.0/dmn" id="dish" name="Dish" namespace="test-drd-2">
  <decision id="dish-decision" name="Dish Decision">
    <informationRequirement>
      <requiredDecision href="#season" />
    </informationRequirement>
    <informationRequirement>
      <requiredDecision href="#guestCount" />
    </informationRequirement>
    <decisionTable id="dishDecisionTable">
      <input id="seasonInput" label="Season">
        <inputExpression id="seasonInputExpression" typeRef="string">
          <text>season</text>
        </inputExpression>
      </input>
      <input id="guestCountInput" label="How many guests">
        <inputExpression id="guestCountInputExpression" typeRef="integer">
          <text>guestCount</text>
        </inputExpression>
      </input>
      <output id="output1" label="Dish" name="desiredDish" typeRef="string" />
      <rule id="row-495762709-1">
        <inputEntry id="UnaryTests_1nxcsjr">
          <text><![CDATA["Winter"]]></text>
        </inputEntry>
        <inputEntry id="UnaryTests_1r9yorj">
          <text><![CDATA[<=8]]></text>
        </inputEntry>
        <outputEntry id="LiteralExpression_1mtwzqz">
          <text><![CDATA["Spareribs"]]></text>
        </outputEntry>
      </rule>
      <rule id="row-495762709-2">
        <inputEntry id="UnaryTests_1lxjbif">
          <text><![CDATA["Winter"]]></text>
        </inputEntry>
        <inputEntry id="UnaryTests_0nhiedb">
          <text><![CDATA[>8]]></text>
        </inputEntry>
        <outputEntry id="LiteralExpression_1h30r12">
          <text><![CDATA["Pasta"]]></text>
        </outputEntry>
      </rule>
      <rule id="row-495762709-3">
        <inputEntry id="UnaryTests_0ifgmfm">
          <text><![CDATA["Summer"]]></text>
        </inputEntry>
        <inputEntry id="UnaryTests_12cib9m">
          <text><![CDATA[>10]]></text>
        </inputEntry>
        <outputEntry id="LiteralExpression_0wgaegy">
          <text><![CDATA["Light salad"]]></text>
        </outputEntry>
      </rule>
      <rule id="row-495762709-7">
        <inputEntry id="UnaryTests_0ozm9s7">
          <text><![CDATA["Summer"]]></text>
        </inputEntry>
        <inputEntry id="UnaryTests_0sesgov">
          <text><![CDATA[<=10]]></text>
        </inputEntry>
        <outputEntry id="LiteralExpression_1dvc5x3">
          <text><![CDATA["Beans salad"]]></text>
        </outputEntry>
      </rule>
      <rule id="row-445981423-3">
        <inputEntry id="UnaryTests_1er0je1">
          <text><![CDATA["Spring"]]></text>
        </inputEntry>
        <inputEntry id="UnaryTests_1uzqner">
          <text><![CDATA[<10]]></text>
        </inputEntry>
        <outputEntry id="LiteralExpression_1pxy4g1">
          <text><![CDATA["Stew"]]></text>
        </outputEntry>
      </rule>
      <rule id="row-445981423-4">
        <inputEntry id="UnaryTests_06or48g">
          <text><![CDATA["Spring"]]></text>
        </inputEntry>
        <inputEntry id="UnaryTests_0wa71sy">
          <text><![CDATA[>=10]]></text>
        </inputEntry>
        <outputEntry id="LiteralExpression_09ggol9">
          <text><![CDATA["Steak"]]></text>
        </outputEntry>
      </rule>
    </decisionTable>
  </decision>
  <decision id="season" name="Season decision">
    <decisionTable id="seasonDecisionTable">
      <input id="temperatureInput" label="Weather in Celsius">
        <inputExpression id="temperatureInputExpression" typeRef="integer">
          <text>temperature</text>
        </inputExpression>
      </input>
      <output id="seasonOutput" label="season" name="season" typeRef="string" />
      <rule id="row-495762709-5">
        <inputEntry id="UnaryTests_1fd0eqo">
          <text><![CDATA[>30]]></text>
        </inputEntry>
        <outputEntry id="LiteralExpression_0l98klb">
          <text><![CDATA["Summer"]]></text>
        </outputEntry>
      </rule>
      <rule id="row-495762709-6">
        <inputEntry id="UnaryTests_1nz6at2">
          <text><![CDATA[<10]]></text>
        </inputEntry>
        <outputEntry id="LiteralExpression_08moy1k">
          <text><![CDATA["Winter"]]></text>
        </outputEntry>
      </rule>
      <rule id="row-445981423-2">
        <inputEntry id="UnaryTests_1a0imxy">
          <text>[10..30]</text>
        </inputEntry>
        <outputEntry id="LiteralExpression_1poftw4">
          <text><![CDATA["Spring"]]></text>
        </outputEntry>
      </rule>
    </decisionTable>
  </decision>
  <decision id="guestCount" name="Guest Count">
    <decisionTable id="guestCountDecisionTable">
      <input id="typeOfDayInput" label="Type of day">
        <inputExpression id="typeOfDayInputExpression" typeRef="string">
          <text>dayType</text>
        </inputExpression>
      </input>
      <output id="guestCountOutput" label="Guest count" name="guestCount" typeRef="integer" />
      <rule id="row-495762709-8">
        <inputEntry id="UnaryTests_0l72u8n">
          <text><![CDATA["WeekDay"]]></text>
        </inputEntry>
        <outputEntry id="LiteralExpression_0wuwqaz">
          <text>4</text>
        </outputEntry>
      </rule>
      <rule id="row-495762709-9">
        <inputEntry id="UnaryTests_03a73o9">
          <text><![CDATA["Holiday"]]></text>
        </inputEntry>
        <outputEntry id="LiteralExpression_1whn119">
          <text>10</text>
        </outputEntry>
      </rule>
      <rule id="row-495762709-10">
        <inputEntry id="UnaryTests_12tygwt">
          <text><![CDATA["Weekend"]]></text>
        </inputEntry>
        <outputEntry id="LiteralExpression_1b5k9t8">
          <text>15</text>
        </outputEntry>
      </rule>
    </decisionTable>
  </decision>
</definitions>"""


def encode_decode_xml_string(xml_string):
    encoded_xml_string = urllib.parse.quote(xml_string)

    print(encoded_xml_string)
    print("\n\n")
    # 解码字符串
    decoded_str = urllib.parse.unquote(encoded_xml_string)

    print(decoded_str)

    print(decoded_str == xml_string)


def invoke_chaincode():
    # 目标URL
    url = "http://127.0.0.1:5000/api/v1/namespaces/default/apis/event2-test/invoke/CreateDMNContent"

    # JSON请求体
    request_body = {"input": {"dmnContent": xml_string, "id": "1"}}

    # 将JSON请求体转为字符串
    json_data = json.dumps(request_body)

    # Headers
    headers = {
        "Content-Type": "application/json",
    }

    # 发送POST请求
    response = requests.post(url, data=json_data, headers=headers)

    # 打印响应
    print(response.status_code)
    print(response.text)


# encode_decode_xml_string(xml_string)
invoke_chaincode()
