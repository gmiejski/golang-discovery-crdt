from requests import get

def get_all_hosts(first_host):
    result = get("{}/cluster/info".format(first_host))

    if result.status_code != 200:
        raise Exception("Wrong response from server: {}".format(result.status_code))

    nodes = result.json()["Nodes"]
    hosts = [node["Url"] for node in nodes] + [result.json()['NodeUrl']]
    return hosts
