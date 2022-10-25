import json
import math
import os

import numpy as np


class Test:
    def __init__(self, json_file_path: str):
        assert os.path.exists(json_file_path), "file {} not exist!"
        self.json_file_path = json_file_path

    def read_json_lines(self):
        f = open(self.json_file_path, 'r')
        lines = f.readlines()
        f.close()
        return lines

# 检查账户入金数额是否正常，让用户输入账户金额，和创建好之后的AccountVal进行比较
    def check_amount(self):
        lines = self.read_json_lines()
        if len(lines) == 0:
            print("Empty Json File!")
            return

        for line in lines:
            try:
                data = json.loads(line)  # data: dict
                data_keys = list(data.keys())  # get all dict keys
                if "AccountVal" in data_keys:
                    if data['message'] == 'Virtual Account Created!':
                        # 用户输入
                        usr_input = input("Please input the AccountVal: ")
                        if int(usr_input) != int(data['AccountVal']):
                            print("(User input){} != (AccountVal){}".format(int(usr_input), int(data['AccountVal'])))
            except Exception as e:
                print(e)

    def forward_look_up(self):
        pass

    @staticmethod
    def offset_look_up(lines: list, current_line_index: int, look_up_key: str, mode: str = 'forward'):
        assert mode in ['forward', 'backward']
        offset = 1
        if mode == 'forward':
            line = lines[current_line_index + offset]
        else:
            line = lines[current_line_index - offset]
        data = json.loads(line)
        while look_up_key not in list(data.keys()):
            offset += 1
            if mode == 'forward':
                line = lines[current_line_index + offset]
            else:  # backward
                line = lines[current_line_index - offset]
            data = json.loads(line)
        return data[look_up_key]
    
    # 账户接受买入下单指令记录持仓是否正常 买入后posidy（今日持仓数会有变化）
    def check_posi(self):
        lines = self.read_json_lines()
        if len(lines) == 0:
            print("Empty Json File!")
            return

        for line_index in range(len(lines)):
            line = lines[line_index]
            try:
                data = json.loads(line)  # data: dict
                if data['message'] == 'Strategy buy':
                    account_offset = 1
                    offset_line = lines[line_index + account_offset]
                    offset_data = json.loads(offset_line)
                    while offset_data['message'] != 'Account':
                        account_offset += 1
                        offset_line = lines[line_index + account_offset]
                        offset_data = json.loads(offset_line)
                    positdy = offset_data['positdy']
                    if positdy == 0:
                        print("买入下单后，持仓不正常，positdy: {}".format(positdy))
                        print("行号：{}".format(line_index + 1))
            except Exception as e:
                print(e)

# 收盘操作后，账户持仓是否进行相应调整测试
    def check_market_close(self):
        lines = self.read_json_lines()
        if len(lines) == 0:
            print("Empty Json File!")
            return

        for line_index in range(len(lines)):
            line = lines[line_index]
            try:
                data = json.loads(line)  # data: dict
                if data['message'] == 'Market Close':
                    previous_positdy = self.offset_look_up(lines, current_line_index=line_index,
                                                           look_up_key='positdy', mode='backward')
                    previous_posipre = self.offset_look_up(lines, current_line_index=line_index,
                                                           look_up_key='posipre', mode='backward')
                    current_positdy = self.offset_look_up(lines, current_line_index=line_index,
                                                           look_up_key='positdy', mode='forward')
                    current_posipre = self.offset_look_up(lines, current_line_index=line_index,
                                                           look_up_key='posipre', mode='forward')
                    if previous_positdy + previous_posipre != current_posipre or current_positdy != 0:
                        print("(previous_positdy){} \n (current_posipre){} \n (current_positdy){}".format(previous_positdy, current_posipre, current_positdy))
                        print("行号：{}".format(line_index + 1))
            except Exception as e:
                print(e)

# 账户是否区分今日持仓和过往持仓
    def check_difference(self):
        lines = self.read_json_lines()
        if len(lines) == 0:
            print("Empty Json File!")
            return

        for line_index in range(len(lines)):
            line = lines[line_index]
            try:
                data = json.loads(line)  # data: dict
                if data['message'] == 'Account':
                    keys = list(data.keys())
                    if "positdy" not in keys:
                        print("positdy not in the line {}".format(line_index + 1))
                    if "posipre" not in keys:
                        print("posipre not in the line {}".format(line_index + 1))
            except Exception as e:
                print(e)

# 检查账户是否不许卖出当日持仓 检查每天的posidy是否减少，今日持仓不需卖出的话，posidy只增不减；
    def check_in_one_day(self):
        lines = self.read_json_lines()
        if len(lines) == 0:
            print("Empty Json File!")
            return

        previous_line_index = None
        previous_timestamp = None
        previous_day = None
        previous_positdy = None

        for line_index in range(len(lines)):
            line = lines[line_index]
            try:
                data = json.loads(line)  # data: dict
                keys = list(data.keys())
                if "positdy" in keys:
                    timestamp = data['TimeStamp']
                    the_day = timestamp.split(' ')[0]
                    positdy = data['positdy']

                    if the_day == previous_day:
                        if positdy < previous_positdy:
                            print("the positdy of {} in line {} > {} in line {}".format(timestamp, line_index + 1,
                                                                                        previous_timestamp,
                                                                                        previous_line_index + 1))
                    previous_line_index = line_index
                    previous_timestamp = timestamp
                    previous_day = the_day
                    previous_positdy = positdy
            except Exception as e:
                print(e)

 # 检查账户针对行情信息进行持仓市值更新是否正常
    def check_shizhi(self):
        lines = self.read_json_lines()
        if len(lines) == 0:
            print("Empty Json File!")
            return

        for line_index in range(len(lines)):
            line = lines[line_index]
            try:
                data = json.loads(line)  # data: dict
                if data['message'] == "Account":
                    keys = list(data.keys())
                    if "Equity" in keys and "Equity4ALL" in keys:
                        equity = data['Equity']
                        equity4ALL = data['Equity4ALL']
                        if equity != equity4ALL:
                            print("Equity != Equity4ALL in line {}".format(line_index + 1))
            except Exception as e:
                print(e)

    # 账户接受买入下单指令资金扣除是否正常
    def check_buy(self):
        lines = self.read_json_lines()
        if len(lines) == 0:
            print("Empty Json File!")
            return

        for line_index in range(len(lines)):
            line = lines[line_index]
            try:
                data = json.loads(line)  # data: dict
                if data['message'] == "Strategy buy":
                    AccountVal = self.offset_look_up(lines, current_line_index=line_index, look_up_key="AccountVal", mode='forward')
                    Fundavail = self.offset_look_up(lines, current_line_index=line_index, look_up_key="Fundavail", mode='forward')
                    Equity4ALL = self.offset_look_up(lines, current_line_index=line_index, look_up_key="Equity4ALL", mode='forward')

                    total = "%.6f" % (Fundavail + Equity4ALL)
                    AccountVal = "%.6f" % AccountVal
                    if total != AccountVal:
                        print("line {} error !!".format(line_index + 1))
            except Exception as e:
                print(e)

# 账户接收卖出下单指令，持仓扣减是否正常
    def check_sell(self):
        lines = self.read_json_lines()
        if len(lines) == 0:
            print("Empty Json File!")
            return

        posipre_list = []
        for line_index in range(len(lines)):
            line = lines[line_index]
            try:
                data = json.loads(line)  # data: dict
                keys = list(data.keys())
                if "Market Close" in keys:
                    posipre_array = np.array(posipre_list)
                    diff = posipre_array[0:posipre_array.shape[0]] - posipre_array[1:]
                    if np.sum(diff < 0) > 0:
                        print("line error {}".format(line_index + 1))
                    posipre_list = []
                if "posipre" in keys:
                    posipre_list.append(data["posipre"])
            except Exception as e:
                print(e)

# 账户接收卖出下单指令，资金划拨是否正常
    def check_sell_money(self):
        lines = self.read_json_lines()
        if len(lines) == 0:
            print("Empty Json File!")
            return

        for line_index in range(len(lines)):
            line = lines[line_index]
            try:
                data = json.loads(line)  # data: dict
                if data['message'] == "Strategy sell":
                    AccountVal = self.offset_look_up(lines, current_line_index=line_index, look_up_key="AccountVal",
                                                     mode='forward')
                    Fundavail = self.offset_look_up(lines, current_line_index=line_index, look_up_key="Fundavail",
                                                    mode='forward')
                    Equity4ALL = self.offset_look_up(lines, current_line_index=line_index, look_up_key="Equity4ALL",
                                                     mode='forward')

                    total = "%.6f" % (Fundavail + Equity4ALL)
                    AccountVal = "%.6f" % AccountVal
                    if total != AccountVal:
                        print("line {} error !!".format(line_index + 1))
            except Exception as e:
                print(e)


if __name__ == '__main__':
    # json_file_path = 'D:\\实验室\\2022-9-30李家铖\\log.json'
    # json_file_path = 'D:\\实验室\\2022-9-30李家铖\\log2.json'
    json_file_path = 'D:\\实验室\\2022-9-30李家铖\\tmp.json'
    test = Test(json_file_path)
    # test.check_posi()
    # test.check_amount()
    # test.check_market_close()
    # test.check_difference()
    # test.check_in_one_day()
    # test.check_shizhi()
    # test.check_buy()
    # test.check_sell()
    # test.check_sell_money()