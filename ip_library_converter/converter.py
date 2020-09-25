import re
from os import listdir
import sys
import sqlite3
from pypinyin import Style, pinyin, lazy_pinyin
import ipaddress

def main(year_month):
    # print(year_month)
    path = "library/%s/省份-运营商-掩码" % year_month
    file_list = listdir(path)
    # print("_".join(lazy_pinyin("中國")))
    # print(lazy_pinyin("Other"))
    # exit()
    conn = sqlite3.connect('sqlite')
    c = conn.cursor()
    for file in file_list:
        state, telecom, _ = file.split('_')
        state = "".join(lazy_pinyin(state)).capitalize()
        telecom = "".join(lazy_pinyin(telecom)).capitalize()
        ip_list = []
        with open(path+"/"+file, 'r') as f:
            for line in f:
                ip = line.encode('utf-8').decode('utf-8-sig').strip()
                ip_list.append(ip)
                # mask = re.search(r"\/(.*)", ip)[1]
        c.execute('insert into ip_table (country, state, telecom, ip) values ("%s", "%s", "%s", "%s")' % ("China", state, telecom, ",".join(ip_list)))
        conn.commit()

if(len(sys.argv) > 1):
    main(sys.argv[1])