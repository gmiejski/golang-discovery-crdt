from crdt_app.host_resolver import get_all_hosts

from requests import get


def test_data_consistency(first_host, local_data):
    print("Local values:")
    print(local_data)

    hosts = get_all_hosts(first_host)
    data_consistent = True

    all_results = [(host, get_data_from_host(host)) for host in hosts]
    for x in all_results:
        if x[1] != local_data:
            print("Not consistent data within host {}".format(x[0]))
            print(x[1])
            data_consistent = False
        else:
            print('Host consistent: {}'.format(x[0]))

    if not data_consistent:
        raise Exception("Data inconsistent!")


def get_data_from_host(host):
    current_value = get("{}/status/readable".format(host))
    json_rs = current_value.json()
    return sorted(list(map(lambda x: int(x), json_rs["Values"])))


if __name__ == "__main__":
    data = ["57","40","13","95","58","29","42","32","12","83","56","10","30","55","54","28","61","35","52","26","63","47","74","59","50","9","68","16","37","48","43","78","87","76","5","72","31","53","34","93","2","14","7","98","60","99","94","17","46","82"]
    int_data = sorted([int(x) for x in data])
    test_data_consistency("http://localhost:8080", int_data)