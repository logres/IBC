#
# SPDX-License-Identifier: Apache-2.0
#
import os
import json
import subprocess
from api.lib.peer.command import Command
from api.config import FABRIC_TOOL


class Channel(Command):
    """Call CMD to perform channel create, join and other related operations"""

    def __init__(self, version="2.2.0", peer=FABRIC_TOOL, **kwargs):
        self.peer = peer + "/peer"
        super(Channel, self).__init__(version, **kwargs)

    def create(self, channel, orderer_url, channel_tx, output_block, time_out="90s"):
        try:
            res = 0x100
            if (
                os.getenv("CORE_PEER_TLS_ENABLED") == "false"
                or os.getenv("CORE_PEER_TLS_ENABLED") is None
            ):
                res = os.system(
                    "{} channel create -c {} -o {} -f {} --outputBlock {} --timeout {}".format(
                        self.peer,
                        channel,
                        orderer_url,
                        channel_tx,
                        output_block,
                        time_out,
                    )
                )
            else:
                ORDERER_CA = os.getenv("ORDERER_CA")
                print(
                    "{} channel create -c {} -o {} -f {} --outputBlock {} --timeout {} --tls --cafile {}".format(
                        self.peer,
                        channel,
                        orderer_url,
                        channel_tx,
                        output_block,
                        time_out,
                        ORDERER_CA,
                    )
                )
                res = os.system(
                    "{} channel create -c {} -o {} -f {} --outputBlock {} --timeout {} --tls --cafile {}".format(
                        self.peer,
                        channel,
                        orderer_url,
                        channel_tx,
                        output_block,
                        time_out,
                        ORDERER_CA,
                    )
                )

            # The return value of os.system is not the result of executing the program. It is a 16 bit number,
            #  and its high bit is the return code
            res = res >> 8
        except Exception as e:
            err_msg = "create channel failed for {}!".format(e)
            raise Exception(err_msg)
        return res

    def list(self):
        try:
            res = subprocess.Popen(
                "{} channel list".format(self.peer),
                shell=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )

            stdout, stderr = res.communicate()
            return_code = res.returncode

            if return_code == 0:
                content = str(stdout, encoding="utf-8")
                content = content.split("\n")
            else:
                stderr = str(stderr, encoding="utf-8")
                return return_code, stderr
        except Exception as e:
            err_msg = "get channel list failed for {}!".format(e)
            raise Exception(err_msg)
        return return_code, content[1:-1]

    # def update(self, channel, channel_tx, orderer_url):
    #     """
    #     Send a configtx update.
    #     params:
    #         channel: channel id.
    #         channel_tx: Configuration transaction file generated by a tool such as configtxgen for submitting to orderer
    #         orderer_url: Ordering service endpoint.
    #     """
    #     try:
    #         res = os.system(
    #             "{} channel update -c {}  -f {} -o {}".format(
    #                 self.peer, channel, channel_tx, orderer_url
    #             )
    #         )
    #     except Exception as e:
    #         err_msg = "update channel failed for {}!".format(e)
    #         raise Exception(err_msg)
    #     res = res >> 8
    #     return res

    def fetch(self, option, channel, path):
        """
        Fetch a specified block, writing it to a file e.g. <channelID>.block.
        params:
            option: block option newest|oldest|config|(block number).
            channel: channel id.
        """
        try:
            res = subprocess.call(
                args=[
                    self.peer,
                    "channel",
                    "fetch",
                    "{}".format(option),
                    path,
                    "-c",
                    channel,
                ]
            )
        except Exception as e:
            err_msg = "fetch a specified block failed {}!".format(e)
            raise Exception(err_msg)
        return res

    def signconfigtx(self, channel_tx):
        """
        Signs a configtx update.
        params:
            channel_tx: Configuration transaction file generated by a tool such as configtxgen for submitting to orderer
        """
        try:
            res = os.system(
                "{} channel signconfigtx -f {}".format(self.peer, channel_tx)
            )
        except Exception as e:
            err_msg = "signs a configtx update failed {}".format(e)
            raise Exception(err_msg)
        res = res >> 8
        return res

    def update(self, channel_tx, channel_id, orderer_host_port):
        """
        Signs and sends the supplied configtx update file to the channel. Requires '-f', '-o', '-c'.
        params:
            channel: channel id.
            channel_tx: Configuration transaction file generated by a tool such as configtxgen for submitting to orderer
            orderer_host_port: Ordering service endpoint.
        """
        try:
            command = "{} channel update -f {} -c {} -o {} --ordererTLSHostnameOverride {} --tls --cafile {}".format(
                self.peer,
                channel_tx,
                channel_id,
                orderer_host_port,
                orderer_host_port.split(":")[0],
                os.getenv("ORDERER_CA"),
            )
            print(command)
            res = os.system(command)
        except Exception as e:
            err_msg = "update a configtx failed {}".format(e)
            raise Exception(err_msg)
        res = res >> 8
        return res

    def join(self, block_file):
        """
        Joins the peer to a channel.
        params:
            block_file: Path to file containing genesis block.
        """
        try:
            res = os.system("{} channel join -b {} ".format(self.peer, block_file))
        except Exception as e:
            err_msg = "join the peer to a channel failed. {}".format(e)
            raise Exception(err_msg)
        res = res >> 8
        return res

    def getinfo(self, channel):
        """
        Get blockchain information of a specified channel.
        params:
            channel: In case of a newChain command, the channel ID to create.
        """
        try:
            res = subprocess.Popen(
                "{} channel getinfo  -c {}".format(self.peer, channel),
                shell=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )

            stdout, stderr = res.communicate()
            return_code = res.returncode

            if return_code == 0:
                content = str(stdout, encoding="utf-8")
                content = content.split("\n")[0].split(":", 1)[1]
                block_info = json.loads(content)
                body = {"block_info": block_info}
            else:
                stderr = str(stderr, encoding="utf-8")
                return return_code, stderr
        except Exception as e:
            err_msg = (
                "get blockchain information of a specified channel failed. {}".format(e)
            )
            raise Exception(err_msg)
        return return_code, body
