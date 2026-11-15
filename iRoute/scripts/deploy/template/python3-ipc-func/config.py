import json
import os

config_file_path = "function/workflow.json"

class Config:
    def __init__(self):
        with open(config_file_path) as f:
            self.config = json.load(f)
        func_name = os.getenv("name", None)
        cc_name = os.getenv("cc", None)
        if func_name:
            self.func_name = func_name
            self.get_index()
        if cc_name:
            self.cc_name = cc_name

    def get_index(self):
        for stage_index, stage_fns in enumerate(self.config["dag"]):
            try:
                my_func_index = stage_fns.index(self.func_name)
                my_stage_index = stage_index
                break
            except ValueError as err:
                pass
        self.stage = my_stage_index
        self.index = my_func_index

        if self.stage > 0:
            self.fan_in = len(self.config["dag"][self.stage - 1])
        else:
            self.fan_in = 1

    def get_funcs(self):
        fns = []
        for stage_fns in self.config["dag"]:
            fns.extend(stage_fns)
        return fns
    
    def get_servers(self):
        try:
            if self.stage > 0:
                return self.config["dag"][self.stage - 1]
            else:
                return []
        except Exception as e:
            return []

    def get_clients(self):
        try:
            if self.stage < len(self.config["dag"]) - 1:
                return self.config["dag"][self.stage + 1]
            else:
                return []
        except Exception as e:
            return []

    # if next stage is fan in
    def get_fan_in(self):
        if self.stage < len(self.config["dag"]) - 1 and len(self.config["dag"][self.stage]) > 1:
            return [0]
        else:
            return []
