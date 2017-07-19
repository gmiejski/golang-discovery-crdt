def getHost(argv):
    print(argv)
    for idx in range(0, len(argv)):
        val = argv[idx]
        val = str(val)
        if val == "-H":
            return argv[idx+1]
        if val.startswith('--host'):
            return val.replace('--host=', '')
    print("Default port: http://localhost:8080")
    return "http://localhost:8080"


if __name__ == "__main__":
    import sys

    print(sys.argv)
    host = getHost(sys.argv)
    print(host)