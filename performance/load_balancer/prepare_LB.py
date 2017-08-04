import sys

from crdt_app.host_resolver import get_all_hosts


def prepare_LB_config(hosts):
    servers = list(map(lambda x: "server {};".format(x.replace('http://','')), hosts))
    text = "\n".join(servers)
    text = text + "\n"

    with open("load_balancer/nginx_config.conf.template") as template:
        lines = template.readlines()
        final = ''.join(lines).replace("${servers}", text)

        with open("load_balancer/nginx_config.conf", mode="w") as target:
            target.write(final)


if __name__ == "__main__":
    first_host = "http://localhost:8080"
    if len(sys.argv) > 1:
        first_host = sys.argv[1]

    print("Getting host from {}".format(first_host))
    hosts = get_all_hosts(first_host)
    prepare_LB_config(hosts)