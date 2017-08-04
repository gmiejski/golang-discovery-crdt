from datetime import datetime, timedelta


class Lwwes():
    def __init__(self):
        self.data_add = {}
        self.data_remove = {}

    def add(self, value, operation):
        if operation == "ADD":
            self.data_add[value] = datetime.now()
        else:
            self.data_remove[value] = datetime.now()

    def get(self):
        result = []
        for k, v in self.data_add.items():
            if k not in self.data_remove or self.data_remove[k] <= v:
                result.append(k)
        return result
